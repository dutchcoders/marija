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
        const { record, columns, node, searches, className } = this.props;
        const { expanded } = this.props;

        let queries = [];
        record.nodes.forEach(n => queries = queries.concat(n.queries));
        queries = uniq(queries);

        let usedSearches = [];
        queries.forEach(query => {
            const search = searches.find(s => s.q === query);

            if (search) {
                usedSearches.push(search);
            }
        });
        usedSearches = uniq(usedSearches);

        const queryElements = [];
        usedSearches.forEach(search => {
            queryElements.push(<Icon name='ion-ios-lightbulb' style={{ color: search.color }} alt={ search.q } key={search.q} />);
        });

        const renderedColumns = (columns || []).map((value) => {
            return (
                <td key={ 'column_' + record.id + value }>
                    <span className={'length-limiter'}>{ fieldLocator(record.fields, value) }</span>
                </td>
            );
        });

        return (
            <tr key={`record_${record.id}`} className={`columns record ${className} ${expanded ? 'expanded' : 'closed'}`}>
                <td width="25" style={{'textAlign': 'center'}}>
                    <Icon onClick={() => this.handleToggleExpand(record.id) }
                          name={expanded ? 'ion-ios-minus' : 'ion-ios-plus'}/>
                    { queryElements }
                </td>
                { renderedColumns}
            </tr>
        );
    }
}
