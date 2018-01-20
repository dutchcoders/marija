import React, { Component } from 'react';
import {connect} from 'react-redux';
import {saveAs} from 'file-saver';
import { forEach, uniqWith, reduce, find, findIndex, pull, concat, map } from 'lodash';
import { Record, RecordDetail, Icon } from '../index';
import { tableColumnAdd, tableColumnRemove } from '../../modules/data/index';

class TableView extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            items: this.getSelectedItems(),
            expandedItems: [],
        };
    }

    toggleExpand(id) {
        if (findIndex(this.state.expandedItems, (o) => o === id) >= 0) {
            // remove 
            this.setState({expandedItems: pull(this.state.expandedItems, id)});
        } else {
            // add
            this.setState({expandedItems: concat(this.state.expandedItems, id)});
        }
    }

    handleTableAddColumn(field) {
        const { dispatch } = this.props;
        dispatch(tableColumnAdd(field));
    }

    handleTableRemoveColumn(dispatch, field) {
        dispatch(tableColumnRemove(field));
    }

    handleCancelEditNode(node) {
        const { dispatch } = this.props;
        this.setState({editNode: null});

    }

    getSelectedItems() {
        const { selectedNodes, items, fields, columns, dispatch, normalizations } = this.props;

        // todo(nl5887): this can be optimized
        let selectedItems = reduce(selectedNodes, (result, node) => {
            for (var itemid of node.items) {
                const i = find(items, (i) => itemid === i.id);
                if (!i) {
                    console.debug("could not find ${itemid} in items?", items);
                    continue;
                }

                // check if already exists
                if (find(result, (i) => itemid == i.id)) {
                    i.nodes.push(node);
                    continue;
                }

                i.nodes = [node];
                result.push(i);
            };

            return result;
        }, []);

        return uniqWith(selectedItems, (value, other) => {
            return (value.id == other.id);
        });
    }

    componentDidUpdate(prevProps) {
        if (prevProps.selectedNodes !== this.props.selectedNodes) {
            this.setState({items: this.getSelectedItems()});
        }
    }

    renderBody() {
        const { columns, searches, dispatch} = this.props;
        const { items } = this.state;

        return map(items, (record, i) => {
                const expanded = (findIndex(this.state.expandedItems, function(o) { return o == record.id; }) >= 0);
                const className = i % 2 === 0 ? 'odd' : 'even';

                return [
                    <Record
                        columns={ columns }
                        record={ record }
                        searches={ searches }
                        toggleExpand = { this.toggleExpand.bind(this) }
                        expanded = { expanded }
                        className={className}
                    />,
                    <RecordDetail
                        columns={ columns }
                        record={ record }
                        onTableAddColumn={(field) => this.handleTableAddColumn(field) }
                        onTableRemoveColumn={(field) => this.handleTableRemoveColumn(field) }
                        expanded = { expanded }
                        className={className}
                    />
                ];
        });
    }

    renderHeader() {
        const { columns, dispatch } = this.props;
        const { handleTableRemoveColumn } = this;

        return map(columns, function (value) {
            return (
                <th key={ 'header_' + value }>
                    { value }
                    <Icon onClick={(e) => handleTableRemoveColumn(dispatch, value)} name="ion-ios-trash-outline"/>
                </th>
            );
        });
    }

    exportCsv() {
        const { items } = this.state;
        const { columns } = this.props;
        const table = [];
        const delimiter = '|';

        table.push(columns.join(delimiter));

        items.forEach(item => {
            const row = [];
            columns.forEach(column => row.push(item.fields[column]));
            table.push(row.join(delimiter));
        });

        const csv = table.join("\n");

        const blob = new Blob(
            [csv],
            {type: "text/csv;charset=utf-8"}
        );

        const now = new Date();
        const dateString = now.getFullYear() + '-' + (now.getMonth() + 1) + '-' + now.getDate();
        const filename = 'marija-export-table-' + dateString + '.csv';

        saveAs(blob, filename);
    }

    render() {
        const { items } = this.state;

        if (!items.length) {
            return <p>Select some nodes first</p>;
        }

        return (
            <div className="form-group">
                <button className="btn btn-primary pull-right" onClick={this.exportCsv.bind(this)}>Export as CSV</button>

                <table>
                    <tbody>
                    <tr>
                        <th width="25">
                        </th>
                        { this.renderHeader() }
                    </tr>
                    {this.renderBody()}
                    </tbody>
                </table>
            </div>
        );
    }
}


function select(state) {
    return {
        selectedNodes: state.entries.node,
        normalizations: state.entries.normalizations,
        items: state.entries.items,
        searches: state.entries.searches,
        fields: state.entries.fields,
        columns: state.entries.columns
    };
}


export default connect(select)(TableView);
