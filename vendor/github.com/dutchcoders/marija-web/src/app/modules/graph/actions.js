import { NODES_DESELECT, NODES_DELETE, NODES_HIGHLIGHT, NODE_UPDATE, NODE_SELECT, NODES_SELECT, SELECTION_CLEAR } from './index';

export function deselectNodes(opts) {
    return {
        type: NODES_DESELECT,
        receivedAt: Date.now(),
        nodes: opts
    };
}

export function deleteNodes(opts) {
    return {
        type: NODES_DELETE,
        receivedAt: Date.now(),
        nodes: opts
    };
}


export function highlightNodes(opts) {
    return {
        type: NODES_HIGHLIGHT,
        receivedAt: Date.now(),
        highlight_nodes: opts
    };
}


export function clearSelection(opts) {
    return {
        type: SELECTION_CLEAR,
        receivedAt: Date.now(),
        ...opts,
    };
}

export function nodeUpdate(node_id, params) {
    return {
        type: NODE_UPDATE,
        receivedAt: Date.now(),
        node_id: node_id,
        params: params
    };
}


export function nodeSelect(opts) {
    return {
        type: NODE_SELECT,
        receivedAt: Date.now(),
        ...opts
    };
}

export function nodesSelect(opts) {
    return {
        type: NODES_SELECT,
        receivedAt: Date.now(),
        nodes: opts
    };
}
