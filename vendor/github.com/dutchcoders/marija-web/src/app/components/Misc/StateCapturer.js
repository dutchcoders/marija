import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Workspaces, Migrations } from '../../domain/index'

class StateCapturer extends Component {

    componentWillReceiveProps(nextProps){
        const { store } = this.props;

        const persistedState = Object.assign({}, store.getState());
        Workspaces.persistCurrentWorkspace(persistedState)
    }

    render() {
        return null;
    }
}

function select(state, ownProps) {
    return {
        ...ownProps,

        columns: state.entries.columns,
        searches: state.entries.searches,
        entries_fields: state.entries.fields,
        date_fields: state.entries.date_fields,
        indexes: state.entries.indexes,
        datasources: state.entries.datasources,
        normalizations: state.entries.normalizations,

        indices: state.indices,
        utils: state.utils,
        servers: state.servers,
        fields: state.fields,
    };
}

export default connect(select)(StateCapturer);
