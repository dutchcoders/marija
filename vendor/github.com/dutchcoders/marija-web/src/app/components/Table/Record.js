import React, { Component } from 'react';
import { uniq, find, map, mapValues, reduce } from 'lodash';
import { fieldLocator } from '../../helpers/index';
import { highlightNodes } from '../../modules/graph/index';
import { Icon } from '../../components/index';

export default class Record extends Component {
    constructor(props) {
        super(props);

        this.state = {
        };

        this.handleToggleExpand.bind(this);
    }

    handleToggleExpand(id) {
        const { toggleExpand } = this.props;
        toggleExpand(id);
    }

    render() {
        const { record, columns, node, searches } = this.props;
        const { expanded } = this.props;

        // make queries uniq
        const queries = (record.nodes || []).map((n) => {
            return (uniq(n.queries) || []).map((q) => {
                const search = find(searches, (s) => s.q == q);
                console.assert(search, "could not find query in search");

                if (find(search.items, (i) => i.id == record.id)) {
                    return (<Icon name='ion-ios-bulb' style={{ color: search.color }} alt={ search.q } />);
                }

                return (null);
            });
        });

        const renderedColumns = (columns || []).map((value) => {
            return (
                <td key={ 'column_' + record.id + value }>
                    <span className={'length-limiter'}>{ fieldLocator(record.fields, value) }</span>
                </td>
            );
        });

        return (
            <tr key={`record_${record.id}`} className={`columns ${expanded ? 'expanded' : 'closed'}`}>
                <td width="25" style={{'textAlign': 'center'}}>
                    <Icon onClick={() => this.handleToggleExpand(record.id) }
                          name={expanded ? 'ion-ios-remove' : 'ion-ios-add'}/>

                    { queries }
                </td>
                { renderedColumns}
            </tr>
        );
    }
}
