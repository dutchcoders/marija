import { concat, merge } from 'lodash';
import { Workspaces, Migrations } from '../domain/index.js';
import Url from "../domain/Url";

export default function persistState() {
    return (next) => (reducer, initialState, enhancer) => {

        initialState = Workspaces.loadCurrentWorkspace(initialState)

        // parse all normalizations
        for (let i=0; i < (initialState.entries.normalizations || []).length; i++) {
            let normalization = initialState.entries.normalizations[i];
            normalization.re = new RegExp(normalization.regex, "i");
        }

        // Add datasources to url to keep local storage and url in sync
        initialState.indices.activeIndices.forEach(datasource => {
            Url.addQueryParam('datasources', datasource);
        });

        // Add datasources to url to keep local storage and url in sync
        initialState.entries.fields.forEach(field => {
            Url.addQueryParam('fields', field.path);
        });

        return next(reducer, initialState, enhancer);
    };
}
