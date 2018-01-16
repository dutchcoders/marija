import React, {Component} from 'react';
import { connect} from 'react-redux';
import Dimensions from 'react-dimensions';
import SkyLight from 'react-skylight';
import SketchPicker from 'react-color';
import * as d3 from 'd3';
import { map, clone, groupBy, reduce, forEach, difference, find, uniq, remove, each, includes, assign, isEqual } from 'lodash';
import moment from 'moment';
import { nodesSelect, highlightNodes, nodeSelect } from '../../modules/graph/index';
import { normalize, fieldLocator } from '../../helpers/index';
import { Icon } from '../../components/index';
import { deleteSearch, setDisplayNodes, searchRequest } from '../../modules/search/index';
import Url from "../../domain/Url";
import Tooltip from 'rc-tooltip';
import {cancelRequest} from "../../utils/actions";

class Queries extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            editSearchValue: null
        };
    }

    componentDidMount() {
    }

    handleEditQuery(query) {
        this.setState({editSearchValue: query});

        this.refs.editDialog.show();
    }

    handleCancelEditQuery() {
        this.setState({editSearchValue: null});
    }

    handleDeleteQuery(query) {
        const { dispatch } = this.props;

        if (!query.completed) {
            dispatch(cancelRequest(query.requestId));
        }

        dispatch(deleteSearch({search: query}));
        Url.removeQueryParam('search', query.q);
    }

    handleDisplayMore(query) {
        const { dispatch } = this.props;
        const nodes = this.countNodesForQuery(query.q);
        const newNumber = query.displayNodes + 100;

        if (newNumber > Math.ceil(nodes / 100) * 100) {
            return;
        }

        dispatch(setDisplayNodes(query.q, newNumber));
    }

    handleDisplayLess(query) {
        const { dispatch } = this.props;
        const displayNodes = this.countDisplayNodesForQuery(query.q);
        let newNumber = query.displayNodes - 100;

        if (displayNodes < newNumber) {
            newNumber = Math.floor(displayNodes/ 100) * 100;
        }

        if (newNumber < 0) {
            return;
        }

        dispatch(setDisplayNodes(query.q, newNumber));
    }

    handleChangeQueryColorComplete(color) {
        const search = this.state.editSearchValue;

        search.color = color.hex;

        this.setState({editSearchValue: search});
    }

    countNodesForQuery(query) {
        const { nodes } = this.props;

        return nodes.filter(node => node.queries.indexOf(query) !== -1).length;
    }

    countDisplayNodesForQuery(query) {
        const { nodesForDisplay } = this.props;

        return nodesForDisplay.filter(node => node.queries.indexOf(query) !== -1).length;
    }

    render() {
        const { queries } = this.props;
        const { editSearchValue } = this.state;

        const editQueryDialogStyles = {
            backgroundColor: '#fff',
            color: '#000',
            width: '400px',
            height: '400px',
            marginTop: '-200px',
            marginLeft: '-200px',
        };

        let edit_query = null;
        if (editSearchValue) {
            edit_query = <form>
                <SketchPicker
                    color={ editSearchValue.color }
                    onChangeComplete={(color) => this.handleChangeQueryColorComplete(color) }
                />
                    <span className="colorBall" style={{backgroundColor: editSearchValue.color}}/>
                    { `${editSearchValue.q} (${editSearchValue.items.length})` }
            </form>;
        }

        return (
            <div className="queries">
                <ul>
                    {map(queries, (query) => {
                        const displayNodes = this.countDisplayNodesForQuery(query.q);
                        const nodes = this.countNodesForQuery(query.q);
                        const lessClass = 'ion ion-ios-minus ' + (displayNodes <= 0 ? 'disabled' : '');
                        const moreClass = 'ion ion-ios-plus ' + (displayNodes === nodes ? 'disabled' : '');
                        const itemClass = query.completed ? '' : 'loading';

                        return (
                                <li key={query.q} style={{backgroundColor: query.color}} className={itemClass}>
                                {query.q}&nbsp;
                                    <span className="count">
                                        {displayNodes}/
                                        {nodes}
                                    </span>

                                    <Tooltip
                                        overlay="Show less results"
                                        placement="bottom"
                                        arrowContent={<div className="rc-tooltip-arrow-inner" />}>
                                        <Icon
                                            onClick={() => this.handleDisplayLess(query) }
                                            name="ion-ios-minus"
                                            className={lessClass}
                                        />
                                    </Tooltip>

                                    <Tooltip
                                        overlay="Show more results"
                                        placement="bottom"
                                        arrowContent={<div className="rc-tooltip-arrow-inner" />}>
                                        <Icon
                                            onClick={() => this.handleDisplayMore(query) }
                                            name="ion-ios-plus"
                                            className={moreClass}
                                        />
                                    </Tooltip>

                                    <Tooltip
                                        overlay="Change color"
                                        placement="bottom"
                                        arrowContent={<div className="rc-tooltip-arrow-inner" />}>
                                        <Icon onClick={(e) => this.handleEditQuery(query, e) }
                                            name="ion-ios-gear"/>
                                    </Tooltip>

                                    <Tooltip
                                        overlay="Delete"
                                        placement="bottom"
                                        arrowContent={<div className="rc-tooltip-arrow-inner" />}>
                                        <Icon onClick={(e) => this.handleDeleteQuery(query) }
                                            name="ion-ios-close"/>
                                    </Tooltip>
                                </li>
                        );
                    })}
                </ul>
                <SkyLight dialogStyles={editQueryDialogStyles} hideOnOverlayClicked ref="editDialog" title="Update query" >
                    { edit_query }
                </SkyLight>
            </div>
        );
    }
}

const select = (state, ownProps) => {
    return {
        ...ownProps,
        queries: state.entries.searches,
        nodesForDisplay: state.entries.nodesForDisplay,
        nodes: state.entries.nodes,
        fields: state.entries.fields,
        activeIndices: state.indices.activeIndices
    };
};

export default connect(select)(Dimensions()(Queries));
