import { each } from 'lodash';

const colors = [
    '#de79f2',
    '#917ef2',
    '#499df2',
    '#49d6f2',
    '#00ccaa',
    '#fac04b',
    '#bf8757',
    '#ff884d',
    '#ff7373',
    '#ff5252',
    '#6b8fb3'
];

export default function getQueryColor(searches) {
    const used = {};

    colors.forEach(color => used[color] = 0);
    searches.forEach(search => used[search.color] ++);

    let leastUsedColor;
    let leastUsedTimes = 10;

    // Find color that's used the least amount of times
    each(used, (times, color) => {
        if (times < leastUsedTimes) {
            leastUsedTimes = times;
            leastUsedColor = color;
        }
    });

    return leastUsedColor;
}