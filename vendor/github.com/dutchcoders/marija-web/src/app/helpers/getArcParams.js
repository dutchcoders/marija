export default function getArcParams(x1, y1, x2, y2, bend) {
    // find mid point
    const averageX = (x1 + x2) / 2;
    const averageY = (y1 + y2) / 2;

    // get vector from start to end
    let deltaX = x2 - x1;
    let deltaY = y2 - y1;

    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);

    // normalise vector
    deltaX /= distance;
    deltaY /= distance;

    // This make the lines flatten at distance
    const flattenedBend = (bend * 300) / Math.pow(distance,1/16);

    // Arc amount bend more at distance
    const x3 = averageX + deltaY * flattenedBend;
    const y3 = averageY - deltaX * flattenedBend;

    // get the radius
    let radius = (0.5 * ((x1-x3) * (x1-x3) + (y1-y3) * (y1-y3)) / (flattenedBend));

    // use radius to get arc center
    const centerX = x3 - deltaY * radius;
    const centerY = y3 + deltaX * radius;

    // radius needs to be positive for the rest of the code
    radius = Math.abs(radius);

    // find angle from center to start and end
    let startAngle = Math.atan2(y1 - centerY, x1 - centerX);
    let endAngle = Math.atan2(y2 - centerY, x2 - centerX);

    // normalise angles
    startAngle = (startAngle + Math.PI * 2) % (Math.PI * 2);
    endAngle = (endAngle + Math.PI * 2) % (Math.PI * 2);

    // ensure angles are in correct directions
    if (bend < 0) {
        if (startAngle < endAngle) {
            startAngle += Math.PI * 2;
        }
    } else {
        if (endAngle < startAngle) {
            endAngle += Math.PI * 2;
        }
    }

    startAngle += 1 / radius * Math.sign(bend);
    endAngle -= 1 / radius * Math.sign(bend);

    return {
        centerX,
        centerY,
        radius,
        startAngle,
        endAngle
    };
}