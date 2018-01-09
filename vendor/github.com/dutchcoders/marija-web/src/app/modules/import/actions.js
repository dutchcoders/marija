import {IMPORT_DATA} from "./index";

export function importData(data) {
    return {
        type: IMPORT_DATA,
        payload: data
    };
}