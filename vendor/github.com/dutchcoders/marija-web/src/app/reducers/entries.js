import { slice, concat, without, reduce, remove, assign, find, forEach, union, filter, uniqBy } from 'lodash';

import {  ERROR, AUTH_CONNECTED, Socket, SearchMessage, DiscoverIndicesMessage, DiscoverFieldsMessage } from '../utils/index';

import {  INDICES_RECEIVE, INDICES_REQUEST } from '../modules/indices/index';
import {  FIELDS_RECEIVE, FIELDS_REQUEST } from '../modules/fields/index';
import {  NODES_DELETE, NODES_HIGHLIGHT, NODE_UPDATE, NODE_SELECT, NODES_SELECT, NODES_DESELECT, SELECTION_CLEAR } from '../modules/graph/index';
import {  SEARCH_DELETE, ITEMS_RECEIVE, ITEMS_REQUEST } from '../modules/search/index';
import {  TABLE_COLUMN_ADD, TABLE_COLUMN_REMOVE, INDEX_ADD, INDEX_DELETE, FIELD_ADD, FIELD_DELETE, DATE_FIELD_ADD, DATE_FIELD_DELETE, NORMALIZATION_ADD, NORMALIZATION_DELETE, INITIAL_STATE_RECEIVE } from '../modules/data/index';

import { normalize, fieldLocator } from '../helpers/index';
import getNodesAndLinks from "../helpers/getNodesAndLinks";
import removeNodesAndLinks from "../helpers/removeNodesAndLinks";
import getHighlightItem from "../helpers/getHighlightItem";


export const defaultState = {
    isFetching: false,
    itemsFetching: false,
    noMoreHits: false,
    didInvalidate: false,
    connected: false,
    total: 0,
    node: [],
    datasources: [],
    highlight_nodes: [],
    columns: [],
    fields: [],
    date_fields: [],
    normalizations: [], 
    indexes: [],
    items: [],
    searches: [],
    nodes: [], // all nodes
    links: [], // relations between nodes
    errors: null
};


export default function entries(state = defaultState, action) {
    switch (action.type) {
        case SELECTION_CLEAR:
            return Object.assign({}, state, {
                node: [],
            });
        case INDEX_DELETE:
            const index = find(state.indexes, (i) => {
                return (i.id == action.index);
            });

            var indexes = without(state.indexes, index);
            return Object.assign({}, state, {
                indexes: indexes,
            });
        case NODES_DELETE:
            var items = concat(state.items);
            var node = concat(state.node);
            var nodes = concat(state.nodes);
            var links = concat(state.links);

            // remove from selection as well
            remove(node, (p) => {
                return find(action.nodes, (o) => {
                    return o.id == p.id;
                });
            });

            remove(nodes, (p) => {
                return find(action.nodes, (o) => {
                    return o.id == p.id;
                });
            });

            remove(links, (p) => {
                return find(action.nodes, (o) => {
                    return p.source == o.id || p.target == o.id;
                });
            });

            return Object.assign({}, state, {
                items: items,
                node: node,
                nodes: nodes,
                links: links
            });
        case SEARCH_DELETE:
            const searches = without(state.searches, action.search);
            var items = concat(state.items);

            items = items.filter(item => item.query !== action.search.q);

            // todo(nl5887): remove related nodes and links
            const result = removeNodesAndLinks(state.nodes, state.links, action.search.q);

            return Object.assign({}, state, {
                searches: searches,
                items: items,
                nodes: result.nodes,
                links: result.links,
                selectedNodes: [],
                highlight_nodes: [],
                node: []
            });
        case TABLE_COLUMN_ADD:
            return Object.assign({}, state, {
                columns: concat(state.columns, action.field),
            });
        case TABLE_COLUMN_REMOVE:
            return Object.assign({}, state, {
                columns: without(state.columns, action.field)
            });
        case FIELD_ADD:
            const existing = state.fields.find(field => field.path === action.path);

            if (existing) {
                // Field was already in store, don't add duplicates
                return state;
            }

            const firstChar = action.path.charAt(0).toUpperCase();
            const fieldsWithSameChar = state.fields.filter(field => field.icon.indexOf(firstChar) === 0);
            let icon;

            if (fieldsWithSameChar.length === 0) {
                icon = firstChar;
            } else {
                // Append a number to the icon if multiple fields share the same
                // first character
                icon = firstChar + (fieldsWithSameChar.length + 1);
            }

            let newField = {
                path: action.path,
                icon: icon
            };

            return Object.assign({}, state, {
                fields: concat(state.fields, newField)
            });
        case FIELD_DELETE:
            return Object.assign({}, state, {
                fields: without(state.fields, action.field)
            });
        case NORMALIZATION_ADD:
            let normalization = action.normalization;
            normalization.re = new RegExp(normalization.regex, "i");
            return Object.assign({}, state, {
                normalizations: concat(state.normalizations, normalization)
            });
        case NORMALIZATION_DELETE:
            return Object.assign({}, state, {
                normalizations: without(state.normalizations, action.normalization)
            });
        case DATE_FIELD_ADD:
            return Object.assign({}, state, {
                date_fields: concat(state.date_fields, action.field)
            });
        case DATE_FIELD_DELETE:
            return Object.assign({}, state, {
                date_fields: without(state.date_fields, action.field)
            });
        case NODES_HIGHLIGHT:
            const highlightItems = [];

            action.highlight_nodes.forEach(node => {
                const item = Object.assign({}, state.items.find(item => item.id === node.items[0]));
                const highlightItem = getHighlightItem(item, node, state.fields);

                highlightItems.push(highlightItem);
            });

            return Object.assign({}, state, {
                highlight_nodes: highlightItems
            });
        case NODES_SELECT:
            const newNodes = [];

            action.nodes.forEach(node => {
                // First check if it's already selected, don't add duplicates
                if (state.node.indexOf(node) === -1) {
                    newNodes.push(node);
                }
            });

            return Object.assign({}, state, {
                node: concat(state.node, newNodes)
            });
        case NODES_DESELECT:
            return Object.assign({}, state, {
                node: filter(state.node, (o) => {
                    return !find(action.nodes, o);
                })
            });
        case NODE_UPDATE:
            let nodes = concat(state.nodes, []);

            let n = find(nodes, {id: action.node_id });
            if (n) {
                n = assign(n, action.params);
            }

            return Object.assign({}, state, {
                nodes: nodes
            });
        case NODE_SELECT:
            return Object.assign({}, state, {
                node: concat(state.node, action.node)
            });
        case ERROR:
            console.debug(action);
            return Object.assign({}, state, {
                errors: action.errors
            });
        case AUTH_CONNECTED:
            return Object.assign({}, state, {
                isFetching: false,
                didInvalidate: false,
                ...action
            });
        case ITEMS_REQUEST: {
            // if we searched before, just retrieve extra results for query
            const search = find(state.searches, (o) => o.q == action.query) || { items: [] };

            let from = search.items.length || 0;

            let message = {
                datasources: action.datasources,
                query: action.query,
                from: from,
                size: action.size,
                // todo: remove this, but requires a backend change. now the server wont respond if we dont send a color
                color: '#de79f2'
            };
            Socket.ws.postMessage(message);

            return Object.assign({}, state, {
                isFetching: true,
                itemsFetching: true,
                didInvalidate: false
            });
        }
        case ITEMS_RECEIVE: {
            const searches = concat(state.searches, []);
            const items = action.items.results === null ? [] : action.items.results;

            // should we update existing search, or add new, do we still need items?
            let search = find(state.searches, (o) => o.q == action.items.query);
            if (search) {
                search.items = concat(search.items, action.items.results);
            } else {
                const colors = [
                    '#de79f2',
                    '#917ef2',
                    '#499df2',
                    '#49d6f2',
                    '#00ccaa',
                    '#fac04b',
                    '#bf8757',
                    '#ff884d',
                    '#ff7373',
                    '#ff5252',
                    '#6b8fb3'
                ];
                // Sequentially uses the available colors, and starts again from the start when we exceed the amount of colors
                const colorIndex = (state.searches.length % colors.length + colors.length) % colors.length;
                const color = colors[colorIndex];

                search = {
                    q: action.items.query,
                    color: color,
                    total: action.items.total,
                    items: items
                };

                searches.push(search);
            }

            // Save per item for which query we received it (so we can keep track of where data came from)
            items.forEach(item => {
                item.query = search.q;
            });

            // todo(nl5887): should we start a webworker here, the webworker can have its own permanent cache?

            // update nodes and links
            const {nodes, links} = getNodesAndLinks(
                state.nodes,
                state.links,
                items,
                state.fields,
                search,
                state.normalizations
            );

            return Object.assign({}, state, {
                errors: null,
                nodes: nodes,
                links: links,
                items: concat(state.items, items),
                searches: searches,
                isFetching: false,
                itemsFetching: false,
                didInvalidate: false
            });
        }
        case INDICES_REQUEST:
            Socket.ws.postMessage(
                {
                    host: [action.payload.server]
                },
                INDICES_REQUEST
            );

            return Object.assign({}, state, {
                isFetching: true,
                didInvalidate: false
            });

        case INDICES_RECEIVE:
            const indices = uniqBy(union(state.indexes, action.payload.indices.map((index) => {
                return {
                    id: `${action.payload.server}${index}`,
                    server: action.payload.server,
                    name: index
                };
            })), (i) => i.id);

            return Object.assign({}, state, {
                indexes: indices,
                isFetching: false
            });

        case FIELDS_REQUEST:
            return Object.assign({}, state, {
                isFetching: true
            });

        case FIELDS_RECEIVE:
            return Object.assign({}, state, {
                isFetching: false
            });

        case INITIAL_STATE_RECEIVE:
        return Object.assign({}, state, {
            datasources: action.initial_state.datasources,
        });

        default:
            return state;
    }
}
