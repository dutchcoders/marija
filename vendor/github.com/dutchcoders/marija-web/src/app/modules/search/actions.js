import { ITEMS_RECEIVE, ITEMS_REQUEST, SEARCH_DELETE } from './index';

const defaultOpts = {
    from: 0,
    size: 500,
    index: "",
    query: "",
    color: ""
};

export function requestItems(opts = defaultOpts) {
    return {
        type: ITEMS_REQUEST,
        receivedAt: Date.now(),
        ...opts
    };
}

export function receiveItems(items, opts = {from: 0}) {
    return {
        type: ITEMS_RECEIVE,
        items: items,
        receivedAt: Date.now()
    };
}

export function deleteSearch(opts) {
    return {
        type: SEARCH_DELETE,
        receivedAt: Date.now(),
        ...opts
    };
}
