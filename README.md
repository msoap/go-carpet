go-carpet - show test coverage for Go source files
==================================================

[![GoDoc](https://godoc.org/github.com/msoap/go-carpet?status.svg)](https://godoc.org/github.com/msoap/go-carpet)
[![Build Status](https://travis-ci.org/msoap/go-carpet.svg?branch=master)](https://travis-ci.org/msoap/go-carpet)
[![Coverage Status](https://coveralls.io/repos/github/msoap/go-carpet/badge.svg?branch=master)](https://coveralls.io/github/msoap/go-carpet?branch=master)
[![Report Card](https://goreportcard.com/badge/github.com/msoap/go-carpet)](https://goreportcard.com/report/github.com/msoap/go-carpet)

To view the test coverage in the terminal, just run go-carpet.

It works not only in the directory GOPATH. And it works recursively for multiple packages.

Install
-------

    go get -u github.com/msoap/go-carpet

Usage
-----

	go-carpet [-options] [paths]
	options:
		-256colors - use more colors on 256-color terminal
		-file string - comma separated list of files to test (defualt: all)

###screenshot

See also
--------

  * [The cover story (blog.golang.org)](https://blog.golang.org/cover)
  * [gocover.io - simple Go test coverage service](https://gocover.io)
  * [coveralls.io - test coverage service](https://coveralls.io)
  * [Package cover](https://godoc.org/golang.org/x/tools/cover)
