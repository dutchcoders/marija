![](https://github.com/dutchcoders/marija-screenshots/blob/master/marija.png?raw=true)

# Marija [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/dutchcoders/marija?utm_source=badge&utm_medium=badge&utm_campaign=&utm_campaign=pr-badge&utm_content=badge) [![Go Report Card](https://goreportcard.com/badge/dutchcoders/marija)](https://goreportcard.com/report/dutchcoders/marija) [![Docker pulls](https://img.shields.io/docker/pulls/marija/marija.svg)](https://hub.docker.com/r/marija/marija/) [![Build Status](https://travis-ci.org/dutchcoders/marija.svg?branch=master)](https://travis-ci.org/dutchcoders/marija)

Marija is a data exploration and visualisation tool for (un)structured data. There are several datasources possible, but currently [Elasticsearch](https://www.elastic.co/), Twitter and Bitcoin are supported. With Marija you'll be able to see relations between data of different datasources without modifying your data or index.

Currently, Marija is being used to identify related spamruns, but it can be used for all kinds of different data sets.

Disclaimer: Marija is still in an alpha stage, expect (many) bugs and changes. Please report bugs in the issue tracker.

# Screenshot

![](https://github.com/dutchcoders/marija-screenshots/blob/master/Screen%20Shot%202016-11-17%20at%2009.46.31.png?raw=true)

## Install

### Using Docker

```
$ docker pull marija/marija
$ docker run -d -p 8080:8080 --name marija marija/marija
```

### Installation from source
Marija depends on Golang version 1.7 or higher. If you do not have a working Golang environment setup or are using an older version, please follow [Golang Installation Guide](https://golang.org/doc/install).
Elasticsearch is also 

Installing Marija is easy.

```
$ go get github.com/dutchcoders/marija
```

### Installation using Homebrew (macOS)

```
$ brew tap dutchcoders/homebrew-marija
$ brew install marija
```

### Running Marija
First, add your datasources to config.toml. To be able to use Elasticsearch as a datasource it should be installed, but Bitcoin and Twitter searches work out of the box. 
When finished, run the application by simply invoking ```marija```.
Optionally specify the configuration file with the parameter ```--config config.toml```. Open http://127.0.0.1:8080 in your browser to view the application.
		
## Usage
 Take the following steps in the configuration window in the application:

* enable the datasources by clicking on the eye icon next to the source.
* click the refresh icon of the FIELDS section to display the list of available fields for that datasource
* add the fields you want to use as nodes in the overview
* additionally you can add the date field you want to use for the histogram
* and add some normalizations (eg removing part of the identifier) using regular expressions

You're all set up now, just type your queries and start exploring your data in the mesh. 

For detailed information, the nodes in the mesh can be selected, which adds them to the Node window on the righthand side of the screen. When selecting one or more nodes, the table view below can be opened. The data itself is displayed here, and columns can be added to the view. 

## Demo
Try our demo with Elasticsearch, Twitter and Bitcoin datasources [here](DEMO.md).


## Configuration
Go to the Marija installation directory and copy the sample configuration file.
``` 
cd go/src/github.com/dutchcoders/marija/
cp config.toml.sample config.toml 
```
If Elasticsearch is not installed, comment it out to prevent errors.
```
#[datasource]
#[datasource.elasticsearch]
#type="elasticsearch"
#url="http://127.0.0.1:9200/demo_index"',
#username=
#password=

[datasource.twitter]
type="twitter"
consumer_key=""
consumer_secret=""
token=""
token_secret=""

[datasource.blockchain]
type="blockchain"

[[logging]]
output = "stdout"
level = "debug"
```
More on configuring datasources [here](CONFIGURATION.md).

## Features

* works on local and remote datasources simultaneously
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

We're working towards a first version. 

* analyze data at realtime
* create specialized tools based on Marija for graphing for example packet traffic flows. 
* see issue list for features and bugs

## Contribute

Contributions are welcome.

### Setup your Marija Github Repository

Fork Marija upstream source repository to your own personal repository. Copy the URL for marija from your personal github repo (you will need it for the git clone command below).

```sh
$ mkdir -p $GOPATH/src/github.com/marija
$ cd $GOPATH/src/github.com/marija
$ git clone <paste saved URL for personal forked marija repo>
$ cd marija
```

###  Developer Guidelines
``Marija`` community welcomes your contribution. To make the process as seamless as possible, we ask for the following:
* Go ahead and fork the project and make your changes. We encourage pull requests to discuss code changes.
    - Fork it
    - Create your feature branch (git checkout -b my-new-feature)
    - Commit your changes (git commit -am 'Add some feature')
    - Push to the branch (git push origin my-new-feature)
    - Create new Pull Request

* If you have additional dependencies for ``Marija``, ``Marija`` manages its dependencies using [govendor](https://github.com/kardianos/govendor)
    - Run `go get foo/bar`
    - Edit your code to import foo/bar
    - Run `make pkg-add PKG=foo/bar` from top-level directory

* If you have dependencies for ``Marija`` which needs to be removed
    - Edit your code to not import foo/bar
    - Run `make pkg-remove PKG=foo/bar` from top-level directory

* When you're ready to create a pull request, be sure to:
    - Have test cases for the new code. If you have questions about how to do it, please ask in your pull request.
    - Run `make verifiers`
    - Squash your commits into a single commit. `git rebase -i`. It's okay to force update your pull request.
    - Make sure `go test -race ./...` and `go build` completes.

* Read [Effective Go](https://github.com/golang/go/wiki/CodeReviewComments) article from Golang project
    - `Marija` project is fully conformant with Golang style
    - if you happen to observe offending code, please feel free to send a pull request

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

**Kevin Hoogerwerf**
- <https://keybase.io/kevinh>

## Copyright and license

Code and documentation copyright 2016 Remco Verhoef.

Code released under [the Apache license](LICENSE).

