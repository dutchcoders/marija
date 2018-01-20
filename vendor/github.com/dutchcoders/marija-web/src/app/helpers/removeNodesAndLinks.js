import { slice, concat, without, pull } from 'lodash';
import {fieldLocator, normalize} from "./index";

const removeNodes = (nodes, removedQuery) => {
    nodes.forEach((node, index) => {
        // Remove this query from the list of queries that the node appeared for
        node.queries = without(node.queries, removedQuery);

        // When there are no more queries where this node appeared for, we can remove the node
        if (node.queries.length === 0) {
            nodes = without(nodes, node);
        }
    });

    return nodes;
};

export default function removeNodesAndLinks(previousNodes, previousLinks, removedQuery) {
    let nodes = concat(previousNodes, []);
    let links = concat(previousLinks, []);

    nodes = removeNodes(nodes, removedQuery);

    links.forEach(link => {
        const sourceNode = nodes.find(node => node.name === link.source);
        const targetNode = nodes.find(node => node.name === link.target);

        // Remove the link when either the source or the target node no longer exists
        if (!sourceNode || !targetNode) {
            links = without(links, link);
        }
    });

    return {
        nodes,
        links
    };
}