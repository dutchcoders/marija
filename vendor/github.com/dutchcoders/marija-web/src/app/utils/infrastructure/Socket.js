import {FlowWS, error} from '../../utils/index';
import {searchReceive, SEARCH_RECEIVE, SEARCH_COMPLETED} from '../../modules/search/index';
import {receiveIndices, INDICES_RECEIVE} from '../../modules/indices/index';
import {receiveFields, FIELDS_RECEIVE} from '../../modules/fields/index';
import {
    receiveInitialState,
    INITIAL_STATE_RECEIVE
} from '../../modules/data/index';
import {searchCompleted} from "../../modules/search/actions";

export const Socket = {
    ws: null,
    wsDispatcher: (message, dispatch) => {
        if (message.error) {
            return dispatch(error(message.error.message));
        }

        switch (message.type) {
            case SEARCH_RECEIVE:
                dispatch(searchReceive({
                    results: message.results,
                    query: message.query,
                }));
                break;

            case INDICES_RECEIVE:
                dispatch(receiveIndices(message.indices));
                break;

            case FIELDS_RECEIVE:
                dispatch(receiveFields(message.fields));
                break;

            case INITIAL_STATE_RECEIVE:
                dispatch(receiveInitialState({
                    datasources: message.datasources,
                    version: message.version
                }));
                break;

            case SEARCH_COMPLETED:
                dispatch(searchCompleted(message['request-id']));
        }
    },
    startWS: (dispatch) => {
        if (!!Socket.ws) {
            return;
        }

        let url;

        if (process.env.WEBSOCKET_URI) {
            url = process.env.WEBSOCKET_URI;
        } else {
            const {location} = window;
            url = ((location.protocol === "https:") ? "wss://" : "ws://") + location.host + "/ws";
        }

        try {
            Socket.ws = new FlowWS(url, null, Socket.wsDispatcher, dispatch);
        } catch (e) {
            dispatch(error(e));
        }
    }
};
