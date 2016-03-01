go-carpet - show test coverage for Go source files
==================================================

[![GoDoc](https://godoc.org/github.com/msoap/go-carpet?status.svg)](https://godoc.org/github.com/msoap/go-carpet)
[![Build Status](https://travis-ci.org/msoap/go-carpet.svg?branch=master)](https://travis-ci.org/msoap/go-carpet)
[![Coverage Status](https://coveralls.io/repos/github/msoap/go-carpet/badge.svg?branch=master)](https://coveralls.io/github/msoap/go-carpet?branch=master)
[![Homebrew formula exists](https://img.shields.io/badge/homebrew-🍺-d7af72.svg)](https://github.com/msoap/go-carpet#install)
[![Report Card](https://goreportcard.com/badge/github.com/msoap/go-carpet)](https://goreportcard.com/report/github.com/msoap/go-carpet)

To view the test coverage in the terminal, just run go-carpet.

It works not only in the directory GOPATH. And it works recursively for multiple packages.

With -256colors option, shades of green indicate the level of coverage.

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
		-file string - comma separated list of files to test (defualt: all)

For view in less, use `-R` option:

    go-carpet | less -R

###Screenshot
<img width="664" alt="go-carpet-screenshot" src="https://cloud.githubusercontent.com/assets/844117/13379229/967db8e0-de28-11e5-8b3f-cba51af3f3bc.png">

See also
--------

  * [The cover story (blog.golang.org)](https://blog.golang.org/cover)
  * [gocover.io - simple Go test coverage service](https://gocover.io)
  * [coveralls.io - test coverage service](https://coveralls.io)
  * [Package cover](https://godoc.org/golang.org/x/tools/cover)
