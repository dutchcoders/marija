import { INDICES_RECEIVE, INDICES_REQUEST, INDEX_ACTIVATED, INDEX_DEACTIVATED } from './index'
import {getFields} from "../fields/actions";

export function requestIndices(server) {
    return {
        type: INDICES_REQUEST,
        payload: {
            server
        }
    }
}

export function receiveIndices(response) {
    return {
        type: INDICES_RECEIVE,
        payload: {
            ...response
        }
    }
}

function indexActivated(index) {
    return {
        type: INDEX_ACTIVATED,
        payload: {
            index: index
        }
    };
}

export function activateIndex(index) {
    return (dispatch, getState) => {
        dispatch(indexActivated(index));

        // Get the updated datasources
        const datasources = getState().indices.activeIndices;

        // Update the fields based on the new datasources
        dispatch(getFields(datasources));
    };
}

function indexDeActivated(index) {
    return {
        type: INDEX_DEACTIVATED,
        payload: {
            index: index
        }
    };
}

export function deActivateIndex(index) {
    return (dispatch, getState) => {
        dispatch(indexDeActivated(index));

        // Get the updated datasources
        const datasources = getState().indices.activeIndices;

        // Update the fields based on the new datasources
        dispatch(getFields(datasources));
    };
}