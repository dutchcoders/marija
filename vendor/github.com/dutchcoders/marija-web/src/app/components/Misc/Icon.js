import React, { Component } from 'react';

export default class Icon extends Component {
    render() {
        const { name } = this.props;
        return (
            <i className={`ion ${name}`} {...this.props} />
        );
    }
}

Icon.propTypes = {
    name: React.PropTypes.string.isRequired
};
