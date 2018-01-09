const baseConfig = require('./webpack.config');
const path = require('path');

module.exports = Object.assign(baseConfig, {
    output: {
        path: [__dirname, 'dist'].join(path.sep),
        filename: 'bundle.js'
    },
});