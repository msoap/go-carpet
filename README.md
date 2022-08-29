go-carpet - show test coverage for Go source files
==================================================

[![Go Reference](https://pkg.go.dev/badge/github.com/msoap/go-carpet.svg)](https://pkg.go.dev/github.com/msoap/go-carpet)
[![Go](https://github.com/msoap/go-carpet/actions/workflows/go.yml/badge.svg)](https://github.com/msoap/go-carpet/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/msoap/go-carpet/badge.svg?branch=master)](https://coveralls.io/github/msoap/go-carpet?branch=master)
[![Report Card](https://goreportcard.com/badge/github.com/msoap/go-carpet)](https://goreportcard.com/report/github.com/msoap/go-carpet)
[![Homebrew formula exists](https://img.shields.io/badge/homebrew-🍺-d7af72.svg)](https://github.com/msoap/go-carpet#install)

To view the test coverage in the terminal, just run `go-carpet`.

It works outside of the `GOPATH` directory. And it works recursively for multiple packages.

With `-256colors` option, shades of green indicate the level of coverage.

By default skip vendor directories (Godeps,vendor), otherwise use `-include-vendor` option.

The `-mincov` option allows you to specify a coverage threshold to limit the files to be displayed.

Usage
-----

    usage: go-carpet [options] [paths]
      -256colors
        	use more colors on 256-color terminal (indicate the level of coverage)
      -args string
        	pass additional arguments for go test
      -file string
        	comma-separated list of files to test (default: all)
      -func string
        	comma-separated functions list (default: all functions)
      -include-vendor
        	include vendor directories for show coverage (Godeps, vendor)
      -mincov float
        	coverage threshold of the file to be displayed (in percent) (default 100)
      -summary
        	only show summary for each file
      -version
        	get version

For view coverage in less, use `-R` option:

    go-carpet | less -R

Install
-------

From source:

    go install github.com/msoap/go-carpet@latest

Download binaries from: [releases](https://github.com/msoap/go-carpet/releases) (OS X/Linux/Windows)

Install from homebrew (OS X):

    brew tap msoap/tools
    brew install go-carpet
    # update:
    brew upgrade go-carpet

### Screenshot

<img width="662" alt="screen shot 2016-03-06" src="https://cloud.githubusercontent.com/assets/844117/13554107/e6c7c82a-e3a7-11e5-82d6-3481f1fead11.png">

See also
--------

  * [blog.golang.org](https://blog.golang.org/cover) - the cover story
  * [gocover.io](https://gocover.io) - simple Go test coverage service
  * [coveralls.io](https://coveralls.io) - test coverage service
  * [package cover](https://godoc.org/golang.org/x/tools/cover) - golang.org/x/tools/cover
  * [gotests](https://github.com/cweill/gotests) - Go commandline tool that generates table driven tests
  * [docker-golang-checks](https://github.com/msoap/docker-golang-checks) - Go-code checks Docker image
