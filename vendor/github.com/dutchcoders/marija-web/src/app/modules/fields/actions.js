import { FIELDS_RECEIVE, FIELDS_REQUEST, FIELDS_CLEAR } from './index';


export function clearAllFields(){
    return {
        type: FIELDS_CLEAR,
    };
}


export function receiveFields(payload) {
    return {
        type: FIELDS_RECEIVE,
        payload: {
            ...payload
        }
    };
}

export function getFields(indexes) {
    return {
        type: FIELDS_REQUEST,
        payload: {
            indexes: indexes
        }
    };
}
