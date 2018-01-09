import React, { Component } from 'react';
import Tooltip from 'rc-tooltip';
import Loader from "../../../components/Misc/Loader";
import {Queries} from "../../../components/index";

export default class SearchBox extends Component {
    constructor(props) {
        super(props);

        this.handleSubmit = this.handleSubmit.bind(this);

        this.state = {
            query: '',
            selectValue: this.props.indexes[0]
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

    render() {
        const { itemsFetching, connected, enabled } = this.props;
        const { query } = this.state;

        let tooltipStyles = {};

        if (enabled) {
            tooltipStyles.visibility = 'hidden';
        }

        return (
            <nav id="searchbox" className="[ navbar ][ navbar-bootsnipp animate ] row" role="navigation">
                <div className="logo-container">
                    <img className={`logo ${connected ? 'connected' : 'not-connected'}`} src="/images/logo.png" title={connected ? "Marija is connected to the backendservice" : "No connection to Marija backend available" } />
                </div>
                <div className="queries-container">
                    <Queries />
                </div>
                <div className="search-container">
                    <div className="form-group">
                        <form onSubmit={this.handleSubmit.bind(this)}>
                            <Tooltip
                                overlay="Select at least one datasource and one field"
                                placement="bottomLeft"
                                overlayStyle={tooltipStyles}
                                arrowContent={<div className="rc-tooltip-arrow-inner" />}>

                                <input
                                    className="form-control"
                                    placeholder="Search something"
                                    value={ query }
                                    onChange={this.handleQueryChange.bind(this)}
                                    disabled={ !enabled }
                                />

                            </Tooltip>
                            <Loader show={itemsFetching} classes={['sk-search-box__loader']}/>
                        </form>
                    </div>
                </div>
            </nav>
        );
    }
}
