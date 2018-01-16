import React, { Component } from 'react';
import { connect } from 'react-redux';
import { find, map, groupBy, reduce, forEach, filter, concat } from 'lodash';
import Dimensions from 'react-dimensions';

import * as d3 from 'd3';
import moment from 'moment';

import { normalize, fieldLocator } from '../../helpers/index';

class Histogram extends React.Component {
    constructor(props) {
        super(props);

        this.draw = this.draw.bind(this);
        this.state = {};
    }

    componentDidMount() {
        this.canvas = this.refs.canvas;
        this.context = this.canvas.getContext('2d');
        this.draw();
    }

    componentDidUpdate(prevProps, prevState) {
        console.debug("componentDidUpdate (histogram");
        // group items to periods using lodash? complete set
        // have selection filter and drag timeline to select nodes
        //

        this.draw();
    }

    draw() {
        const { node, normalizations, fields, date_fields, items } = this.props;
        const { canvas } = this;

        if (!items.length || !date_fields.length) {
            return;
        }

        let itemsCopy = concat([], items);

        const context = this.context;
        context.save();

        var margin = {top: 20, right: 20, bottom: 30, left: 40},
            width = canvas.width - margin.left - margin.right,
            height = canvas.height - margin.top - margin.bottom;

        this.context.clearRect(0, 0, canvas.width, canvas.height);

        const x = d3.scaleBand()
            .rangeRound([0, width])
            .padding(0.1);

        const y = d3.scaleLinear()
            .rangeRound([height, 0]);

        context.translate(margin.left, margin.top);

        const groupedResults = groupBy(itemsCopy, (result) => {
            for (var date_field of date_fields) {
                let date = fieldLocator(result.fields, date_field.path);
                if (!date) {
                    continue;
                }

                return moment(date).year() + '-' + moment(date).month();
            }
        });

        const minX = reduce(itemsCopy, (min, result) => {
            for (var date_field of date_fields) {
                let date = fieldLocator(result.fields, date_field.path);
                if (!date) {
                    continue;
                }

                return (moment(date) < min ? moment(date) : min);
            }
        }, moment());

        const maxX = reduce(itemsCopy, (max, result) => {
            for (var date_field of date_fields) {
                let date = fieldLocator(result.fields, date_field.path);
                if (!date) {
                    continue;
                }

                return (moment(date) > max ? moment(date) : max);
            }
        }, 0);

        const periods = [];

        let year = minX.year();
        let month = minX.month();
        for (; year < maxX.year() || (year == maxX.year() && month < maxX.month());) {
            month++;
            if (month > 12) {
                year++;
                month = 1;
            }

            periods.push(year + "-" + month);
        }

        x.domain(periods);

        const maxValue = reduce(groupedResults, (max, n, m) => (n.length > max ? n.length : max), 0);
        y.domain([0, maxValue]);

        const yTickCount = 10,
            yTicks = y.ticks(yTickCount),
            yTickFormat = y.tickFormat(yTickCount);

        context.beginPath();
        x.domain().forEach((d) => {
            context.moveTo(x(d) + x.bandwidth() / 2, height);
            context.lineTo(x(d) + x.bandwidth() / 2, height + 6);
        });

        // todo(nl5887): add selection
        // todo(nl5887): show difference between queries

        context.strokeStyle = "#b5b5b5";
        context.stroke();

        context.textAlign = "center";
        context.fillStyle = "#b5b5b5";
        context.textBaseline = "top";
        x.domain().forEach((d) => {
            context.fillText(d, x(d) + x.bandwidth() / 2, height + 6);
        });

        context.beginPath();
        yTicks.forEach((d) => {
            context.moveTo(0, y(d) + 0.5);
            context.lineTo(-6, y(d) + 0.5);
        });
        context.strokeStyle = "#b5b5b5";
        context.stroke();

        context.textAlign = "right";
        context.textBaseline = "middle";
        context.textColor = "#b5b5b5";
        yTicks.forEach((d) => {
            context.fillText(yTickFormat(d), -9, y(d));
        });

        context.beginPath();
        context.moveTo(-6.5, 0 + 0.5);
        context.lineTo(0.5, 0 + 0.5);
        context.lineTo(0.5, height + 0.5);
        context.lineTo(-6.5, height + 0.5);
        context.strokeStyle = "#b5b5b5";
        context.stroke();

        context.fillStyle = "#57c17b";

        forEach(groupedResults, (d, v) => {
            context.fillRect(x(v), y(d.length), x.bandwidth(), height - y(d.length));
        });

        if (node.length > 0) {
            itemsCopy = filter(itemsCopy, (item) => {
                // check if node contains item
                return reduce(fields, (found, field) => {
                    const val = fieldLocator(item.fields, field.path);
                    if (!val) {return found;}

                    return found || find(node, (o) => {
                        return (o.id === normalize(normalizations, val));
                    });
                }, false);
            });

            const groupedResultsSelection = groupBy(itemsCopy, (result) => {
                for (var date_field of date_fields) {
                    let date = fieldLocator(result.fields, date_field.path);
                    if (!date) {
                        continue;
                    }

                    return moment(date).year() + '-' + moment(date).month();
                }
            });

            context.fillStyle = "#b5b5b5";
            forEach(groupedResultsSelection, (d, v) => {
                context.fillRect(x(v), y(d.length), x.bandwidth(), height - y(d.length));
            });
        }


        context.restore();
    }

    render() {
        const { items, containerHeight, containerWidth } = this.props;

        let noitems = null; 
        if (items.length == 0) {
            noitems = <div>No items</div>;
        }

        return (
            <div>
                { noitems }
                <canvas
                    width={ containerWidth }
                    height={ containerHeight }
                    ref="canvas"
                />
            </div>
        );

    }
}

const select = (state, ownProps) => {
    return {
        ...ownProps,
        node: state.entries.node,
        queries: state.entries.queries,
        fields: state.entries.fields,
        normalizations: state.entries.normalizations,
        date_fields: state.entries.date_fields,
        items: state.entries.items,
        highlight_nodes: state.entries.highlight_nodes
    };
};

export default connect(select)(Dimensions()(Histogram));
