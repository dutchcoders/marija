import { slice, concat, without, reduce, remove, assign, find, forEach, union, filter, uniqBy } from 'lodash';

import {  ERROR, AUTH_CONNECTED, Socket, SearchMessage, DiscoverIndicesMessage, DiscoverFieldsMessage } from '../utils/index';

import {  INDICES_RECEIVE, INDICES_REQUEST } from '../modules/indices/index';
import {  FIELDS_RECEIVE, FIELDS_REQUEST } from '../modules/fields/index';
import {  NODES_DELETE, NODES_HIGHLIGHT, NODE_UPDATE, NODE_SELECT, NODES_SELECT, NODES_DESELECT, SELECTION_CLEAR } from '../modules/graph/index';
import {  SEARCH_DELETE, ITEMS_RECEIVE, ITEMS_REQUEST } from '../modules/search/index';
import {  TABLE_COLUMN_ADD, TABLE_COLUMN_REMOVE, INDEX_ADD, INDEX_DELETE, FIELD_ADD, FIELD_DELETE, DATE_FIELD_ADD, DATE_FIELD_DELETE, NORMALIZATION_ADD, NORMALIZATION_DELETE, INITIAL_STATE_RECEIVE } from '../modules/data/index';

import { normalize, fieldLocator } from '../helpers/index';


export const defaultState = {
    isFetching: false,
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
            var searches = without(state.searches, action.search);

            var items = concat(state.items);
            remove(items, (p) => {
                return (p.q === action.search.q);
            });

            // todo(nl5887): remove related nodes and links

            return Object.assign({}, state, {
                searches: searches,
                items: items
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
            return Object.assign({}, state, {
                fields: concat(state.fields, action.field)
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
            return Object.assign({}, state, {
                highlight_nodes: action.highlight_nodes
            });
        case NODES_SELECT:
            return Object.assign({}, state, {
                node: concat(action.nodes)
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

            let message = {datasources: action.datasources, query: action.query, from: from, size: action.size, color: action.color};
            Socket.ws.postMessage(message);

            return Object.assign({}, state, {
                isFetching: true,
                didInvalidate: false
            });
        }
        case ITEMS_RECEIVE: {
            var searches = concat(state.searches, []);

            // should we update existing search, or add new, do we still need items?
            let search = find(state.searches, (o) => o.q == action.items.query);
            if (search) {
                search.items = concat(search.items, action.items.results);
            } else {
                searches.push({
                    q: action.items.query,
                    color: action.items.color,
                    total: action.items.total,
                    items: action.items.results
                });
            }

            // let search = find(searches, (o) => return o.q == action.items.query) {

            const { normalizations } = state;


            // todo(nl5887): should we start a webworker here, the webworker can have its own permanent cache?

            // update nodes and links
            var items = action.items.results;

            var nodes = concat(state.nodes, []);
            var links = concat(state.links, []);

            let nodeCache = {};
            for (let node of nodes) {
                nodeCache[node.id] = node;
            }

            let linkCache = {};
            for (let link of links) {
                linkCache[link.source + link.target] = link;
            }

            const fields = state.fields;
            forEach(items, (d, i) => {
                forEach(fields, (source) => {
                    let sourceValue = fieldLocator(d.fields, source.path);
                    if (sourceValue === null) {
                        return;
                    }

                    if (!Array.isArray(sourceValue)) {
                        sourceValue = [sourceValue];
                    }

                    for (let sv of sourceValue) {
                        switch (typeof sv) {
                        case "boolean":
                            sv = (sv?"true":"false");
                        }

                        const normalizedSourceValue = normalize(normalizations, sv);
                        if (normalizedSourceValue === "") {
                            continue;
                        }

                        let n = nodeCache[normalizedSourceValue];
                        if (n) {
                            if (n.items.indexOf(d.id) == -1){
                                n.items.push(d.id);
                            }

                            if (n.fields.indexOf(source.path) == -1){
                                n.fields.push(source.path);
                            }

                            n.queries.push(action.items.query);
                        } else {
                            let n = {
                                id: normalizedSourceValue,
                                queries: [action.items.query],
                                items: [d.id],
                                name: normalizedSourceValue,
                                description: '',
                                icon: source.icon,
                                fields: [source.path],
                            };

                            nodeCache[n.id] = n;
                            nodes.push(n);
                        }

                        forEach(fields, (target) => {
                            let targetValue = fieldLocator(d.fields, target.path);
                            if (targetValue === null) {
                                return;
                            }

                            if (!Array.isArray(targetValue)) {
                                targetValue = [targetValue];
                            }

                            // todo(nl5887): issue with normalizing is if we want to use it as name as well.
                            // for example we don't want to have the first name only as name.
                            //
                            // we need to keep track of the fields the value is in as well.
                            for (let tv of targetValue) {
                                switch (typeof tv) {
                                case "boolean":
                                    tv = (tv?"true":"false");
                                }

                                const normalizedTargetValue = normalize(normalizations, tv);
                                if (normalizedTargetValue === "") {
                                    continue;
                                }

                                let n = nodeCache[normalizedTargetValue];
                                if (n) {
                                    if (n.items.indexOf(d.id) == -1){
                                        n.items.push(d.id);
                                    }

                                    if (n.fields.indexOf(target.path) == -1){
                                        n.fields.push(target.path);
                                    }

                                    // should add counter instead of thousands same query being added
                                    n.queries.push(action.items.query);
                                } else {
                                    let n = {
                                        id: normalizedTargetValue,
                                        queries: [action.items.query],
                                        items: [d.id],
                                        name: normalizedTargetValue,
                                        description: '',
                                        icon: [target.icon],
                                        fields: [target.path],
                                    };

                                    nodeCache[n.id] = n;
                                    nodes.push(n);
                                }

                                if (sourceValue.length > 1) {
                                    // we don't want all individual arrays to be linked together
                                    // those individual arrays being linked are (I assume) irrelevant
                                    // otherwise this needs to be a configuration option
                                    continue;
                                }

                                if (linkCache[normalizedSourceValue + normalizedTargetValue]) {
                                    // link already exists
                                    continue;
                                }

                                const link = {
                                    source: normalizedSourceValue,
                                    target: normalizedTargetValue,
                                    color: '#ccc'
                                };

                                links.push(link);
                                linkCache[link.source + link.target] = link;
                            }
                        });
                    }
                });
            });

            return Object.assign({}, state, {
                errors: null,
                nodes: nodes,
                links: links,
                items: concat(state.items, items),
                searches: searches,
                isFetching: false,
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
