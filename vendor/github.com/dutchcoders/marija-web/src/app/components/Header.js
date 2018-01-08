import React, { Component } from 'react';
import { connect } from 'react-redux';

import { SearchBox, requestItems } from '../modules/search/index';
import { Icon } from '../components/index';

import { generateColour } from '../helpers/index';
import { openPane } from '../utils/index';

class Header extends Component {

    constructor(props) {
        super(props);

    }

    getColour(str) {
        return generateColour(str);
    }

    onSearchSubmit(q, index) {
        const { dispatch, activeIndices, queries } = this.props;

        const colors = ['#de79f2', '#917ef2', '#6d83f2', '#499df2', '#49d6f2', '#00ccaa', '#fac04b', '#bf8757', '#ff884d', '#ff7373', '#ff5252', '#6b8fb3'];

        dispatch(requestItems({
            query: q,
            from: 0, 
            size: 500,
            datasources: activeIndices,
            color: colors[queries.length % colors.length]
        }));
    }

    openPane() {
        const { dispatch } = this.props;
        dispatch(openPane('configuration'));
    }

    render() {
        const { connected, isFetching, total, indexes } = this.props;

        let errors = null;

        if (this.props.errors) {
            errors = <div className="alert alert-danger"><strong>Error executing query: </strong>{ this.props.errors }</div>;
        }

        return (
            <header className="row">
                <SearchBox
                    isFetching={isFetching}
                    total={total}
                    onSubmit={(q, index) => this.onSearchSubmit(q, index)}
                    connected={connected}
                    indexes={indexes}
                />
                { errors }
            </header>
        );
    }
}


function select(state) {
    return {
        isFetching: state.entries.isFetching,
        connected: state.entries.connected,
        errors: state.entries.errors,
        indexes: state.entries.indexes,
        activeIndices: state.indices.activeIndices,
        queries: state.entries.searches,
        total: state.entries.total
    };
}


export default connect(select)(Header);
