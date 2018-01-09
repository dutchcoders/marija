import {fieldLocator} from './index';
import {forEach} from 'lodash';

export default function getHighlightItem (item, node, fields) {
    const highlightItem = {
        id: item.id,
        fields: {}
    };

    // Only keep the fields that the user configured for brevity
    forEach(fields, (field) => {
        const value = fieldLocator(item.fields, field.path);

        if (value !== null) {
            highlightItem.fields[field.path] = value;
        }
    });

    highlightItem.x = node.x;
    highlightItem.y = node.y;
    highlightItem.matchFields = node.fields;

    return highlightItem;
}