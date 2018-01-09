import React, { Component } from 'react';
import { connect } from 'react-redux';
import { find, sortBy, map, slice } from 'lodash';
import { requestIndices } from '../../modules/indices/index';
import { fieldAdd, fieldDelete, dateFieldAdd, dateFieldDelete, normalizationAdd, normalizationDelete, indexAdd, indexDelete } from '../../modules/data/index';
import { serverAdd, serverRemove } from '../../modules/servers/index';
import { activateIndex, deActivateIndex } from '../../modules/indices/index';
import { batchActions } from '../../modules/batch/index';
import { getFields, clearAllFields, Field } from '../../modules/fields/index';
import { Icon } from '../index';
import Url from "../../domain/Url";
import Loader from "../Misc/Loader";
import 'rc-tooltip/assets/bootstrap.css';
import {Workspaces} from "../../domain/index";
import {saveAs} from 'file-saver';
import {importData} from "../../modules/import/actions";

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

        Url.addQueryParam('fields', path);

        dispatch(fieldAdd(path));
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

        Url.addQueryParam('date-fields', path);

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

        Url.removeQueryParam('fields', field.path);

        dispatch(fieldDelete(field));
    }

    handleDeleteDateField(field) {
        const { dispatch } = this.props;

        Url.removeQueryParam('date-fields', field.path);

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

    handleDatasourceChange(event, id) {
        const { dispatch } = this.props;

        if (event.target.checked) {
            dispatch(activateIndex(id));
            Url.addQueryParam('datasources', id);
        } else {
            dispatch(deActivateIndex(id));
            Url.removeQueryParam('datasources', id);
        }
    }

    renderDateFields(fields, availableFields) {
        const { currentDateFieldSearchValue } = this.state;

        const options = map(fields, (field) => {
            return (
                <li key={'date_field_' + field.path} value={ field.path }>
                    <i className="glyphicon">{ field.icon }</i>{ field.path }
                    <Icon onClick={() => this.handleDeleteDateField(field)} name="ion-ios-trash-outline"/>
                </li>
            );
        });

        let no_date_fields = null;

        if (fields.length === 0) {
            no_date_fields = <div className='text-warning'>No date fields configured.</div>;
        }

        const availableDateFields = availableFields.filter(field => field.type === 'date');

        const search = (
            <form>
                <div className="row">
                    <div className="col-xs-12">
                        <input className="form-control" value={this.state.currentDateFieldSearchValue}
                               onChange={this.handleDateFieldSearchChange.bind(this)} type="text" ref="date_field"
                               placeholder={'Search ' + availableDateFields.length + ' date fields'} />
                    </div>
                </div>
            </form>
        );



        const available = (
            <ul>
                {slice(availableDateFields.filter((item) => {
                    const inSearch = item.path.toLowerCase().indexOf(currentDateFieldSearchValue.toLowerCase()) === 0;
                    const inCurrentFields = fields.reduce((value, field) => {
                        if (value) {
                            return true;
                        }

                        return field.path === item.path;
                    }, false);

                    return inSearch && !inCurrentFields;
                }), 0, 10).map((item) => {
                    return (
                        <Field
                            key={'available_date_fields_' + item.path}
                            item={item} handler={() => this.handleAddDateField(item.path)}
                            icon={'ion-ios-plus'}/>
                    );
                })}
            </ul>
        );

        let noAvailableDateFields = null;

        if (availableDateFields.length === 0) {
            noAvailableDateFields = <p>The selected datasources don't have any date fields.</p>;
        }

        return (
            <div>
                <ul>{ options }</ul>
                { noAvailableDateFields === null ? no_date_fields : null }
                { noAvailableDateFields }
                { availableDateFields.length > 0 ? search : null }
                { available }
            </div>
        );
    }


    renderFields(fields, availableFields) {
        const { currentFieldSearchValue } = this.state;

        const options = map(fields, (field) => {
            return (
                <li key={'field_' + field.path} value={ field.path }>
                    { field.path }
                    <i className="fieldIcon">{ field.icon }</i>
                    <Icon onClick={() => this.handleDeleteField(field)} name="ion-ios-trash-outline"/>
                </li>
            );
        });

        const search = (
            <form>
                <div className="row">
                    <div className="col-xs-12">
                        <input className="form-control" value={this.state.currentFieldSearchValue}
                               onChange={this.handleFieldSearchChange.bind(this)} type="text" ref="field"
                               placeholder={'Search ' + availableFields.length + ' fields'} />
                    </div>
                </div>
            </form>
        );

        const available = (
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
                            icon={'ion-ios-plus'}/>
                    );
                })}
            </ul>
        );

        let selectDatasourceMessage = null;

        if (availableFields.length === 0 && fields.length === 0) {
            selectDatasourceMessage = <p>First select a datasource.</p>;
        }

        return (
            <div>
                <ul>{ options }</ul>
                { availableFields.length > 0 ? search : null }
                { available }
                { selectDatasourceMessage }
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

                    <input type="checkbox" defaultChecked={active} onChange={(event) => this.handleDatasourceChange(event, datasource.id)} />

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

    getAtLeastOneAlert() {
        return (
            <span className="heading-alert">
                Select at least one
            </span>
        );
    }

    resetConfig() {
        Workspaces.deleteWorkspace();
        // Remove all data from url and refresh the page for simplicity
        // ideally this whole reset button is removed in the future
        // it was now mainly added to deal with erroneous data in either local
        // storage or the url
        window.location = '/';
    }

    exportJson() {
        const { entireState } = this.props;

        const blob = new Blob(
            [JSON.stringify(entireState)],
            {type: "text/json;charset=utf-8"}
        );

        const now = new Date();
        const dateString = now.getFullYear() + '-' + now.getMonth() + '-' + now.getDate();
        const filename = 'marija-export-' + dateString + '.json';

        saveAs(blob, filename);
    }

    chooseImportFile() {
        this.refs.importFile.click();
    }

    importJson(event) {
        const { dispatch } = this.props;

        const reader = new FileReader();

        reader.onload = function(){
            const store = JSON.parse(reader.result);

            dispatch(importData(store));
        };

        reader.readAsText(event.target.files[0]);
    }

    render() {
        const { fields, date_fields, normalizations, datasources, availableFields, activeIndices, dispatch, fieldsFetching } = this.props;

        return (
            <div>
                <div className="form-group">
                    <h2>
                        Datasources
                        {activeIndices.length === 0 ? this.getAtLeastOneAlert() : null}
                    </h2>
                    { this.renderDatasources(datasources) }
                </div>

                <div className="form-group">
                    <h2>
                        Fields
                        <Loader show={fieldsFetching} />
                        {activeIndices.length > 0 && fields.length === 0 ? this.getAtLeastOneAlert() : null}
                    </h2>

                    { this.renderFields(fields, availableFields) }


                </div>

                <div className="form-group">
                    <h2>
                        Date fields
                        <Loader show={fieldsFetching} />
                    </h2>
                    <p>The date fields are being used for the histogram.</p>

                    { this.renderDateFields(date_fields, availableFields) }
                </div>

                <div className="form-group">
                    <h2>Normalizations</h2>
                    <p>Normalizations are regular expressions being used to normalize the node identifiers and
                        fields.</p>
                    { this.renderNormalizations(normalizations) }
                </div>

                <div className="form-group">
                    <button className="btn btn-primary" onClick={this.exportJson.bind(this)}>Export</button>
                    <input type="file" ref="importFile" className="importFile" onChange={this.importJson.bind(this)} />
                    <button className="btn btn-primary" onClick={this.chooseImportFile.bind(this)}>Import</button>
                </div>

                <div className="form-group">
                    <button className="btn btn-primary" onClick={this.resetConfig.bind(this)}>Reset config</button>
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
        datasources: state.entries.datasources,
        fieldsFetching: state.fields.fieldsFetching,

        // For the export to JSON feature
        entireState: state
    };
}


export default connect(select)(ConfigurationView);
