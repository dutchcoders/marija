import { concat, without } from 'lodash'

import { SERVER_ADD, SERVER_REMOVE } from '../modules/servers/index'

const defaultState = [
    "http://127.0.0.1:9200/"
];


export default function servers(state = defaultState, action) {
    switch (action.type) {
        case SERVER_ADD:
            const newServers = concat(state, action.server);
            return newServers;

        case SERVER_REMOVE:
            const filteredServers = without(state, action.server);
            return filteredServers;

        default:
            return state;
    }
}