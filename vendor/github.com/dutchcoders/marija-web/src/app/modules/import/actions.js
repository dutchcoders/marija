import {IMPORT_DATA, EXPORT_DATA} from "./index";

export function importData(data) {
    return {
        type: IMPORT_DATA,
        payload: data
    };
}

export function exportData() {
    return {
        type: EXPORT_DATA
    };
}