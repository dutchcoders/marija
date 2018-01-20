import getConnectedComponents from "./getConnectedComponents";

test('should identify the separate components that make up the graph', () => {
    const nodes = [
        {
            id: 'a'
        },
        {
            id: 'b'
        },
        {
            id: 'c'
        },
        {
            id: 'd'
        },
    ];

    const links = [
        {
            source: 'a',
            target: 'b'
        },
        {
            source: 'c',
            target: 'd'
        }
    ];

    const components = getConnectedComponents(nodes, links);

    // Expect to get 2 components, with both 2 nodes
    expect(components.length).toBe(2);
    expect(components[0].length).toBe(2);
    expect(components[1].length).toBe(2);
});