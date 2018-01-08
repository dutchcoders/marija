import React from 'react';

import { Icon } from '../index';
import { closePane, openPane } from '../../utils/index';

export function Pane(props) {
    const { handle, children, name, panes, icon, dispatch } = props;

    const isOpen = panes.reduce((value, item) => {
        if ((item.name == handle && item.state === true) || value === true) {
            return true;
        }
        return false;
    }, false);

    const open = () => {
        dispatch(openPane(handle));
    };

    const close = () => {
        dispatch(closePane(handle));
    };

    return (
        <div className={`pane ${handle} ${isOpen ? 'open' : 'closed'}`}>
            <div className="container-fluid">
                <div className="row pane-header">
                    <div className="col-md-12">
                        {name}
                        <Icon onClick={() => close()} name="ion-ios-remove shut"/>
                    </div>
                </div>
                <div className="row">
                    <div className="col-md-12 pane-holder">

                        {icon ?
                            <div onClick={() => open()} className="open-tag">
                                <Icon name={icon}/>
                            </div> :
                            null
                        }

                        <div className="col-md-12 pane-content">
                            {isOpen ? children : null}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
