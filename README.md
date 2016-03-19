go-carpet - show test coverage for Go source files
==================================================

[![GoDoc](https://godoc.org/github.com/msoap/go-carpet?status.svg)](https://godoc.org/github.com/msoap/go-carpet)
[![Build Status](https://travis-ci.org/msoap/go-carpet.svg?branch=master)](https://travis-ci.org/msoap/go-carpet)
[![Coverage Status](https://coveralls.io/repos/github/msoap/go-carpet/badge.svg?branch=master)](https://coveralls.io/github/msoap/go-carpet?branch=master)
[![Homebrew formula exists](https://img.shields.io/badge/homebrew-🍺-d7af72.svg)](https://github.com/msoap/go-carpet#install)
[![Report Card](https://goreportcard.com/badge/github.com/msoap/go-carpet)](https://goreportcard.com/report/github.com/msoap/go-carpet)

To view the test coverage in the terminal, just run go-carpet.

It works outside of the directory GOPATH. And it works recursively for multiple packages.

With -256colors option, shades of green indicate the level of coverage.

By default skip vendor directories (Godeps,vendor), otherwise use -include-vendor option.

Install
-------

From source:

    go get -u github.com/msoap/go-carpet
    ln -s $GOPATH/bin/go-carpet /usr/local/bin/go-carpet

Download binaries from: [releases](https://github.com/msoap/go-carpet/releases) (OS X/Linux/Windows)

Or install from homebrew (OS X):

    brew tap msoap/tools
    brew install go-carpet
    # update:
    brew update; brew upgrade go-carpet

Usage
-----

    go-carpet [-options] [paths]
    options:
        -256colors - use more colors on 256-color terminal (indicate the level of coverage)
        -file string - comma separated list of files to test (default: all)
        -include-vendor - include vendor directories for show coverage (Godeps, vendor)

For view in less, use `-R` option:

    go-carpet | less -R

###Screenshot
<img width="662" alt="screen shot 2016-03-06" src="https://cloud.githubusercontent.com/assets/844117/13554107/e6c7c82a-e3a7-11e5-82d6-3481f1fead11.png">

TODO
----

  * option `-func` for filter by functions

See also
--------

  * [blog.golang.org](https://blog.golang.org/cover) - the cover story
  * [gocover.io](https://gocover.io) - simple Go test coverage service
  * [coveralls.io](https://coveralls.io) - test coverage service
  * [package cover](https://godoc.org/golang.org/x/tools/cover) - golang.org/x/tools/cover
  * [gotests](https://github.com/cweill/gotests) - Go commandline tool that generates table driven tests
