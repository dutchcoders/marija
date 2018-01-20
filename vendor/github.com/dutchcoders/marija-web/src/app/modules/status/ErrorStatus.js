import React from 'react';

export default class ErrorStatus extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const { error } = this.props;

        if (!error) {
            return null;
        }

        return <div className="alert alert-danger">{ error } </div>;
    }
}
