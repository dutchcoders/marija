import { FIELDS_RECEIVE, FIELDS_REQUEST, FIELDS_CLEAR } from '../modules/fields/index'

import { concat, without, map, filter, union, reduce, merge } from 'lodash'
import { Fields } from '../domain/index'

import { Socket } from '../utils/index';

const defaultState = {
    availableFields: [],
    fieldsFetching: false
};

export default function fields(state = defaultState, action) {
    switch (action.type) {
        case FIELDS_RECEIVE:
            const newAvailableFields = Fields.extractNewFields(action.payload.fields, state.availableFields);

            return Object.assign({}, state, {
                availableFields: state.availableFields.concat(newAvailableFields),
                fieldsFetching: false
            });

        case FIELDS_REQUEST:
            if (action.payload.indexes.length === 0) {
                return Object.assign({}, state, { availableFields: defaultState.availableFields });
            }

            Socket.ws.postMessage(
                {
                    datasources: action.payload.indexes
                },
                FIELDS_REQUEST
            );

            return Object.assign({}, state, {
                fieldsFetching: true
            });

        case FIELDS_CLEAR:
            return Object.assign({}, state, { availableFields: defaultState.availableFields });
            break;

        default:
            return state;
    }
}
