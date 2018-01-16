import {fieldLocator} from './index';
import {forEach} from 'lodash';

/**
 * If the value exceeds the maximum length, attempt to display the relevant part
 * of the query.
 *
 * @param value
 * @param query
 * @param maxValueLength
 * @returns string
 */
function getValue(value, query, maxValueLength) {
    if (value.length <= maxValueLength) {
        return value;
    }

    const matchIndex = value.toLowerCase().indexOf(query.toLowerCase());

    if (matchIndex === -1) {
        return value.substring(0, maxValueLength - 1) + "...";
    }

    const match = value.substring(matchIndex, matchIndex + query.length);
    const contextLength = Math.round((maxValueLength - query.length) / 2);
    const contextLeft = value.substring(matchIndex - contextLength, matchIndex);
    const contextRight = value.substring(matchIndex + query.length, matchIndex + query.length + contextLength);

    let ret = '';

    if (contextLeft) {
        if (value.indexOf(contextLeft) > 0) {
            ret += '...';
        }

        ret += contextLeft;
    }

    ret += match;

    if (contextRight) {
        ret += contextRight;

        if (value.indexOf(contextRight) + contextRight.length < value.length) {
            ret += '...';
        }
    }

    return ret;
}

export default function getHighlightItem (item, node, fields, maxValueLength) {
    const highlightItem = {
        id: item.id,
        fields: {}
    };

    // Only keep the fields that the user configured for brevity
    forEach(fields, (field) => {
        const value = fieldLocator(item.fields, field.path);

        if (value !== null) {
            highlightItem.fields[field.path] = getValue(value, item.query, maxValueLength);
        }
    });

    highlightItem.x = node.x;
    highlightItem.y = node.y;
    highlightItem.matchFields = node.fields;
    highlightItem.query = item.query;

    return highlightItem;
}