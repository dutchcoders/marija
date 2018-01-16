import {saveAs} from 'file-saver';

export default function exportJson(data) {
    const blob = new Blob(
        [JSON.stringify(data)],
        {type: "text/json;charset=utf-8"}
    );

    const now = new Date();
    const dateString = now.getFullYear() + '-' + now.getMonth() + '-' + now.getDate();
    const filename = 'marija-export-' + dateString + '.json';

    saveAs(blob, filename);
}