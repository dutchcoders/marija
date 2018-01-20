import {authConnected, error} from '../index';
import Websocket from 'reconnecting-websocket';

import { SEARCH_REQUEST } from '../../modules/search/index'

export default class FlowWS {
    constructor(url, token, dispatcher, storeDispatcher) {
        this.url = url;

        const websocket = new Websocket(url);

        this.opened = new Promise(resolve => {
            websocket.onopen = function (event) {
                storeDispatcher(authConnected({connected: true, errors: null}));
                resolve();
            };
        });

        websocket.onclose = function (event) {
            var reason = "";

            if (event.code == 1000)
                reason = "Normal closure, meaning that the purpose for which the connection was established has been fulfilled.";
            else if (event.code == 1001)
                reason = "An endpoint is \"going away\", such as a server going down or a browser having navigated away from a page.";
            else if (event.code == 1002)
                reason = "An endpoint is terminating the connection due to a protocol error";
            else if (event.code == 1003)
                reason = "An endpoint is terminating the connection because it has received a type of data it cannot accept (e.g., an endpoint that understands only text data MAY send this if it receives a binary message).";
            else if (event.code == 1004)
                reason = "Reserved. The specific meaning might be defined in the future.";
            else if (event.code == 1005)
                reason = "No status code was actually present.";
            else if (event.code == 1006)
                reason = "The connection was closed abnormally, e.g., without sending or receiving a Close control frame";
            else if (event.code == 1007)
                reason = "An endpoint is terminating the connection because it has received data within a message that was not consistent with the type of the message (e.g., non-UTF-8 [http://tools.ietf.org/html/rfc3629] data within a text message).";
            else if (event.code == 1008)
                reason = "An endpoint is terminating the connection because it has received a message that \"violates its policy\". This reason is given either if there is no other sutible reason, or if there is a need to hide specific details about the policy.";
            else if (event.code == 1009)
                reason = "An endpoint is terminating the connection because it has received a message that is too big for it to process.";
            else if (event.code == 1010) // Note that this status code is not used by the server, because it can fail the WebSocket handshake instead.
                reason = "An endpoint (client) is terminating the connection because it has expected the server to negotiate one or more extension, but the server didn't return them in the response message of the WebSocket handshake. <br /> Specifically, the extensions that are needed are: " + event.reason;
            else if (event.code == 1011)
                reason = "A server is terminating the connection because it encountered an unexpected condition that prevented it from fulfilling the request.";
            else if (event.code == 1015)
                reason = "The connection was closed due to a failure to perform a TLS handshake (e.g., the server certificate can't be verified).";
            else
                reason = "Unknown reason";

            // todo(nl5887): differentiate connection errors and query errors
            storeDispatcher(authConnected({connected: false, errors: reason}));
        };
        websocket.onerror = function (event) {
            // the close event contains more feedback
        };
        websocket.onmessage = function (event) {
            const parsed = JSON.parse(event.data);
            console.log('receive', parsed.type);

            if (parsed.type === 'ERROR') {
                console.error('received error over socket: ', parsed);
            }

            dispatcher(parsed, storeDispatcher);
        };

        this.websocket = websocket;
    }

    postMessage(data, type = SEARCH_REQUEST) {
        console.log('send ' + type, data);

        this.opened.then(() => {
            this.websocket.send(
                JSON.stringify({
                    type: type,
                    ...data
                })
            );
        });
    }

    close() {
        this.websocket.close();
    }
}
