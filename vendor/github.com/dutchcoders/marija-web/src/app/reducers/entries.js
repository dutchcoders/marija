import { slice, concat, without, reduce, remove, assign, find, forEach, union, filter, uniqBy, uniqueId } from 'lodash';

import {  ERROR, AUTH_CONNECTED, Socket, SearchMessage, DiscoverIndicesMessage, DiscoverFieldsMessage } from '../utils/index';

import {  INDICES_RECEIVE, INDICES_REQUEST } from '../modules/indices/index';
import {  FIELDS_RECEIVE, FIELDS_REQUEST } from '../modules/fields/index';
import {  NODES_DELETE, NODES_HIGHLIGHT, NODE_UPDATE, NODE_SELECT, NODES_SELECT, NODES_DESELECT, SELECTION_CLEAR } from '../modules/graph/index';
import {  SEARCH_DELETE, SEARCH_RECEIVE, SEARCH_REQUEST, SEARCH_COMPLETED, SET_DISPLAY_NODES } from '../modules/search/index';
import {  TABLE_COLUMN_ADD, TABLE_COLUMN_REMOVE, INDEX_ADD, INDEX_DELETE, FIELD_ADD, FIELD_DELETE, DATE_FIELD_ADD, DATE_FIELD_DELETE, NORMALIZATION_ADD, NORMALIZATION_DELETE, INITIAL_STATE_RECEIVE } from '../modules/data/index';

import {
    normalize, fieldLocator, getNodesForDisplay,
    removeDeadLinks, applyVia, getQueryColor
} from '../helpers/index';
import getNodesAndLinks from "../helpers/getNodesAndLinks";
import removeNodesAndLinks from "../helpers/removeNodesAndLinks";
import getHighlightItem from "../helpers/getHighlightItem";
import {VIA_ADD, VIA_DELETE} from "../modules/data/constants";


export const defaultState = {
    isFetching: false,
    itemsFetching: false,
    noMoreHits: false,
    didInvalidate: false,
    connected: false,
    total: 0,
    node: [],
    datasources: [],
    highlight_nodes: {},
    columns: [],
    fields: [],
    date_fields: [],
    normalizations: [], 
    indexes: [],
    items: [],
    searches: [],
    nodes: [], // all nodes
    links: [], // relations between nodes
    nodesForDisplay: [], // nodes that will be rendered
    linksForDisplay: [], // links that will be rendered
    errors: null,
    via: [],
    version: ''
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
        case NODES_DELETE: {
            const items = concat([], state.items);
            const node = concat([], state.node);
            const nodes = concat([], state.nodes);
            const links = concat([], state.links);

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

            const nodesForDisplay = getNodesForDisplay(nodes, state.searches);
            const linksForDisplay = removeDeadLinks(nodesForDisplay, links);

            return Object.assign({}, state, {
                items: items,
                node: node,
                nodes: nodes,
                links: links,
                nodesForDisplay: nodesForDisplay,
                linksForDisplay: linksForDisplay
            });
        }
        case SEARCH_DELETE:
            const searches = without(state.searches, action.search);
            var items = concat(state.items);

            if (!action.search.completed) {
                // Tell the server it can stop sending results for this query
                Socket.ws.postMessage(
                    {
                        'request-id': action.search['request-id']
                    },
                    INDICES_REQUEST
                );
            }

            items = items.filter(item => item.query !== action.search.q);

            // todo(nl5887): remove related nodes and links
            const result = removeNodesAndLinks(state.nodes, state.links, action.search.q);
            const nodesForDisplay = getNodesForDisplay(result.nodes, state.searches);
            const linksForDisplay = removeDeadLinks(nodesForDisplay, result.links);

            return Object.assign({}, state, {
                searches: searches,
                items: items,
                nodes: result.nodes,
                links: result.links,
                nodesForDisplay: nodesForDisplay,
                linksForDisplay: linksForDisplay,
                selectedNodes: [],
                highlight_nodes: {},
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
        case VIA_ADD:
            return Object.assign({}, state, {
                via: concat(state.via, action.via)
            });
        case VIA_DELETE:
            return Object.assign({}, state, {
                via: without(state.via, action.via)
            });
        case DATE_FIELD_ADD: {
            const existing = state.date_fields.find(search => search.path === action.field.path);

            if (typeof existing !== 'undefined') {
                return state;
            }

            return Object.assign({}, state, {
                date_fields: concat(state.date_fields, action.field)
            });
        }
        case DATE_FIELD_DELETE:
            return Object.assign({}, state, {
                date_fields: without(state.date_fields, action.field)
            });
        case NODES_HIGHLIGHT:
            const highlightItems = {};

            forEach(action.highlight_nodes, node => {
                const item = Object.assign({}, state.items.find(item => item.id === node.items[0]));
                highlightItems[node.hash] = getHighlightItem(item, node, state.fields, 50);
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
        case SEARCH_REQUEST: {
            // if we searched before, just retrieve extra results for query
            // const search = find(state.searches, (o) => o.q == action.query) || { items: [] };
            const searches = concat(state.searches, []);

            let search = find(state.searches, (o) => o.q === action.query);

            if (!search) {
                const color = getQueryColor(state.searches);

                search = {
                    q: action.query,
                    color: color,
                    total: 0,
                    displayNodes: action.displayNodes,
                    items: [],
                    requestId: uniqueId(),
                    completed: false
                };

                searches.push(search);
            }

            let from = search.items.length || 0;

            const fieldPaths = [];
            action.fields.forEach(field => fieldPaths.push(field.path));

            let message = {
                datasources: action.datasources,
                query: action.query,
                from: from,
                size: action.size,
                // todo: remove this, but requires a backend change. now the server wont respond if we dont send a color
                color: '#de79f2',
                fields: fieldPaths,
                'request-id': search.requestId
            };
            Socket.ws.postMessage(message);

            return Object.assign({}, state, {
                isFetching: true,
                itemsFetching: true,
                didInvalidate: false,
                searches: searches
            });
        }
        case SEARCH_RECEIVE: {
            const searches = concat(state.searches, []);
            const items = action.items.results === null ? [] : action.items.results;

            // should we update existing search, or add new, do we still need items?
            let search = find(state.searches, (o) => o.q === action.items.query);
            if (search) {
                search.items = concat(search.items, action.items.results);
            } else {
                console.error('received items for a query we were not searching for: ' + action.items.query);
                return state;
            }

            // Save per item for which query we received it (so we can keep track of where data came from)
            items.forEach(item => {
                item.query = search.q;
            });

            // todo(nl5887): should we start a webworker here, the webworker can have its own permanent cache?

            // update nodes and links
            const result = getNodesAndLinks(
                state.nodes,
                state.links,
                items,
                state.fields,
                search,
                state.normalizations
            );
            const { nodes, links } = applyVia(result.nodes, result.links, state.via);
            const nodesForDisplay = getNodesForDisplay(nodes, state.searches || []);

            let linksForDisplay;

            if (nodesForDisplay.length < nodes.length) {
                linksForDisplay = removeDeadLinks(nodesForDisplay, links);
            } else {
                linksForDisplay = links;
            }

            return Object.assign({}, state, {
                errors: null,
                nodes: nodes,
                links: links,
                nodesForDisplay: nodesForDisplay,
                linksForDisplay: linksForDisplay,
                items: concat(state.items, items),
                searches: searches,
                isFetching: false,
                itemsFetching: false,
                didInvalidate: false
            });
        }
        case SEARCH_COMPLETED: {
            const index = state.searches.findIndex(search => search['request-id'] === action['request-id']);

            if (index === -1) {
                // Could not find out which search was completed
                return state;
            }

            const search = state.searches[index];
            const newSearch = Object.assign(search, { completed: true });
            const newSearches = concat([], state.searches);

            newSearches[index] = newSearch;

            return Object.assign({}, state, {
                searches: newSearches
            });
        }
        case SET_DISPLAY_NODES: {
            const searches = concat([], state.searches);

            const search = state.searches.find(search => search.q === action.query);
            const newSearch = Object.assign({}, search, {
                displayNodes: action.newAmount
            });

            const index = searches.indexOf(search);
            searches[index] = newSearch;

            const nodesForDisplay = getNodesForDisplay(state.nodes, searches);
            const linksForDisplay = removeDeadLinks(nodesForDisplay, state.links);

            return Object.assign({}, state, {
                searches: searches,
                nodesForDisplay: nodesForDisplay,
                linksForDisplay: linksForDisplay
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
                version: action.initial_state.version
            });

        default:
            return state;
    }
}
