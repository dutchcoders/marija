import React, { Component } from 'react';
import { connect} from 'react-redux';
import Tooltip from 'rc-tooltip';
import SkyLight from 'react-skylight';
import SketchPicker from 'react-color';
import {Query} from "../../../components/index";
import {editSearch} from "../actions";
import {isEqual} from 'lodash';
import {headerHeightChange} from "../../../utils/actions";

class SearchBox extends Component {
    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            query: '',
            selectValue: this.props.indexes[0],
            editSearchValue: null
        };
    }

    handleSubmit(e) {
        e.preventDefault();

        const { query } = this.state;

        if (query === '') {
            return;
        }

        this.setState({query: ''});
        this.props.onSubmit(query, this.state.selectValue);
    }

    handleQueryChange(e) {
        this.setState({query: e.target.value});
    }

    handleEditSearch(search) {
        this.setState({editSearchValue: search});
        this.refs.editDialog.show();
    }

    handleChangeQueryColorComplete(color) {
        const { dispatch } = this.props;
        const search = Object.assign({}, this.state.editSearchValue);

        search.color = color.hex;

        dispatch(editSearch(search.q, {
            color: color.hex
        }));

        this.setState({editSearchValue: search});
    }

    componentDidUpdate(prevProps) {
        const { searches, dispatch } = this.props;

        if (!isEqual(searches, prevProps.searches)) {
            const height = this.refs.header.getBoundingClientRect().height;

            dispatch(headerHeightChange(height));
        }
    }

    render() {
        const { connected, enabled, searches } = this.props;
        const { query, editSearchValue } = this.state;

        let tooltipStyles = {};

        if (enabled) {
            tooltipStyles.visibility = 'hidden';
        }

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
            <nav id="searchbox" className="[ navbar ][ navbar-bootsnipp animate ] row" role="navigation" ref="header">
                <div className="logoContainer">
                    <img className={`logo ${connected ? 'connected' : 'not-connected'}`} src="/images/logo.png" title={connected ? "Marija is connected to the backendservice" : "No connection to Marija backend available" } />
                </div>
                <div className="queriesContainer">
                    {searches.map(search => <Query search={search} key={search.q} handleEdit={() => this.handleEditSearch(search)}/>)}

                    <form onSubmit={this.handleSubmit.bind(this)}>
                        <Tooltip
                            overlay="Select at least one datasource and one field"
                            placement="bottomLeft"
                            overlayStyle={tooltipStyles}
                            arrowContent={<div className="rc-tooltip-arrow-inner" />}>
                            <input
                                className="queryInput"
                                placeholder="Search"
                                value={ query }
                                onChange={this.handleQueryChange.bind(this)}
                                disabled={ !enabled }
                            />
                        </Tooltip>
                    </form>
                </div>
                <SkyLight dialogStyles={editQueryDialogStyles} hideOnOverlayClicked ref="editDialog" title="Update query" >
                    { edit_query }
                </SkyLight>
            </nav>
        );
    }
}

const select = (state, ownProps) => {
    return {
        ...ownProps,
        searches: state.entries.searches
    };
};

export default connect(select)(SearchBox);