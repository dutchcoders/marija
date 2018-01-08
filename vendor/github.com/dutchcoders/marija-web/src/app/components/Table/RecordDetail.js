import React, { Component } from 'react';
import { map, mapValues, reduce } from 'lodash';
import { fieldLocator } from '../../helpers/index';
import { highlightNodes } from '../../modules/graph/index';
import { Icon } from '../../components/index';
import DOMPurify from 'dompurify';

export default class Record extends Component {
    constructor(props) {
        super(props);

        this.state = {
            editNode: null,
            expanded: false
        };
    }

    handleTableAddColumn(field) {
        const { onTableAddColumn } = this.props;
        onTableAddColumn(field);
    }

    handleTableRemoveColumn(field) {
        const { onTableRemoveColumn } = this.props;
        onTableRemoveColumn(field);
    }

    handleMouseOver(id) {
        /*
        const { onMouseOver } = this.props;
        onMouseOver({nodes: [id]});
        */
    }

    extractAllFields(fields, keySeed = false) {
        return reduce(mapValues(fields, (value, key) => {
            const keyParts = [key];
            if (keySeed) {
                keyParts.unshift(keySeed);
            }

            const useKey = keyParts.join('.');
            if (value === null ) {
               return []; 
            } else if (typeof value.map == 'function') {
                return [useKey];
            } else if (typeof value == 'object') {
                return [].concat(this.extractAllFields(value, useKey));
            } else {
                return [useKey];
            }
        }), (result, value) => {
            return result.concat(value);
        });
    }

    renderDetails(columns) {
        const { record } = this.props;

        const allFields = this.extractAllFields(record.fields, false);

        const expandedFields = map(allFields, (value, key) => {
            const highlight =  (record.highlight || {});
            const field_value = highlight[value] || fieldLocator(record.fields, value) ;

            const clean = DOMPurify.sanitize(field_value, { ALLOWED_TAGS: ['p', 'br', 'em']});
            return (
                <tr key={ 'field_' + value }>
                    <td width="110">{value}
                        <Icon onClick={() => this.handleTableAddColumn(value)}
                            name="ion-ios-add-circle"
                            style={{marginLeft: '8px', lineHeight: '20px', fontSize: '12px'}}/>
                    </td>
                    <td colSpan="3" dangerouslySetInnerHTML={{ __html: clean }}></td>
                </tr>
            );
        });

        return ([
            <td>
            </td>,
            <td colSpan={columns.length ? columns.length : 1 }>
                <table className="details">
                    <tbody>{ expandedFields }</tbody>
                </table>
            </td>
        ]);
    }


    render() {
        const { record, columns, node, expanded } = this.props;
        if (!expanded) {
            return null;
        }

        return (
            <tr>
                { this.renderDetails(columns) }
            </tr>
        );
    }
}
