import filterSecondaryComponents from "./filterSecondaryComponents";

test('should filter out components that dont contain nodes from the primary query', () => {
    const components = [
        [
            {
                id: 'a',
                queries: ['first search']
            },
            {
                id: 'b',
                queries: ['second search']
            }
        ],
        [
            {
                id: 'c',
                queries: ['second search']
            },
            {
                id: 'd',
                queries: ['second search']
            }
        ]
    ];

    const filtered = filterSecondaryComponents('first search', components);

    expect(filtered.length).toBe(1);
    expect(filtered[0].length).toBe(2);
    expect(filtered[0].find(node => node.queries.indexOf('first search'))).toBeDefined();
});