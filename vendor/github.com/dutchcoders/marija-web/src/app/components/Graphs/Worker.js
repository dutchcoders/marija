import * as d3 from 'd3';
import { map, clone, groupBy, reduce, forEach, difference, find, uniq, remove, each, includes, assign, isEqual } from 'lodash';

let simulation = null;
let timer = null;
let nodes = [];
let links = [];

onmessage = function(event) {
    if (event.data.type === "restart") {
        let { nodes } = event.data;

        for (let n1 of this.nodes) {
            for (let n2 of nodes) {
                if (n1.id !== n2.id)
                    continue;

                n1.fx = n2.fx;
                n1.fy = n2.fy;
            }
        }
        
        simulation.alpha(0.3).restart();
    } else if (event.data.type === 'stop') {
        let { nodes } = event.data;

        for (let n1 of this.nodes) {
            for (let n2 of nodes) {
                if (n1.id !== n2.id)
                    continue;

                n1.fx = n2.fx;
                n1.fy = n2.fy;
            }
        }
        
        // simulation.alpha(0);
    } else if (event.data.type === 'init') {
        let { clientWidth, clientHeight } = event.data;

        simulation = d3.forceSimulation()
            .stop()
            .force("link", d3.forceLink().distance(link => {
                if (!link.label) {
                    return 80;
                }

                let label = link.label;

                if (Array.isArray(label)) {
                    label = label.join('');
                }

                return label.length * 10 + 30;
            }).id(function (d) {
                return d.id;
            }))
            .force("charge", d3.forceManyBody().strength(-100).distanceMax(500))
            .force("center", d3.forceCenter(clientWidth / 2, clientHeight / 2))
            .force("vertical", d3.forceY().strength(0.018))
            .force("horizontal", d3.forceX().strength(0.006))
            .on("tick", () => {
                postMessage({type: "tick", nodes: this.nodes, links: this.links });
            });

        this.nodes = [];
        this.links = [];
    } else if (event.data.type === 'tick') {
    } else if (event.data.type === 'update') {
        let { nodes, links } = event.data;

        const sizeRange = [15, 30];

        let forceScale = function (node) {
            var scale = d3.scaleLog().domain(sizeRange).range(sizeRange.slice().reverse());
            return node.r + scale(node.r);
        };

        var countExtent = d3.extent(nodes, (n) => {
            return n.count;
        }),
            radiusScale = d3.scalePow().exponent(2).domain(countExtent).range(sizeRange);

        var newNodes = false;

        var that = this;

        // remove deleted nodes
        remove(this.nodes, (n) => {
            return !find(nodes, (o) => {
                return (o.id==n.id);
            });
        });

        for (let i=0; i < nodes.length; i++) {
            let node = nodes[i];
            // todo(nl5887): cleanup

            var n = find(that.nodes, {id: node.id});
            if (n) {
                n = assign(n, node);
                n = assign(n, {force: forceScale(n), r: radiusScale(n.count)});

                newNodes = true;
                continue;
            }

            let node2 = clone(node);
            node2 = assign(node2, {force: forceScale(node2), r: radiusScale(node2.count)});

            that.nodes.push(node2);

            newNodes = true;
        }

        remove(this.links, (link) => {
            return !find(links, (o) => {
                return (link.source.id == o.source && link.target.id == o.target);
            });
        });

        for (let i=0; i < links.length; i++) {
            let link = links[i];

            var n = find(that.links, (o) => {
                return o.source.id == link.source && o.target.id == link.target;
            });
            
            if (n) {
                link.color = n.color;
                continue;
            }
            
            // todo(nl5887): why?
            that.links.push({
                source: link.source,
                target: link.target,
                color: link.color,
                label: link.label,
                total: link.total,
                current: link.current
            });
        }

        simulation
            .nodes(this.nodes);

        simulation.force("link")
            .links(this.links);

        simulation.alpha(0.3).restart();
    }
}
