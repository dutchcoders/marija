import removeNodesAndLinks from './removeNodesAndLinks';
import { uniqueId } from 'lodash';

const generateNode = (name, queries) => {
    return {
        id: uniqueId(),
        queries: queries,
        name: 'test' + uniqueId(),
        description: 'vfnjdsvnfds',
        icon: 'a',
        fields: [
            'text'
        ]
    };
};

const generateLink = (source, target) => {
    return {
        color: '#ccc',
        source: source,
        target: target,
        queries: []
    };
};

test('should remove nodes', () => {
    const previousNodes = [
        generateNode('test1', ['test query 1']),
        generateNode('test2', ['test query 2'])
    ];

    const {nodes, links} = removeNodesAndLinks(previousNodes, [], 'test query 1');

    expect(nodes.length).toBe(1);
});

test('should remove links to nodes that no longer exist', () => {
    const previousNodes = [
        generateNode('test1', ['test query 1']),
        generateNode('test2', ['test query 2'])
    ];

    const previousLinks = [
        generateLink('test1', 'test2')
    ];

    const {nodes, links} = removeNodesAndLinks(previousNodes, previousLinks, 'test query 1');

    expect(links.length).toBe(0);
});