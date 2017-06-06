<h1 align="center"><img align="middle" src="https://github.com/dutchcoders/marija-screenshots/blob/master/marija.png?raw=true">Marija</h1>

Marija is a data exploration and visualization tool for (un)structured data. Many [data sources](wiki/Datasources) will be possible in the future, but currently [Elasticsearch](https://www.elastic.co/), Twitter and Bitcoin are supported. With Marija you'll be able to see relations between data of different data sources without modifying your data or index. This works by creating a node for each datapoint and connecting related nodes to each other, creating a mesh structure. Each node can be clicked to view more detailed information and if spatial information is available, a histogram can be created for that node. Currently, Marija is being used to identify relations between spam campaigns, but it can be used for all kinds of different data sets. Read an article on using Marija on the Freedom Hosting hack data [here](https://hackernoon.com/analysing-freedom-hosting-ii-data-with-marija-fe64984a4e7f).

Disclaimer: Marija is still in an alpha stage, expect (many) bugs and changes. Please report bugs in the issue tracker.

<h3 align="center">Preview</h3>

![](https://github.com/dutchcoders/marija-screenshots/blob/master/Screen%20Shot%202016-11-17%20at%2009.46.31.png?raw=true)

## Demo
Try our demo with Elasticsearch, Twitter and Bitcoin datasource [here](http://demo.marija.io). Usage information on the demo application is also [provided](Demo). 

## Installation
[Instructions](Installation) are provided for installation with Docker, macOS (Homebrew) and from source. For an installation from source, [Golang](https://golang.org/) must be installed.

## Configuration
Configuration is easy, just modify the TOML file and change the data sources as you see fit.
View the [configuration](CONFIGURATION) and [data sources](Datasources) topics for the details.

## Features

* works on local and remote data sources simultaneously
* multiple fields can be used as a node identifier
* identifiers can be normalized through regular expressions
* each unique datasource field has its own icon
* indexes can be queried using regular Elasticsearch syntax
* a histogram view displays nodes on a time scale
* select and delete nodes
* select related nodes, deselect all but selected nodes
* zoom and move nodes
* navigate through selected data using the table view

## Workspace

Currently only one single workspace is supported. The workspace is being stored in the local storage of your browser. Next versions will support loading and saving multiple workspaces.

## Todo

* Optimize, optimize, optimize.

## Roadmap

We are working towards a stable version. These are some of the items on our feature wishlist. See the [issue list](https://github.com/Einzelganger/marija/issues) for features and bugs.

* analyze data in real-time
* create specialized tools based on Marija, for example a tool to graph packet traffic flows. 


## Contribute

Contributions are [welcome](Contribution_Guide).

## About

Marija is created by Dutchcoders. Find us [here](About).

