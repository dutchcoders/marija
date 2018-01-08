export default function fieldLocator(document, field) {
    if (!field) {return false;}

    const field_levels = field.split('.');

    const value = field_levels.reduce((currentLevelInDocument, currentField) => {
        if (!currentLevelInDocument) {
            return null;
        }

        if (typeof currentLevelInDocument[currentField] !== 'undefined') {
            return currentLevelInDocument[currentField];
        }

        return null;
    }, document);

    return (value);
}
