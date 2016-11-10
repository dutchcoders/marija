var webpack = require('webpack');
var path = require('path');

var buildPath = path.resolve(__dirname, 'dist');
var mainPath = path.resolve(__dirname, 'src', 'app', 'main.js');

var ExtractTextPlugin = require('extract-text-webpack-plugin');

module.exports = {
    entry: './src/app/main.js',
    target: 'web',
    devtool: 'source-map',
    plugins: [
        new ExtractTextPlugin('../dist/app.css'),
    ],
    resolve: {
        modulesDirectories: ['node_modules']
    },
    output: {
        path: [__dirname, 'dist'].join(path.sep),
        filename: 'bundle.js'
    },
    module: {
        loaders: [
            /*{
                test: /\.js$/,
                loaders: ["babel-loader", "eslint-loader"],
                exclude: /node_modules/
            },*/
            {
                test: /.js?$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
                query: {
                    presets: ['es2015', 'stage-1', 'react'],
                }
            },
            {
                test: /\.scss$/,
                loader: ExtractTextPlugin.extract('style-loader', 'css-loader?sourceMap!sass-loader?outputStyle=expanded&sourceMap=true&sourceMapContents=true')
            },
            {
                test: /\.(html)$/i,
                loader: "file-loader?name=/[name].[ext]"
            },
            {
                test: /\.(eot|svg|ttf|woff|woff2)([\?]?.*)$/,
                loader: 'file?name=fonts/[name].[ext]'
            },
            {
                test: /\.(png|jpg|jpeg)$/,
                loader: 'file?name=images/[name].[ext]'
            }
        ]
    }
};

