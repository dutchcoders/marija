import {ERROR, AUTH_CONNECTED, OPEN_PANE, CLOSE_PANE, CANCEL_REQUEST, HEADER_HEIGHT_CHANGE} from './index';

export function error(msg) {
    return {
        type: ERROR,
        receivedAt: Date.now(),
        errors: msg
    };
}

export function cancelRequest(requestId) {
    return {
        type: CANCEL_REQUEST,
        requestId: requestId
    };
}

export function authConnected(p) {
    return {
        type: AUTH_CONNECTED,
        receivedAt: Date.now(),
        ...p
    };
}

export function openPane(pane) {
    return {
        type: OPEN_PANE,
        pane: pane
    };
}

export function closePane(pane) {
    return {
        type: CLOSE_PANE,
        pane: pane
    };
}

export function headerHeightChange(height) {
    return {
        type: HEADER_HEIGHT_CHANGE,
        height: height
    };
}
