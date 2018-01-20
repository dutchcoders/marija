import React from 'react';

import { Icon } from '../../../components/index';
import { FieldType } from '../index';

export default function Field(props) {
    const { item, handler, icon } = props;

    return (
        <li className="field">
            <FieldType field={item} />

            {item.path}

            <Icon
                onClick={() => handler() }
                name={icon}/>
        </li>
    );
}
