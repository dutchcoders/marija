![](https://github.com/dutchcoders/marija-screenshots/blob/master/marija.png?raw=true)

# Marija [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/dutchcoders/marija?utm_source=badge&utm_medium=badge&utm_campaign=&utm_campaign=pr-badge&utm_content=badge) [![Go Report Card](https://goreportcard.com/badge/dutchcoders/marija)](https://goreportcard.com/report/dutchcoders/marija) [![Docker pulls](https://img.shields.io/docker/pulls/marija/marija.svg)](https://hub.docker.com/r/marija/marija/) [![Build Status](https://travis-ci.org/dutchcoders/marija.svg?branch=master)](https://travis-ci.org/dutchcoders/marija)

Marija is a data exploration and visualisation tool for (un)structured Elasticsearch data. Using Marija you'll be able to see relations between data of different datasources without any modifications to your data or index.

[Demo](https://demo.marija.io/)

# Screenshot

![](https://github.com/dutchcoders/marija-screenshots/blob/master/Screen%20Shot%202018-01-20%20at%2015.14.12.png?raw=true)

## Features

* work on multiple servers and indexes at the same time
* different fields can be used as node identifier
* identifiers can be normalized through normalization regular expressions
* each field will have its own icon
* query indexes natively to the datasource
* histogram view to identify nodes in time
* select and delete nodes
* select related nodes, deselect all but selected nodes
* zoom and move nodes
* navigate through selected data using the tableview
* multiple datasources supported like Elasticsearch, Splunk, Blockchain.info, Twitter
* submit nodes in realtime

## Install

### Using Docker

```
$ docker pull marija/marija
$ vim config-docker.toml # update elasticsearch configuration
$ docker run -d -p 8080:8080 -v $(pwd)/config-docker.toml:/config/config.toml --name marija marija/marija
```

### Installation from source

#### Install Golang

If you do not have a working Golang environment setup please follow [Golang Installation Guide](https://golang.org/doc/install).

#### Install Marija

Installation of Marija is easy.

```
$ go get github.com/dutchcoders/marija
$ marija
```

### Installation using Homebrew (macOS)

```
$ brew tap dutchcoders/homebrew-marija
$ brew install marija
```

## Configuration

### Elasticsearch
```
[datasource]

[datasource.elasticsearch]
type="elasticsearch"
url="http://127.0.0.1:9200/demo_index"
#username=
#password=

[[logging]]
output = "stdout"
level = "debug"
```


### Splunk

```
[datasource]

[datasource.splunk]
type="splunk"
url="https://localhost:8089"
username="admin"
password="admin"
```

## Contribute to Marija

Please follow Marija [Contributor's Guide](CONTRIBUTING.md)

## Community
Join the [marija-users](https://groups.google.com/forum/#!forum/marija-users) mailing list to discuss all things Marija.

## Copyright and license

Code and documentation copyright 2016-2018 Remco Verhoef [twitter](https://twitter.com/remco_verhoef).

Code released under [AGPL-3.0](LICENSE).

