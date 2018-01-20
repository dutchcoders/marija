import { find } from 'lodash';

export default function getRelatedNodes(nodes, allNodes, allLinks) {
    const related_nodes = [];

    let x = (n) => {
        related_nodes.push(n);

        for (let link of allLinks) {
            if (link.source === n.id) {
                // check if already visited
                if (find(related_nodes, (o) => {
                        return (link.target === o.id);
                    })) {
                    continue;
                }

                const target_node = find(allNodes, (n2) => {
                    return (link.target === n2.id);
                });

                x(target_node);
            }

            if (link.target === n.id) {
                if (find(related_nodes, (o) => {
                        return (link.source === o.id);
                    })) {
                    continue;
                }

                const source_node = find(allNodes, (n2) => {
                    return (link.source === n2.id);
                });

                x(source_node);
            }
        }
    };

    for (let n of nodes) {
        x(n);
    }

    return related_nodes;
}