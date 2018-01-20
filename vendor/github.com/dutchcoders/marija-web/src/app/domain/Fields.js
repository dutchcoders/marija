import { reduce, merge, filter } from 'lodash';

export default class Fields {

    /**
     * getTypes
     * @param mappings
     * @returns {Array|*}
     */
    static getTypes(mappings) {
        const defaultType = '_default_';

        return filter(Object.keys(mappings), (type) => {
            return type !== defaultType;
        });
    }


    /**
     * getFieldsFromResult
     * @param fields
     * @returns {*}
     */
    static getFieldsFromResult(fields) {
        return fields;
    }


    /**
     * recurseFields
     * @param type
     * @param fieldsContainer
     * @param base
     * @param fields
     * @returns {*|Object}
     */
    static recurseFields(type, fieldsContainer, base, fields = []) {

        let foundFields = [];

        const shouldExtract = ['properties', 'fields'].reduce((check, field) => {
            if (typeof fieldsContainer[field] === 'object') {
                return fieldsContainer[field];
            }
            return check;
        }, false);


        if (shouldExtract) {
            const field_keys = Object.keys(shouldExtract);

            foundFields = reduce(field_keys, (results, field) => {
                const innerFields = Fields.recurseFields(
                    type,
                    shouldExtract[field],
                    field,
                    []
                );

                const combinedFields = fields.concat(innerFields, results);
                if (typeof shouldExtract[field].type != 'undefined') {
                    combinedFields.push({
                        path: base ? [base, field].join('.') : field,
                        document_type: type,
                        type: shouldExtract[field].type || "nested",
                        format: shouldExtract[field].format || null
                    });
                }

                return combinedFields;
            }, []);
        }

        return merge(fields, foundFields);
    }


    /**
     * extractNewFields
     * @param fields
     * @param currentFields
     * @returns {*}
     */
    static extractNewFields(fields, currentFields) {
        return reduce(fields, (allFields, newItem) => {
            const notAllreadySavedToState = typeof currentFields.find((item) => item.path === newItem.path) == 'undefined';
            const notCurrentlySavingToState = typeof allFields.find((item) => item.path === newItem.path) == 'undefined';
            
            if (notAllreadySavedToState && notCurrentlySavingToState) {
                allFields.push(newItem);
            }

            return allFields;
        }, []);
    }

}
