import createBrowserHistory from 'history/createBrowserHistory';
import queryString from 'query-string';

const history = createBrowserHistory();
let currentLocation = history.location;

// Update our current location whenever the url changes
history.listen(location => {
    currentLocation = location;
});

export default class Url {
    /**
     * Add a new query param to the url
     *
     * @param name
     * @param value
     */
    static addQueryParam(name, value) {
        const newQueryString = this.addToQueryString(currentLocation.search, name, value);

        this.updateQueryString(newQueryString);
    }

    /**
     * Remove a query param from the url
     *
     * @param name
     * @param value
     */
    static removeQueryParam(name, value) {
        const newQueryString = this.removeFromQueryString(currentLocation.search, name, value);

        this.updateQueryString(newQueryString);
    }

    static removeAllQueryParams(name) {
        const newQueryString = this.removeAllFromQueryString(currentLocation.search, name);

        this.updateQueryString(newQueryString);
    }

    /**
     * Pushes the new query string on to the history object
     *
     * @param newQueryString
     */
    static updateQueryString(newQueryString) {
        history.push(currentLocation.pathname + newQueryString);
    }

    /**
     * Pure function that adds a {value} to {name} in a query string. Returns a new query string.
     *
     * @param currentQueryString
     * @param name
     * @param value
     * @returns {string}
     */
    static addToQueryString(currentQueryString, name, value) {
        const queryParams = queryString.parse(currentQueryString);

        // Check if this name was already in the url
        if (queryParams[name]) {
            const items = queryParams[name].split(',');

            // Only add if it doesnt exist, don't add duplicates
            if (items.indexOf(value) === -1) {
                items.push(value);
            }

            queryParams[name] = items.join(',');
        } else {
            queryParams[name] = value;
        }

        return '?' + queryString.stringify(queryParams);
    }

    /**
     * Pure function that removes a {value} from {name} in a query string. Returns a new query string.
     *
     * @param currentQueryString
     * @param name
     * @param value
     * @returns {string}
     */
    static removeFromQueryString(currentQueryString, name, value) {
        const queryParams = queryString.parse(currentQueryString);

        // Check if this name was already in the url
        if (queryParams[name]) {
            let items = queryParams[name].split(',');

            // Remove from the array
            items = items.filter(item => item !== value);

            if (items.length > 0) {
                queryParams[name] = items.join(',');
            } else {
                delete queryParams[name];
            }
        }

        return '?' + queryString.stringify(queryParams);
    }

    static removeAllFromQueryString(currentQueryString, name) {
        const queryParams = queryString.parse(currentQueryString);

        delete queryParams[name];

        return '?' + queryString.stringify(queryParams);
    }
}