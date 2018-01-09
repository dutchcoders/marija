import { entries, enableBatching, utils, servers, indices, fields, defaultState } from './';
import { dispatch, compose, createStore, combineReducers, applyMiddleware } from 'redux';
import { syncHistoryWithStore, routerReducer } from 'react-router-redux';
import {IMPORT_DATA} from "../modules/import/constants";

const appReducer = enableBatching(
    combineReducers({
        entries,
        utils,
        servers,
        indices,
        fields,
        routing: routerReducer
    })
);

export default function root(state, action) {
    if (action.type === IMPORT_DATA) {
        state = action.payload;
        console.log(state);
    }

    return appReducer(state, action);
}