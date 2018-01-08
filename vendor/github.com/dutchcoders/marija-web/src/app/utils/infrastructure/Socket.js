import { FlowWS, error } from '../../utils/index';
import { receiveItems, ITEMS_RECEIVE } from '../../modules/search/index';
import { receiveIndices, INDICES_RECEIVE} from '../../modules/indices/index';
import { receiveFields, FIELDS_RECEIVE} from '../../modules/fields/index';
import { receiveInitialState, INITIAL_STATE_RECEIVE } from '../../modules/data/index';

export const Socket = {
    ws: null,
    wsDispatcher: (message, dispatch) => {
        if (message.error) {
            return dispatch(error(message.error.message));
        }

        switch (message.type) {
            case ITEMS_RECEIVE:
		dispatch(receiveItems(message.items));
                break;

            case INDICES_RECEIVE:
		dispatch(receiveIndices(message.indices));
                break;

            case FIELDS_RECEIVE:
		dispatch(receiveFields(message.fields));
                break;

            case INITIAL_STATE_RECEIVE:
		dispatch(receiveInitialState(message.state));
                break;
        }
    },
    startWS: (dispatch) => {
        if (!!Socket.ws) {
            return;
        }
        
        const { location }  = window;
        const url = ((location.protocol === "https:") ? "wss://" : "ws://") + location.host + "/ws";

        try {
            Socket.ws = new FlowWS(url, null, Socket.wsDispatcher, dispatch);
        } catch (e) {
            dispatch(error(e));
        }
    }
};
