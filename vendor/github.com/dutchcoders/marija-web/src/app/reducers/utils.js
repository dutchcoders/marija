import { OPEN_PANE, CLOSE_PANE, CANCEL_REQUEST, Socket } from '../utils/index';

function setPaneTo(panes, pane, state) {
    return panes.map((item) => {
        if (item.name == pane) {
            return Object.assign({}, item, {state: state});
        }
        return item;
    });
}

export default function utils(state = {panes: []}, action) {

    switch (action.type) {
        case OPEN_PANE:
            return Object.assign({}, state, {panes: setPaneTo(state.panes, action.pane, true)});
        case CLOSE_PANE:
            return Object.assign({}, state, {panes: setPaneTo(state.panes, action.pane, false)});
        case CANCEL_REQUEST:
            Socket.ws.postMessage(
                {
                    'request-id': action.requestId
                },
                CANCEL_REQUEST
            );
            return state;
        default:
            return state;
    }
}
