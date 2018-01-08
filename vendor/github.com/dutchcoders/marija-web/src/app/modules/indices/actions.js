import { INDICES_RECEIVE, INDICES_REQUEST, INDEX_ACTIVATED, INDEX_DEACTIVATED } from './index'

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

export function activateIndex(index) {
    return {
        type: INDEX_ACTIVATED,
        payload: {
            index: index
        }
    }
}
export function deActivateIndex(index) {
    return {
        type: INDEX_DEACTIVATED,
        payload: {
            index: index
        }
    }
}