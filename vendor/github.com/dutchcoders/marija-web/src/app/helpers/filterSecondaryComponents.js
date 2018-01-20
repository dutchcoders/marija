export default function filterSecondaryComponents(primaryQuery, components) {
    return components.filter(component => {
        const match = component.find(node => node.queries.indexOf(primaryQuery) !== -1);

        return typeof match !== 'undefined';
    });
}