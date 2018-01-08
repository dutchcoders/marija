import React, {Component} from 'react';
import { connect} from 'react-redux';
import Dimensions from 'react-dimensions';

import * as d3 from 'd3';
import { map, clone, groupBy, reduce, forEach, difference, find, uniq, remove, each, includes, assign, isEqual } from 'lodash';
import moment from 'moment';

import { nodesSelect, highlightNodes, nodeSelect } from '../../modules/graph/index';
import { normalize, fieldLocator } from '../../helpers/index';

var Worker = require("worker-loader!./Worker");

class Graph extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            nodes: [],
            links: [],
            highlight_nodes: [],
            edges: [],
            queries: [],
            clusters: {},
            start: new Date(),
            time: 0,
            n: {
            },
            ticks: 0
        };

        const { containerHeight, containerWidth } = props;

        this.network = {
            graph: {
                nodes: [],
                links: [],
                highlight_nodes: [],
                selection: null,
                queries: [],
                selectedNodes: [],
                tooltip: null,
                transform: d3.zoomIdentity
            },
            simulation: {},
            lines: {
                stroke: {
                    color: "#ccc",
                    thickness: 1
                }
            },
            nodes: {
                fill: {
                    color: "#333"
                },
                stroke: {
                    color: "#fff",
                    thickness: 3
                },
                sizeRange: [15, 30]
            },
            zoomed: function () {
                this.graph.transform = d3.event.transform;
            },
            setup: function (el) {
                let { clientHeight, clientWidth } = el;

                this.render = this.render.bind(this);
                this.drawNode = this.drawNode.bind(this);
                this.drawLink = this.drawLink.bind(this);

                this.canvas = el;
                this.context = this.canvas.getContext('2d');
                var canvas = d3.select(this.canvas);

                this.graph.transform.x = clientWidth / 2;
                this.graph.transform.y = clientHeight / 2;

                canvas.on("mousedown", this.mousedown.bind(this))
                    .on("mousemove", this.mousemove.bind(this))
                    .on("mouseup", this.mouseup.bind(this))
                    .call(d3.drag()
                          .filter(() => {
                              return d3.event.altKey;
                          })
                          .subject(this.dragsubject.bind(this))
                          .on("start", this.dragstarted.bind(this))
                          .on("drag", this.dragged.bind(this))
                          .on("end", this.dragended.bind(this))
                         )
                    .call(d3.zoom()
                          .filter(() => {
                              return d3.event.altKey;
                          })
                          .scaleExtent([1 / 2, 8])
                          .on("zoom", this.zoomed.bind(this))
                         );

                let { height, width } = this.canvas;
                
                this.worker = new Worker();
                this.worker.onmessage = function(event) {
                    switch (event.data.type) {
                    case "tick":
                        return this.ticked(event.data);
                    case "end":
                        return this.ended(event.data);
                    }
                    return false;
                }.bind(this);
                
                this.worker.postMessage({
                    clientWidth: width,
                    clientHeight: height,
                    type: "init"
                });

                requestAnimationFrame(this.render);
            },
            ticked: function(data) {
                this.graph.nodes = data.nodes;
                this.graph.links = data.links;
            },
            select: function (nodes) {
                this.graph.selectedNodes = nodes;
            },
            highlight: function (nodes) {
                this.graph.highlight_nodes = nodes;
            },
            updateNodes: function (graph) {
                this.graph.queries = graph.queries;

                this.worker.postMessage({
                    nodes: graph.nodes,
                    links: graph.links,
                    type: "update"
                });

            },
            render: function () {
                if (!this.graph) {
                    return;
                }

                // only when there is activity?
                requestAnimationFrame(this.render);

                // optimizations to do:
                // * group alike drawings, prevent switching of colors, fonts, etc
                // * stroke and fill all at once
                // * don't draw out of the viewing context
                //

                const { canvas, context, graph, lines } = this;

                // todo(nl5887): is this slow?
                const {width, height}  = canvas;

                context.save();
                context.clearRect(0, 0, width, height);

                context.translate((0) + graph.transform.x, (0) + graph.transform.y);

                context.scale(graph.transform.k, graph.transform.k);

                // draw links
                context.beginPath();

                // group all colors
                let color = lines.stroke.color;

                context.lineWidth = lines.stroke.thickness;

                for (const link of graph.links) {
                    this.drawLink(link);

                    if (link.color !== color) {
                        context.strokeStyle = color;
                        context.stroke();

                        color = link.color;
                    }
                };

                context.strokeStyle = color;
                context.stroke();

                // draw nodes
                for (let i = 0; i < graph.queries.length; i++) {
                    this.context.beginPath();

                    for (let j = 0; j < graph.nodes.length; j++) {
                        const d = graph.nodes[j];
                        this.drawNode(d, graph.queries[i]);
                    }

                    const color = graph.queries[i].color;

                    this.context.fillStyle = color;
                    this.context.fill();

                    this.context.strokeStyle = color;
                    this.context.stroke();
                }

                this.context.fillStyle = '#fff';

                for (let j = 0; j < graph.nodes.length; j++) {
                    const d = graph.nodes[j];
                    const fontHeight = 6 + Math.floor(0.8*d.r);
                    this.context.font=fontHeight + "px glyphicons halflings";
                    this.context.fillText(d.icon, d.x - ((fontHeight - 0.5) /2), d.y + (fontHeight + 0.5)/2);
                }

                // todo(nl5887): we're having graph and react nodes here, go fix.
                if (this.graph.selectedNodes) {
                    this.context.strokeStyle = '#993833';
                    this.context.lineWidth = this.nodes.stroke.thickness;

                    this.context.beginPath();
                    for (const selectedNode of this.graph.selectedNodes) {
                        const d = find(this.graph.nodes, (n) => {
                            return n.id == selectedNode.id;
                        });

                        if (!d) 
                            continue;

                        this.context.moveTo(d.x + d.r, d.y);
                        this.context.arc(d.x, d.y, d.r, 0, 2 * Math.PI);
                    }
                    this.context.stroke();
                }

                if (graph.selection) {
                    context.beginPath();
                    context.strokeStyle = '#c0c0c0';
                    context.fillStyle = "rgba(224, 224, 224, 0.6)";
                    //context.fillStyle = '#eee';
                    context.lineWidth = 1;
                    // context.setLineDash([6]);

                    context.rect(graph.selection.x1, graph.selection.y1, graph.selection.x2 - graph.selection.x1, graph.selection.y2 - graph.selection.y1);
                    context.fill();
                    context.stroke();
                }

                if (graph.tooltip) {
                    context.beginPath();
                    context.lineWidth="1";
                    context.strokeStyle="#cecece";
                    context.fillStyle="#fff";

                    const {width} = context.measureText(graph.tooltip.node.name);
                    context.rect(graph.tooltip.x,graph.tooltip.y-25,width,30);
                    context.stroke();
                    context.fill();

                    context.fillStyle = '#000'; //d.color[0];
                    context.font = "14px Arial";
                    context.fillText(graph.tooltip.node.name, graph.tooltip.x + 5, graph.tooltip.y - 5);
                }


                context.restore();
            },
            drawLink: function (d) {
                this.context.moveTo(d.source.x, d.source.y);
                this.context.lineTo(d.target.x, d.target.y);
            },
            drawNode: function (d, q) {
                // this.context.moveTo(d.x + d.r, d.y);
                // for each different query, show a part. This will show that the edge
                //  has been found in multiple queries.
                // this can be optimized by combining all same color fills apart
                for (var i = 0; i < d.queries.length; i++) {
                    var j = i;
                    for (; j < d.queries.length; j++) {
                        if (d.queries[i] !== d.queries[j]) {
                            break;
                        }
                    }

                    if (q.q === d.queries[i]) {
                        this.context.moveTo(d.x, d.y);
                        this.context.arc(d.x, d.y, d.r, 2 * Math.PI * (i / d.queries.length), 2 * Math.PI * ( (j + 1) / d.queries.length));
                        this.context.lineTo(d.x, d.y);
                    }

                    i = j;
                }
            },
            find(x, y) {
                var i = 0,
                    n = this.graph.nodes.length,
                    dx,
                    dy,
                    d2,
                    node,
                    closest;

                let radius = 20;
                if (radius == null) 
                    radius = Infinity;
                else 
                    radius *= radius;

                for (i = 0; i < n; ++i) {
                    node = this.graph.nodes[i];
                    dx = x - node.x;
                    dy = y - node.y;
                    d2 = dx * dx + dy * dy;
                    if (d2 < (node.r * node.r)) closest = node;
                }

                return closest;
            },
            mousedown: function () {
                const { graph } = this;

                if (d3.event.altKey) {
                    return;
                }

                var x = graph.transform.invertX(d3.event.layerX),
                    y = graph.transform.invertY(d3.event.layerY);

                var subject = this.find(x, y);
                if (!subject) {
                    graph.selection = {x1: x, y1: y, x2: x, y2: y};
                    this.onNodesSelect([]);
                    return;
                } else {
                    if (!includes(graph.selectedNodes, subject)) {
                        graph.selectedNodes.push(subject);
                    } else {
                        remove(graph.selectedNodes, subject);
                    }

                    this.onNodesSelect(graph.selectedNodes);
                }
            },
            mouseup: function () {
                const { graph } = this;

                if (d3.event.altKey) {
                    return;
                }

                var x = graph.transform.invertX(d3.event.layerX),
                    y = graph.transform.invertY(d3.event.layerY);

                if (graph.selection) {
                    graph.selection = assign(graph.selection, {x2: x, y2: y});

                    graph.nodes.forEach((d)=> {
                        if ((d.x > graph.selection.x1 && d.x < graph.selection.x2) &&
                            (d.y > graph.selection.y1 && d.y < graph.selection.y2)) {
                            if (!includes(graph.selectedNodes, d)) {
                                graph.selectedNodes.push(d);
                            }
                        }

                        if ((d.x > graph.selection.x2 && d.x < graph.selection.x1) &&
                            (d.y > graph.selection.y2 && d.y < graph.selection.y1)) {
                            if (!includes(graph.selectedNodes, d)) {
                                graph.selectedNodes.push(d);
                            }
                        }
                    });

                    this.onNodesSelect(graph.selectedNodes);

                    graph.selection = null;
                }
            },
            mousemove: function (n) {
                const { graph } = this;

                if (d3.event.altKey) {
                    return;
                }

                var x = graph.transform.invertX(d3.event.layerX),
                    y = graph.transform.invertY(d3.event.layerY);

                if (graph.selection) {
                    graph.selection = assign(graph.selection, {x2: x, y2: y});
                }

                var subject = this.find(x, y);
                if (subject === undefined) {
                    graph.tooltip = null;
                    this.onHighlightNode([]);
                } else if (!graph.tooltip || graph.tooltip.node !== subject) {
                    graph.tooltip = {node: subject, x: x, y: y};
                    this.onmousemove(subject);
                    this.onHighlightNode([subject]);
                }
            },
            dragstarted: function () {
                const { graph, simulation } = this;

                graph.selection = null;
                graph.tooltip = null;

                var x = graph.transform.invertX(d3.event.sourceEvent.layerX),
                    y = graph.transform.invertY(d3.event.sourceEvent.layerY);

                d3.event.subject.fx = (x);
                d3.event.subject.fy = (y);

                this.worker.postMessage({
                    nodes: [d3.event.subject],
                    type: 'restart'
                });
            },
            dragged: function () {
                const { graph, simulation } = this;

                var x = graph.transform.invertX(d3.event.sourceEvent.layerX),
                    y = graph.transform.invertY(d3.event.sourceEvent.layerY);

                d3.event.subject.fx = (x);
                d3.event.subject.fy = (y);

                this.worker.postMessage({
                    nodes: [d3.event.subject],
                    type: 'restart'
                });
            },
            dragended: function () {
                const { graph, simulation } = this;

                this.worker.postMessage({
                    nodes: [d3.event.subject],
                    type: 'stop'
                });
            },
            dragsubject: function () {
                const { graph, simulation } = this;

                const x = graph.transform.invertX(d3.event.x),
                      y = graph.transform.invertY(d3.event.y);

                return this.find(x, y);
            },
            mousemoved: function () {
            }
        };
    }

    componentDidMount() {
        const { network } = this;
        const { dispatch } = this.props;

        network.dispatch = dispatch;
        network.onmouseclick = this.onMouseClick.bind(this);
        network.onmousemove = this.onMouseMove.bind(this);
        network.onHighlightNode = this.onHighlightNode.bind(this);
        network.onNodesSelect = this.onNodesSelect.bind(this);

        network.setup(this.refs.canvas);
    }

    onMouseClick(node) {
        const { dispatch } = this.props;
        dispatch(nodeSelect({node: node}));
    }

    onHighlightNode(nodes) {
        // todo(nl5887): dispatch actual react (this.props.nodes, not graph nodes)
        if (isEqual(nodes, this.props.nodes)) {
            return;
        }

        const { dispatch } = this.props;
        dispatch(highlightNodes(nodes));
    }

    onNodesSelect(nodes) {
        // todo(nl5887): dispatch actual react (this.props.nodes, not graph nodes)
        const { dispatch } = this.props;
        dispatch(nodesSelect(nodes));
    }

    onMouseMove(node) {
        //const { dispatch } = this.props;
        //dispatch(nodeSelect({node: node}));
    }

    onMouseOver(node) {
        //const { dispatch } = this.props;
        //dispatch(nodeSelect({node: node}));
    }

    componentWillReceiveProps(nextProps) {
    }


    componentDidUpdate(prevProps, prevState) {
        const { network } = this;
        const { fields } = this.props;

        // todo(nl5887): only when adding or removing new nodes
        if (prevProps.nodes !== this.props.nodes) {
            network.updateNodes({
                nodes: this.props.nodes,
                links: this.props.links,
                queries: this.props.queries
            });
        }

        // todo(nl5887): we're having graph and react nodes here, go fix.
        network.select(this.props.node);

        network.highlight(this.props.highlight_nodes);
    }

    shouldComponentUpdate(){
        // todo(nl5887): not always update, only on changes
        return true;
    }

    render() {
        const { containerHeight, containerWidth } = this.props;

        return (
                <canvas
                    style={{fontFamily: 'glyphicons halflings'}}
                    width={containerWidth}
                    height={containerHeight}
                    ref="canvas">
                        histogram
                </canvas>
        );
    }
}

const select = (state, ownProps) => {
    return {
        ...ownProps,
        node: state.entries.node,
        nodes: state.entries.nodes,
        links: state.entries.links,
        queries: state.entries.searches,
        fields: state.entries.fields,
        items: state.entries.items,
        highlight_nodes: state.entries.highlight_nodes
    };
};

export default connect(select)(Dimensions()(Graph));
