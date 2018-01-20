import { OPEN_PANE, CLOSE_PANE, CANCEL_REQUEST, Socket, HEADER_HEIGHT_CHANGE } from '../utils/index';

function setPaneTo(panes, pane, state) {
    return panes.map((item) => {
        if (item.name == pane) {
            return Object.assign({}, item, {state: state});
        }
        return item;
    });
}

const defaultState = {
    panes: [],
    headerHeight: 56
};

export default function utils(state = defaultState, action) {

    switch (action.type) {
        case OPEN_PANE:
            return Object.assign({}, state, {panes: setPaneTo(state.panes, action.pane, true)});
        case CLOSE_PANE:
            return Object.assign({}, state, {panes: setPaneTo(state.panes, action.pane, false)});
        case HEADER_HEIGHT_CHANGE:
            return Object.assign({}, state, {headerHeight: action.height});
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
