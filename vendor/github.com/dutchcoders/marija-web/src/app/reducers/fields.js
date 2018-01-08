import { FIELDS_RECEIVE, FIELDS_REQUEST, FIELDS_CLEAR } from '../modules/fields/index'

import { concat, without, map, filter, union, reduce, merge } from 'lodash'
import { Fields } from '../domain/index'

import { Socket } from '../utils/index';

const defaultState = {
    availableFields: []
};

export default function fields(state = defaultState, action) {
    switch (action.type) {
        case FIELDS_RECEIVE:
            const newAvailableFields = Fields.extractNewFields(action.payload.fields, state.availableFields);
            return Object.assign({}, state, {availableFields: state.availableFields.concat(newAvailableFields)});
        case FIELDS_REQUEST:
            Socket.ws.postMessage(
                {
                    datasources: action.payload.indexes
                },
                FIELDS_REQUEST
            );
            return state;

        case FIELDS_CLEAR:
            return Object.assign({}, state, { availableFields: defaultState.availableFields });
            break;
        default:
            return state;
    }
}
