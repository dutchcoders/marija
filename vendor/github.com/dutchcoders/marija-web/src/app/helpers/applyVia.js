import { concat } from 'lodash';

function getConnectedNodes(node, nodes, links) {
    const connected = [];

    links.forEach(link => {
        if (link.target === node.id) {
            connected.push(nodes.find(search => search.id === link.source));
        } else if (link.source === node.id) {
            connected.push(nodes.find(search => search.id === link.target));
        }
    });

    return connected;
}

function getLinkLabel(label) {
    label += '';

    const maxLength = 20;

    if (label.length <= maxLength) {
        return label;
    }

    const shortened = label.substring(0, maxLength - 1);

    return shortened + '...';
}

export default function applyVia(nodes, links, via) {
    links = concat([], links);
    nodes = concat([], nodes);

    const removeNode = (node) => {
        links = links.filter(link =>
            link.source !== node.id && link.target !== node.id
        );
        nodes = nodes.filter(search => search.id !== node.id);
    };

    const removeUnlabeledLinks = (sourceId, targetId) => {
        links = links.filter(link => {
            const remove =
                (link.source === sourceId && link.target === targetId)
                || (link.source === targetId && link.target === sourceId);

            return !remove;
        });
    };

    const nodesToRemove = [];
    const linksToRemove = [];
    const newLinks = [];

    via.forEach(viaItem => {
        const step1Nodes = nodes.filter(node =>
            node.fields.length === 1
            && node.fields[0] === viaItem.endpoints[0]
        );

        step1Nodes.forEach(step1Node => {
            const connected = getConnectedNodes(step1Node, nodes, links);

            const step2Nodes = connected.filter(search =>
                search.fields.length === 1
                && search.fields[0] === viaItem.label
            );

            step2Nodes.forEach(step2Node => {
                const connected = getConnectedNodes(step2Node, nodes, links);

                const step3Nodes = connected.filter(search =>
                    search.fields.length === 1
                    && search.id !== step1Node.id
                    && search.fields[0] === viaItem.endpoints[1]
                );

                step3Nodes.forEach(step3Node => {
                    // Now we found something that can be replaced by a labeled link
                    // but first check if we didn't already do that, but then
                    // the other way around
                    const existing = newLinks.find(link =>
                        link.target === step1Node.id && link.source === step3Node.id
                    );

                    if (typeof existing === 'undefined') {
                        newLinks.push({
                            source: step1Node.id,
                            target: step3Node.id,
                            label: getLinkLabel(step2Node.name)
                        });
                    }

                    nodesToRemove.push(step2Node);
                    linksToRemove.push({
                        from: step1Node,
                        to: step3Node
                    });
                });
            });
        });
    });

    nodesToRemove.forEach(node => removeNode(node));
    linksToRemove.forEach(data => {
        removeUnlabeledLinks(data.from.id, data.to.id);
    });

    const counter = {};

    newLinks.forEach(link => {
        const key = link.source + link.target;
        const current = counter[key] ? counter[key] + 1 : 1;

        link.current = current;
        counter[key] = current;

        links.push(link);
    });

    newLinks.forEach(link => {
        link.total = counter[link.source + link.target];
    });

    return {
        nodes,
        links
    };
}