import getNodesAndLinks from './getNodesAndLinks';
import { uniqueId } from 'lodash';
import applyVia from "./applyVia";

const generateItem = (fields) => {
    if (typeof fields === 'undefined') {
        fields = {
            text: 'test' + uniqueId()
        };
    }

    return {
        highlight: null,
        id: uniqueId(),
        fields: fields,
        query: undefined
    };
};

const generateField = (field) => {
    return {
        icon: 'a',
        path: field
    };
};

const generateQuery = (items) => {
    return {
        color: '#aaaaaa',
        q: 'my search',
        total: 100,
        items: items
    };
};

// test if a link exists between a source and a target
const expectLink = (links, source, target) => {
    const link = links.find(link =>
        (link.source === source && link.target === target)
        || (link.target === source && link.source === target)
    );
    expect(link).toBeDefined();
};

test('should output nodes', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        generateItem(),
        generateItem(),
        generateItem()
    ];

    const fields = [
        generateField('text')
    ];

    const query = generateQuery(items);
    const { nodes } = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations);

    expect(nodes.length).toBe(3);
});

test('should output 1 node for every field in an item', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        generateItem({
            text: 'hello',
            user: 'thomas'
        }),
    ];

    const fields = [
        generateField('text'),
        generateField('user')
    ];

    const query = generateQuery(items);
    const { nodes } = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations);

    expect(nodes.length).toBe(2);
});

test('should output links between related nodes', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        generateItem({
            text: 'lalala',
            user: 'thomas'
        }),
    ];

    const fields = [
        generateField('text'),
        generateField('user')
    ];

    const query = generateQuery(items);
    const { nodes, links } = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations);

    expect(links).toBeDefined();
    expect(links.length).toBe(1);
});


test('when nodes have exactly the same fields they should not be duplicated', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        generateItem({text: 'same'}),
        generateItem({text: 'same'})
    ];

    const fields = [
        generateField('text'),
    ];

    const query = generateQuery(items);
    const { nodes } = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations);

    expect(nodes.length).toBe(1);
});

test('should output links with labels when via info is specified', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        {
            id: uniqueId(),
            fields: {
                'port': 80,
                'client': 1,
                'server': 2
            }
        }
    ];

    const fields = [
        generateField('client'),
        generateField('server'),
        generateField('port'),
    ];

    const query = generateQuery(items);

    const via = [
        {
            endpoints: ['client', 'server'],
            label: 'port'
        }
    ];

    const result = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations);
    const { nodes, links } = applyVia(result.nodes, result.links, via);

    expect(nodes.length).toBe(2);

    expect(links).toBeDefined();
    expect(links.length).toBe(1);

    expectLink(links, 1, 2);

    links.forEach(link => {
        expect(link.label).toBe('80');
    });
});

test('should be able to draw multiple labeled lines between 2 nodes', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        {
            id: uniqueId(),
            fields: {
                'port': 80,
                'client': 1,
                'server': 2
            }
        },
        {
            id: uniqueId(),
            fields: {
                'port': 1337,
                'client': 1,
                'server': 2
            }
        }
    ];

    const fields = [
        generateField('client'),
        generateField('server'),
        generateField('port'),
    ];

    const query = generateQuery(items);

    const via = [{
        endpoints: ['client', 'server'],
        label: 'port'
    }];

    const result = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations, via);
    const { nodes, links } = applyVia(result.nodes, result.links, via);

    expect(nodes.length).toBe(2);
    expect(links).toBeDefined();
    expect(links.length).toBe(2);
    expect(links.find(link => link.label === '80')).toBeDefined();
    expect(links.find(link => link.label === '1337')).toBeDefined();
});

test('should not mess up when multiple via configs are present', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        {
            id: uniqueId(),
            fields: {
                'port': 80,
                'source': 1,
                'target': 2,
                'source2': 10,
                'target2': 20
            }
        },
        {
            id: uniqueId(),
            fields: {
                'port': 1337,
                'source': 1,
                'target': 2
            }
        }
    ];

    const fields = [
        generateField('source'),
        generateField('target'),
        generateField('source2'),
        generateField('target2'),
        generateField('port')
    ];

    const query = generateQuery(items);

    const via = [
        {
            endpoints: ['source', 'target'],
            label: 'port'
        },
        {
            endpoints: ['source2', 'target2'],
            label: 'port'
        }
    ];

    const result = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations, via);
    const { nodes, links } = applyVia(result.nodes, result.links, via);

    expect(nodes.length).toBe(4);
    expect(links).toBeDefined();
    expect(links.length).toBe(7);
});

test('should generate labeled links between endpoints of the same type', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        {
            id: uniqueId(),
            fields: {
                'ip': 1,
                'port': 80,
            }
        },
        {
            id: uniqueId(),
            fields: {
                'ip': 2,
                'port': 80
            }
        }
    ];

    const fields = [
        generateField('ip'),
        generateField('port')
    ];

    const query = generateQuery(items);

    const via = [
        {
            endpoints: ['ip', 'ip'],
            label: 'port'
        }
    ];

    let { nodes, links } = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations, via);
    const viaResult = applyVia(nodes, links, via);
    nodes = viaResult.nodes;
    links = viaResult.links;

    expect(nodes.length).toBe(2);
    expect(links).toBeDefined();
    expect(links.length).toBe(1);
    expect(links[0].label).toBe('80');
});

test('should not remove too many links when via info is specified', () => {
    const previousNodes = [];
    const previousLinks = [];
    const normalizations = [];

    const items = [
        {
            id: uniqueId(),
            fields: {
                'port': 80,
                'client': 1,
                'server': 2,
                'country': 'nl'
            }
        }
    ];

    const fields = [
        generateField('client'),
        generateField('server'),
        generateField('port'),
        generateField('country'),
    ];

    const query = generateQuery(items);

    const via = [
        {
            endpoints: ['client', 'server'],
            label: 'port'
        }
    ];

    const result = getNodesAndLinks(previousNodes, previousLinks, items, fields, query, normalizations);

    const { nodes, links } = applyVia(result.nodes, result.links, via);

    expect(nodes.length).toBe(3);

    expect(links).toBeDefined();
    expect(links.length).toBe(3);

    expectLink(links, 1, 2);
    expectLink(links, 'nl', 1);
    expectLink(links, 'nl', 2);
});