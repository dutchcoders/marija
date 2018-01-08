import React, { Component } from 'react';
import SketchPicker from 'react-color';
import { connect } from 'react-redux';

import { map } from 'lodash';

import { deleteSearch } from '../index';
import { Icon } from '../../../components/index';

class Searches extends Component {

    constructor(props) {
        super(props);

        this.state = {
            editSearchValue: null
        };
    }

    handleEditSearch(search) {
        this.setState({editSearchValue: search});
    }

    handleCancelEditSearch() {
        this.setState({editSearchValue: null});
    }

    handleDeleteSearch(search) {
        const { dispatch } = this.props;
        dispatch(deleteSearch({search: search}));
    }

    handleChangeSearchColorComplete(color) {
        const search = this.state.editSearchValue;

        search.color = color.hex;

        this.setState({editSearchValue: search});
    }

    render() {
        const { searches } = this.props;
        const { editSearchValue } = this.state;

        return (
            <div className="form-group">
                <ul>
                    {map(searches, (search) => {
                        if (editSearchValue === search) {
                            return (
                                <li key={search.q}>
                                    <SketchPicker
                                        color={ search.color }
                                        onChangeComplete={(color) => this.handleChangeSearchColorComplete(color) }
                                    />
                                    <span className="colorBall" style={{backgroundColor: search.color}}/>

                                    { `${search.q} (${search.items.length})` }

                                    <Icon onClick={() => this.handleCancelEditSearch() }
                                          name="ion-ios-checkmark-circle"/>
                                </li>
                            );
                        } else {
                            return (
                                <li key={search.q}>
                                    <span className="colorBall" style={{backgroundColor: search.color}}/>

                                    { `${search.q} (${search.items.length})` }

                                    <Icon style={{'marginRight': '15px'}}
                                          onClick={(e) => this.handleEditSearch(search, e) }
                                          name="ion-ios-brush"/>
                                    <Icon onClick={(e) => this.handleDeleteSearch(search) }
                                          name="ion-ios-trash-outline"/>
                                </li>
                            );
                        }
                    })}
                </ul>
            </div>
        );

    }
}

function select(state) {
    return {
        searches: state.entries.searches
    };
}


export default connect(select)(Searches);
