import React from 'react';
import { Icon } from '../../components/index';

export default class ConnectionStatus extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        const {connected} = this.props;
        return (
            <div>
                {connected ?
                    <Icon
                        title="Marija backend server is functioning properly"
                        name="ion-ios-checkmark-circle status"
                        style={{'color': '#57c17b'}}/>
                    :
                    <Icon
                        title="Marija backend server is down"
                        name="ion-ios-close-circle status"
                        style={{'color': '#ea4552'}}/>
                }
            </div>
        );
    }
}
