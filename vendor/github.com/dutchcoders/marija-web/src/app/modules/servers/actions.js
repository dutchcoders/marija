import { SERVER_ADD, SERVER_REMOVE } from './index'


export function serverAdd(server) {
    return {
        type: SERVER_ADD,
        receivedAt: Date.now(),
        server: server
    };
}

export function serverRemove(server) {
    return {
        type: SERVER_REMOVE,
        receivedAt: Date.now(),
        server: server
    };
}
