import { TABLE_COLUMN_ADD, TABLE_COLUMN_REMOVE, INDEX_ADD, INDEX_DELETE, FIELD_ADD, FIELD_DELETE, DATE_FIELD_ADD, DATE_FIELD_DELETE, NORMALIZATION_ADD, NORMALIZATION_DELETE, INITIAL_STATE_RECEIVE } from './index';

export function tableColumnRemove(field) {
    return {
        type: TABLE_COLUMN_REMOVE,
        receivedAt: Date.now(),
        field: field
    };
}

export function tableColumnAdd(field) {
    return {
        type: TABLE_COLUMN_ADD,
        receivedAt: Date.now(),
        field: field
    };
}

export function indexAdd(index) {
    return {
        type: INDEX_ADD,
        receivedAt: Date.now(),
        index: index
    };
}

export function indexDelete(index) {
    return {
        type: INDEX_DELETE,
        receivedAt: Date.now(),
        index: index
    };
}

export function dateFieldAdd(field) {
    return {
        type: DATE_FIELD_ADD,
        receivedAt: Date.now(),
        field: field
    };
}

export function dateFieldDelete(field) {
    return {
        type: DATE_FIELD_DELETE,
        receivedAt: Date.now(),
        field: field
    };
}


export function fieldAdd(field) {
    return {
        type: FIELD_ADD,
        receivedAt: Date.now(),
        field: field
    };
}

export function fieldDelete(field) {
    return {
        type: FIELD_DELETE,
        receivedAt: Date.now(),
        field: field
    };
}

export function normalizationAdd(normalization) {
    return {
        type: NORMALIZATION_ADD,
        receivedAt: Date.now(),
        normalization: normalization
    };
}

export function normalizationDelete(normalization) {
    return {
        type: NORMALIZATION_DELETE,
        receivedAt: Date.now(),
        normalization: normalization
    };
}

export function receiveInitialState(initial_state) {
    return {
        type: INITIAL_STATE_RECEIVE,
        receivedAt: Date.now(),
        initial_state: initial_state
    };
}

