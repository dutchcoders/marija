export default function abbreviateNodeName(value, query, maxValueLength) {
    value += '';

    if (value.length <= maxValueLength) {
        return value;
    }

    const matchIndex = value.toLowerCase().indexOf(query.toLowerCase());

    if (matchIndex === -1) {
        return value.substring(0, maxValueLength - 1) + "...";
    }

    const match = value.substring(matchIndex, matchIndex + query.length);
    const contextLength = Math.round((maxValueLength - query.length) / 2);
    const contextLeft = value.substring(matchIndex - contextLength, matchIndex);
    const contextRight = value.substring(matchIndex + query.length, matchIndex + query.length + contextLength);

    let ret = '';

    if (contextLeft) {
        if (value.indexOf(contextLeft) > 0) {
            ret += '...';
        }

        ret += contextLeft;
    }

    ret += match;

    if (contextRight) {
        ret += contextRight;

        if (value.indexOf(contextRight) + contextRight.length < value.length) {
            ret += '...';
        }
    }

    return ret;
}