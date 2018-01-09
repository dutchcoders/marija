import {uniqueId} from 'lodash';
import getHighlightItem from "./getHighlightItem";

const generateItemAndNode = (fields) => {
    const id = uniqueId();

    const item = {
        fields: fields,
        highlight: null,
        id: id,
        matchFields: Object.keys(fields),
        query: 'wilders',
        x: 0,
        y: 0
    };

    const node = {
        description: '',
        fields: Object.keys(fields),
        force: 1,
        icon: ['a'],
        id: 'yolo',
        index: 1,
        items: [id],
        name: 'yolo',
        queries: ['wilders'],
        r: 15,
        vx: 1,
        vy: 1,
        x: 1,
        y: 1,
    };

    return {
        item,
        node
    };
};

const generateFields = (fields) => {
    const ret = [];

    fields.forEach(field => {
        ret.push({
            path: field,
            icon: 'a'
        });
    });

    return ret;
};

test('highlight item should contain fields of item', () => {
    const fields = generateFields(['text', 'user']);

    const {item, node} = generateItemAndNode({
        text: 'yolo',
        user: 'thomas'
    });

    const highlighItem = getHighlightItem(item, node, fields);

    expect(highlighItem.fields).toEqual({
        text: 'yolo',
        user: 'thomas'
    });
});

test('should delete fields that are not configured', () => {
    const fields = generateFields(['text']);

    const {item, node} = generateItemAndNode({
        text: 'yolo',
        user: 'thomas'
    });

    const highlighItem = getHighlightItem(item, node, fields);

    expect(highlighItem.fields).toEqual({
        text: 'yolo'
    });
});

test('should not contain fields that were not in the item', () => {
    const fields = generateFields(['text', 'unconfiguredField']);

    const {item, node} = generateItemAndNode({
        text: 'yolo',
        user: 'thomas'
    });

    const highlighItem = getHighlightItem(item, node, fields);

    expect(highlighItem.fields).toEqual({
        text: 'yolo'
    });
});

test('should keep fields in nested objects, like user.name', () => {
    const fields = generateFields(['text', 'user.name']);

    const {item, node} = generateItemAndNode({
        text: 'yolo',
        user: {
            name: 'thomas'
        }
    });

    const highlighItem = getHighlightItem(item, node, fields);

    expect(highlighItem.fields).toEqual({
        'text': 'yolo',
        'user.name': 'thomas'
    });
});