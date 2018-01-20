import React, { Component } from 'react';
import { connect } from 'react-redux';

import { SearchBox, searchRequest } from '../modules/search/index';
import { Icon } from '../components/index';

import { generateColour } from '../helpers/index';
import { openPane } from '../utils/index';
import Url from "../domain/Url";

class Header extends Component {

    getColour(str) {
        return generateColour(str);
    }

    onSearchSubmit(q, index) {
        const { dispatch, activeIndices, fields } = this.props;

        Url.addQueryParam('search', q);

        dispatch(searchRequest({
            query: q,
            from: 0, 
            size: 500,
            datasources: activeIndices,
            fields: fields
        }));
    }

    openPane() {
        const { dispatch } = this.props;
        dispatch(openPane('configuration'));
    }

    render() {
        const { connected, total, indexes, fields, activeIndices } = this.props;

        let errors = null;

        if (this.props.errors) {
            errors = <div className="alert alert-danger"><strong>Error executing query: </strong>{ this.props.errors }</div>;
        }

        return (
            <header className="row">
                <SearchBox
                    total={total}
                    onSubmit={(q, index) => this.onSearchSubmit(q, index)}
                    connected={connected}
                    indexes={indexes}
                    enabled={fields.length > 0 && activeIndices.length > 0}
                />
                { errors }
            </header>
        );
    }
}


function select(state) {
    return {
        itemsFetching: state.entries.itemsFetching,
        connected: state.entries.connected,
        errors: state.entries.errors,
        indexes: state.entries.indexes,
        activeIndices: state.indices.activeIndices,
        queries: state.entries.searches,
        total: state.entries.total,
        fields: state.entries.fields
    };
}


export default connect(select)(Header);
