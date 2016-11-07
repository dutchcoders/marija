# Marija [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/dutchcoders/marija?utm_source=badge&utm_medium=badge&utm_campaign=&utm_campaign=pr-badge&utm_content=badge) [![Go Report Card](https://goreportcard.com/badge/dutchcoders/marija)](https://goreportcard.com/report/dutchcoders/marija)

Marija is a graphing solution for (un)structured Elasticsearch data. Using Marija you'll be able to see connections 
between data of different indexes and datasources without any modifications to your data or index.

Currently Marija is being used to identify related spamruns, but can be used for all kind of different set, but can be used for all kind of different sets

# Screenshot

![](https://github.com/dutchcoders/marija-screenshots/blob/master/Screen%20Shot%202016-11-07%20at%2017.25.21.png?raw=true)

## Install

Installation of Marija is easy.

```
$ go get github.com/dutchcoders/marija
$ marija
```

## Usage

There are a few steps you need to take before you can start.

* add the elasticsearch server to the configuration, using the cloud icon you can provision the indexes
* enable the index(es) you want to search in
* add the correct fields, those fields are being used as unique identifier for each node 
* now you can enter your queries, and analyse the data

## Features

* normalization of identifiers
* support for multiple servers and indexes
* support for links on multiple fields
* add different icons to fields
* query the index using elasticsearch queries
* histogram
* identify datasources using the table pane

## Roadmap

We're working towards a first version. 

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

**Kevin Hoogerwerf**

## Copyright and license

Code and documentation copyright 2016 Remco Verhoef.
Code released under [the Apache license](LICENSE).

