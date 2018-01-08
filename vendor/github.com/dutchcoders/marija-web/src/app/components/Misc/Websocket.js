import React, { Component } from 'react';
import { connect } from 'react-redux';

import { Socket } from '../../utils/index'

class Websocket extends Component {

    componentWillMount() {
	const { dispatch } = this.props;
        Socket.startWS(dispatch);
    }

    render() {
        return null;
    }
}

function select(state, ownProps) {
    return {
        ...ownProps
    };
}

export default connect(select)(Websocket);
