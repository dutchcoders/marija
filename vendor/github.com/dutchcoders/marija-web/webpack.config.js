const webpack = require('webpack');
const path = require('path');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const dotenv = require('dotenv');
const { gitDescribeSync } = require('git-describe');

dotenv.config();
const gitInfo = gitDescribeSync();

module.exports = {
    entry: './src/app/main.js',
    target: 'web',
    devtool: 'source-map',
    plugins: [
        new ExtractTextPlugin('../dist/app.css'),
        new webpack.DefinePlugin({
            "process.env": { 
                NODE_ENV: JSON.stringify(process.env.NODE_ENV || "development"),
                WEBSOCKET_URI: process.env.WEBSOCKET_URI ? JSON.stringify(process.env.WEBSOCKET_URI) : null,
                CLIENT_VERSION: JSON.stringify(gitInfo.raw)
            }
        })
    ],
    resolve: {
        modules: ['node_modules', 'src'],
        extensions: ['.js', '.scss']
    },
    module: {
        loaders: [
            {
                test: /.js?$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
                query: {
                    presets: ['es2015', 'stage-1', 'react'],
                }
            },
            {
                test: /\.(scss|css)$/,
                loaders: ['style-loader','css-loader','sass-loader']
            },
            {
                test: /\.(html)$/i,
                loader: 'file-loader?name=/[name].[ext]'
            },
            {
                test: /\.(eot|svg|ttf|woff|woff2)([\?]?.*)$/,
                loader: 'file-loader?name=fonts/[name].[ext]'
            },
            {
                test: /\.(png|jpg|jpeg)$/,
                loader: 'file-loader?name=images/[name].[ext]'
            }
        ]
    },
    node: {
        fs: 'empty',
        child_process: 'empty'
    },
};

