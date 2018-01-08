import { concat, merge } from 'lodash';
import { Workspaces, Migrations } from '../domain/index.js';

export default function persistState() {
    return (next) => (reducer, initialState, enhancer) => {

        initialState = Workspaces.loadCurrentWorkspace(initialState)

        // parse all normalizations
        for (let i=0; i < (initialState.entries.normalizations || []).length; i++) {
            let normalization = initialState.entries.normalizations[i];
            normalization.re = new RegExp(normalization.regex, "i");
        }

        return next(reducer, initialState, enhancer);
    };
}
