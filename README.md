# Marija [![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/dutchcoders/marija?utm_source=badge&utm_medium=badge&utm_campaign=&utm_campaign=pr-badge&utm_content=badge) [![Go Report Card](https://goreportcard.com/badge/dutchcoders/marija)](https://goreportcard.com/report/dutchcoders/marija)

[![Join the chat at https://gitter.im/dutchcoders/marija](https://badges.gitter.im/dutchcoders/marija.svg)](https://gitter.im/dutchcoders/marija?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Marija is a graphing solution for (un)structured Elasticsearch data. Using Marija you'll be able to see connections 
between data of different indexes and datasources without any modifications to your data or index.

Currently Marija is being used to identify related spamruns, but can be used for all kind of different set, but can be used for all kind of different sets

# Screenshot

## Install

Installation of Marija is easy.

```
$ go get github.com/dutchcoders/marija
$ marija
```

## Usage

There are a few steps you need to take before you can start.

* add the indexes in the configuration
* add the correct fields, those fields are being used as node identifier
* now you can enter your queries, and analyse the data


## Roadmap

We're working towards a first version. From there we'll work on the following:

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

**Kevin Hoogerwerf**

## Copyright and license

Code and documentation copyright 206 Remco Verhoef.
Code released under [the Apache license](LICENSE).

