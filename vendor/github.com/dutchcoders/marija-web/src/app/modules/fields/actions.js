import { FIELDS_RECEIVE, FIELDS_REQUEST, FIELDS_CLEAR } from './index';


export function clearAllFields(){
    return {
        type: FIELDS_CLEAR,
    };
}


export function receiveFields(fields) {
    return {
        type: FIELDS_RECEIVE,
        payload: {
            fields: fields
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
