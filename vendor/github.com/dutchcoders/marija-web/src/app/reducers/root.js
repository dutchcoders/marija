import { entries, enableBatching, utils, servers, indices, fields, defaultState } from './';
import { dispatch, compose, createStore, combineReducers, applyMiddleware } from 'redux';
import { syncHistoryWithStore, routerReducer } from 'react-router-redux';
import {EXPORT_DATA, IMPORT_DATA} from "../modules/import/constants";
import exportJson from "../helpers/exportJson";

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
    } else if (action.type === EXPORT_DATA) {
        exportJson(state);
    }

    return appReducer(state, action);
}