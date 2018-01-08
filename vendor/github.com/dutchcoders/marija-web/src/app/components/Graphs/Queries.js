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
import { deleteSearch } from '../../modules/search/index';


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
        dispatch(deleteSearch({search: query}));
    }

    handleChangeQueryColorComplete(color) {
        const search = this.state.editSearchValue;

        search.color = color.hex;

        this.setState({editSearchValue: search});
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
                    onChangeComplete={(color) => this.handleChangeSearchColorComplete(color) }
                />
                    <span className="colorBall" style={{backgroundColor: editSearchValue.color}}/>
                    { `${editSearchValue.q} (${editSearchValue.items.length})` }
            </form>;
        }

        return (
                <div className="queries">
                    <ul>
                    {map(queries, (query) => {
                        return (
                                <li key={query.q} style={{backgroundColor: query.color}}>
                                { `${query.q}` }&nbsp;<span className="count">{ `${query.items.length}`}<b>{`(${query.total})` }</b></span>
                                <Icon style={{'marginRight': '10px'}}
                                    onClick={(e) => this.handleEditQuery(query, e) }
                                    name="ion-ios-brush"/>
                                <Icon style={{'marginRight': '10px'}} onClick={(e) => this.handleDeleteQuery(query) }
                                    name="ion-ios-trash-outline"/>
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
    };
};

export default connect(select)(Dimensions()(Queries));
