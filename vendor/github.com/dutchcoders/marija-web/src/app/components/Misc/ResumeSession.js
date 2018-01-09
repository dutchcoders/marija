import React, { Component } from 'react';
import { connect } from 'react-redux';
import queryString from 'query-string';
import {dateFieldAdd, fieldAdd} from "../../modules/data/actions";
import {activateIndex} from "../../modules/indices/actions";
import {requestItems} from "../../modules/search/actions";
import Url from "../../domain/Url";

class ResumeSession extends Component {
    componentDidMount() {
        const { history, location, dispatch } = this.props;

        const parsed = queryString.parse(location.search);

        if (parsed.fields) {
            this.addFields(parsed.fields.split(','));
        }

        if (parsed['date-fields']) {
            this.addDateFields(parsed['date-fields'].split(','));
        }

        if (parsed.datasources) {
            this.addDatasources(parsed.datasources.split(','));
        }

        if (parsed.search) {
            this.search(parsed.search.split(','));
        }
    }

    addFields(fields) {
        const { dispatch } = this.props;

        fields.forEach(field => {
            dispatch(fieldAdd(field));
        });
    }

    addDateFields(fields) {
        const { dispatch } = this.props;

        fields.forEach(field => {
            dispatch(dateFieldAdd(field));
        });
    }

    addDatasources(ids) {
        const { dispatch } = this.props;

        ids.forEach(id => {
            dispatch(activateIndex(id));
        });
    }

    search(terms) {
        const { dispatch } = this.props;

        dispatch((dispatch, getState) => {
            const state = getState();

            // Get the updated state, because it might have changed by
            // adding datasources from the url
            const datasources = state.indices.activeIndices;

            if (datasources.length === 0) {
                this.cancelSearch();
                return;
            }

            const fields = state.entries.fields;

            if (fields.length === 0) {
                this.cancelSearch();
                return;
            }

            terms.forEach(term => {
                dispatch(requestItems({
                    query: term,
                    datasources: datasources,
                    from: 0,
                    size: 500
                }));
            });
        });
    }

    /**
     * Tried to perform search, but couldnt continue. Remove the search queries
     * from the url
     */
    cancelSearch() {
        Url.removeAllQueryParams('search');
    }

    render() {
        return null;
    }
}

// Empty function, needed for mapping dispatch method to the props
const select = () => {
    return {};
};

export default connect(select)(ResumeSession);