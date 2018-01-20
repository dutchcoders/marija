export default function removeDeadLinks(nodes, links) {
    return links.filter(link => {
        const source = nodes.find(node => node.id === link.source);
        const target = nodes.find(node => node.id === link.target);

        return typeof source !== 'undefined' && typeof target !== 'undefined';
    });
}