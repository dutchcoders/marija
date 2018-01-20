import { INDEX_ACTIVATED, INDEX_DEACTIVATED } from '../modules/indices/index';

import { union, without } from 'lodash';

const defaultState = {
    activeIndices: []
};

export default function indices(state = defaultState, action) {
    switch (action.type) {
        case INDEX_ACTIVATED:
            const newActiveIndices = union(state.activeIndices, [action.payload.index]);
            return Object.assign({}, state, {activeIndices: newActiveIndices});

        case INDEX_DEACTIVATED:
            const filteredIndices = without(state.activeIndices, action.payload.index);
            return Object.assign({}, state, {activeIndices: filteredIndices});

        default:
            return state;
    }
}
