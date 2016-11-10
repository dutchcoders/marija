![](https://github.com/dutchcoders/marija-screenshots/blob/master/marija.png?raw=true)

# Marija [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/dutchcoders/marija?utm_source=badge&utm_medium=badge&utm_campaign=&utm_campaign=pr-badge&utm_content=badge) [![Go Report Card](https://goreportcard.com/badge/dutchcoders/marija)](https://goreportcard.com/report/dutchcoders/marija)

Marija is a data exploration and visualisation tool for (un)structured Elasticsearch data. Using Marija you'll be able to see relations 
between data of different datasources without any modifications to your data or index.

Currently Marija is being used to identify related spamruns, but can be used for all kind of different data sets.

# Screenshot

![](https://github.com/dutchcoders/marija-screenshots/blob/master/Screen%20Shot%202016-11-07%20at%2017.25.21.png?raw=true)

## Install

Currently installation is only supported using source.

### Install Golang

If you do not have a working Golang environment setup please follow [Golang Installation Guide](https://golang.org/doc/install).

### Install Marija

Installation of Marija is easy.

```
$ go get github.com/dutchcoders/marija
$ marija
```

## Usage

There are a few steps you need to take before you can start.

* add your elasticsearch server to the configuration
* use the cloud icon to retrieve the indexes
* enable the index(es) you want to search in using the eye icon
* add the fields you want to use as nodes

* additionally you can add the date field you want to use for the histogram
* and add some normalizations (eg removing part of the identifier) using regular expressions

You're all setup now, just type your queries.

## Features

* work on multiple servers and indexes at the same time
* different fields can be used as node identifier
* identifiers can be normalized through normalization regular expressions
* each field will have its own icon
* query indexes using elasticsearch queries like your used to do
* histogram view to identify nodes in time
* select and delete nodes
* select related nodes, deselect all but selected nodes
* zoom and move nodes
* navigate through selected data using the tableview

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

$ mkdir -p $GOPATH/src/github.com/marija
$ cd $GOPATH/src/github.com/marija
$ git clone <paste saved URL for personal forked marija repo>
$ cd marija

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

## Copyright and license

Code and documentation copyright 2016 Remco Verhoef.
Code released under [the Apache license](LICENSE).

