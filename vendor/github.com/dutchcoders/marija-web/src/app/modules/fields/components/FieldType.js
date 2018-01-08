import React from 'react';

import { Icon } from '../../../components/index';

const fieldTypes = {
    string: 'ion-ios-text-outline',
    byte: 'ion-ios-information',
    short: 'ion-ios-text-outline',
    integer: 'ion-ios-information',
    long: 'ion-ios-text-outline',
    float: 'ion-ios-information',
    double: 'ion-ios-information',
    boolean: 'ion-ios-checkmark',
    date: 'ion-ios-calendar-outline',
    geo_point: 'ion-ios-navigate-outline',
    ip: 'ion-ios-locate-outline'
};

const unknownType = 'ion-ios-help';

/**
 * FieldType
 * @param props
 * @returns {XML}
 * @constructor
 */
export default function FieldType(props) {
    const icon = fieldTypes[props.field.type] || unknownType;
    return <Icon name={`${icon} type-indicator`} title={props.field.type} />
}
