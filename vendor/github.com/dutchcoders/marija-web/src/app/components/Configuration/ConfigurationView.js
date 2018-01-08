import React, { Component } from 'react';
import { connect } from 'react-redux';
import { find, sortBy, map, includes, slice } from 'lodash';

import { requestIndices } from '../../modules/indices/index';
import { fieldAdd, fieldDelete, dateFieldAdd, dateFieldDelete, normalizationAdd, normalizationDelete, indexAdd, indexDelete } from '../../modules/data/index';
import { serverAdd, serverRemove } from '../../modules/servers/index';
import { activateIndex, deActivateIndex } from '../../modules/indices/index';
import { batchActions } from '../../modules/batch/index';
import { getFields, clearAllFields, Field } from '../../modules/fields/index';
import { Icon } from '../index';

class ConfigurationView extends React.Component {
    constructor(props) {
        super(props);


        this.state = {
            normalization_error: '',
            currentFieldSearchValue: '',
            currentDateFieldSearchValue: ''
        };
    }

    handleAddField(path) {
        const { dispatch } = this.props;

        const icons = ["\u20ac", "\ue136", "\ue137", "\ue138", "\ue139", "\ue140", "\ue141", "\ue142", "\ue143"];

        const icon = icons[Math.floor((Math.random() * icons.length))];

        dispatch(fieldAdd({
            icon: icon,
            path: path
        }));
    }

    handleFieldSearchChange(event) {
        this.setState({currentFieldSearchValue: event.target.value});
    }

    handleDateFieldSearchChange(event) {
        this.setState({currentDateFieldSearchValue: event.target.value});
    }

    handleAddDateField(path) {
        const { date_field } = this.refs;
        const { dispatch } = this.props;

        dispatch(dateFieldAdd({
            path: path
        }));
    }

    handleAddNormalization(e) {
        e.preventDefault();

        const { regex, replaceWith  } = this.refs;
        const { dispatch } = this.props;

        if (regex.value === '') {
            return;
        }

        try {
            new RegExp(regex.value, "i");
        } catch (e) {
            this.setState({'normalization_error': e.message});
            return;
        }

        this.setState({'normalization_error': null});

        dispatch(normalizationAdd({
            regex: regex.value,
            replaceWith: replaceWith.value
        }));
    }

    handleAddIndex(e) {
        e.preventDefault();
        const { index } = this.refs;
        const { dispatch } = this.props;

        if (index.value === '') {
            return;
        }

        dispatch(batchActions(
            indexAdd(index.value),
            activateIndex(index.value)
        ));
    }

    handleAddServer(e) {
        e.preventDefault();

        const { server } = this.refs;
        const { dispatch } = this.props;

        if (server.value === '') {
            return;
        }

        dispatch(batchActions(
            serverAdd(server.value),
            requestIndices(server.value)
        ));
    }

    handleDeleteServer(server) {
        const { dispatch } = this.props;
        dispatch(serverRemove(server));
    }

    handleDeleteField(field) {
        const { dispatch } = this.props;
        dispatch(fieldDelete(field));
    }

    handleDeleteDateField(field) {
        const { dispatch } = this.props;
        dispatch(dateFieldDelete(field));
    }

    handleDeleteNormalization(normalization) {
        const { dispatch } = this.props;
        dispatch(normalizationDelete(normalization));
    }

    handleDeleteIndex(index) {
        const { dispatch } = this.props;
        dispatch(batchActions(
            indexDelete(index.id),
            deActivateIndex(index.id)
        ));
    }

    handleRequestIndices(server) {
        const { dispatch } = this.props;
        dispatch(requestIndices(server));
    }

    renderDateFields(fields) {
        const options = map(fields, (field) => {
            return (
                <li key={'date_field_' + field.path} value={ field.path }>
                    <i className="glyphicon">{ field.icon }</i>{ field.path }
                    <Icon onClick={() => this.handleDeleteDateField(field)} name="ion-ios-trash-outline"/>
                </li>
            );
        });

        let no_date_fields = null;

        if (fields.length == 0) {
            no_date_fields = <div className='text-warning'>No date fields configured.</div>;
        }

        return (
            <div>
                <ul>{ options }</ul>
                { no_date_fields }
                <form>
                    <div className="row">
                        <div className="col-xs-12">
                            <input className="form-control" value={this.state.currentDateFieldSearchValue}
                                   onChange={this.handleDateFieldSearchChange.bind(this)} type="text" ref="date_field"
                                   placeholder="Search date fields"/>
                        </div>
                    </div>
                </form>
            </div>
        );
    }


    renderFields(fields) {
        const style = {'marginRight': '20px'};

        const options = map(fields, (field) => {
            return (
                <li key={'field_' + field.path} value={ field.path }>
                    { field.path }
                    <i className="icon" style={ style }>{ field.icon }</i>
                    <Icon onClick={() => this.handleDeleteField(field)} name="ion-ios-trash-outline"/>
                </li>
            );
        });

        let no_fields = null;

        if (fields.length == 0) {
            no_fields = <div className='text-warning'>No fields configured.</div>;
        }

        return (
            <div>
                <ul>{ options }</ul>
                { no_fields }
                <form>
                    <div className="row">
                        <div className="col-xs-12">
                            <input className="form-control" value={this.state.currentFieldSearchValue}
                                   onChange={this.handleFieldSearchChange.bind(this)} type="text" ref="field"
                                   placeholder="Search fields"/>
                        </div>
                    </div>
                </form>
            </div>
        );
    }


    renderNormalizations(normalizations) {
        const { normalization_error } = this.state;

        const options = map(normalizations, (normalization) => {
            return (
                <li key={normalization.path} value={ normalization.path }>
                    <span>
                       Regex '<b>{normalization.regex}</b>' will be replaced with value '<b>{normalization.replaceWith}</b>'.
                    </span>
                    <Icon onClick={() => this.handleDeleteNormalization(normalization)} name="ion-ios-trash-outline"/>
                </li>
            );
        });

        let no_normalizations = null;

        if (normalizations.length == 0) {
            no_normalizations = <div className='text-warning'>No normalizations configured.</div>;
        }

        return (
            <div>
                <ul>{ options }</ul>
                { no_normalizations }
                <form onSubmit={this.handleAddNormalization.bind(this)}>
                    <div className="row">
                        <span className='text-danger'>{ normalization_error }</span>
                    </div>
                    <div className="row">
                        <div className="col-xs-10">
                            <input className="form-control" type="text" ref="regex" placeholder="regex"/>
                            <input className="form-control" type="text" ref="replaceWith" placeholder="replace value"/>
                        </div>
                        <div className="col-xs-1">
                            <Icon onClick={this.handleAddNormalization.bind(this)}
                                  name="ion-ios-add-circle-outline add"/>
                        </div>
                    </div>
                </form>
            </div>
        );
    }

    renderDatasources(datasources) {
        const { dispatch, activeIndices } = this.props;

        const options = map(sortBy(datasources, ["name"]), (datasource) => {
            const indexName = datasource.name;
            const active = find(activeIndices, (a) => a === datasource.id);

            return (
                <li key={ datasource.id } value={ indexName }>
                    <div className="index-name" title={indexName }>
                        { indexName }
                    </div>

                    { active ?
                        <Icon onClick={() => dispatch(deActivateIndex(datasource.id)) } name="ion-ios-eye"/>
                        :
                        <Icon onClick={() => dispatch(activateIndex(datasource.id)) } name="ion-ios-eye-off"/>
                    }
                </li>
            );
        });

        let no_datasources = null;
        if (datasources.length == 0) {
            no_datasources = <div className='text-warning'>No datasources configured.</div>;
        }

        return (
            <div>
                <ul>{options}</ul>
                { no_datasources }
            </div>
        );
    }

    render() {
        const { fields, date_fields, normalizations, datasources, availableFields, activeIndices, dispatch } = this.props;
        const { currentFieldSearchValue, currentDateFieldSearchValue } = this.state;


        return (
            <div>
                <div className="form-group">
                    <h2>Datasources</h2>
                    <p>Select and add the datasources to query.</p>
                    { this.renderDatasources(datasources) }
                </div>

                <div className="form-group">
                    <h2>
                        Fields
                        <Icon onClick={() => { dispatch(batchActions(clearAllFields(), getFields(activeIndices))); } } name="ion-ios-refresh" 
                            style={{float: "right", fontSize:"23px"}}/>
                    </h2>

                    <p>The fields are used as node id.</p>
                    { this.renderFields(fields) }

                    <ul>
                        {slice(availableFields.filter((item) => {
                            const inSearch = item.path.toLowerCase().indexOf(currentFieldSearchValue.toLowerCase()) === 0;
                            const inCurrentFields = fields.reduce((value, field) => {
                                if (value) {
                                    return true;
                                }
                                return field.path == item.path;
                            }, false);

                            return inSearch && !inCurrentFields;
                        }), 0, 10).map((item) => {
                            return (
                                <Field
                                    key={'available_fields_' + item.path}
                                    item={item} handler={() => this.handleAddField(item.path)}
                                    icon={'ion-ios-add-circle-outline'}/>
                            );
                        })}
                    </ul>
                </div>

                <div className="form-group">
                    <h2>
                        Date fields
                        <Icon onClick={() => { dispatch(batchActions(clearAllFields(), getFields(activeIndices))); } } name="ion-ios-refresh"
                              style={{float: "right", fontSize:"23px"}}/>
                    </h2>
                    <p>The date fields are being used for the histogram.</p>

                    { this.renderDateFields(date_fields) }

                    <ul>
                        {slice(availableFields.filter((item) => {
                            if (item.type !== 'date') {
                                return false;
                            }
                            

                            const inSearch = item.path.toLowerCase().indexOf(currentDateFieldSearchValue.toLowerCase()) === 0;
                            const inCurrentFields = fields.reduce((value, field) => {
                                if (value) {
                                    return true;
                                }

                                return field.path == item.path;
                            }, false);

                            return inSearch && !inCurrentFields;
                        }), 0, 10).map((item) => {
                            return (
                                <Field
                                    key={'available_date_fields_' + item.path}
                                    item={item} handler={() => this.handleAddDateField(item.path)}
                                    icon={'ion-ios-add-circle-outline'}/>
                            );
                        })}
                    </ul>
                </div>

                <div className="form-group">
                    <h2>Normalizations</h2>
                    <p>Normalizations are regular expressions being used to normalize the node identifiers and
                        fields.</p>
                    { this.renderNormalizations(normalizations) }
                </div>
            </div>
        );
    }
}


function select(state) {
    return {
        fields: state.entries.fields,
        availableFields: state.fields.availableFields,
        date_fields: state.entries.date_fields,
        normalizations: state.entries.normalizations,
        activeIndices: state.indices.activeIndices,
        datasources: state.entries.datasources
    };
}


export default connect(select)(ConfigurationView);
