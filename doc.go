/*
go-carpet - show test coverage for Go source files

It works not only in the directory GOPATH. And it works recursively for multiple packages.
With -256colors option, shades of green indicate the level of coverage.

Install/update:

	go get -u github.com/msoap/go-carpet
	ln -s $GOPATH/bin/go-carpet ~/bin/go-carpet

Usage:

	go-carpet [-options] [paths]
	options:
	    -256colors - use more colors on 256-color terminal (indicate the level of coverage)
	    -args - pass additional arguments for go test (for example "-short" or "-i -timeout t")
	    -file string - comma-separated list of files to test (default: all)
	    -func string - comma-separated functions list (default: all functions)
	    -include-vendor - include vendor directories for show coverage (Godeps, vendor)
	    -summary - only show summary for each file
	    -version - get version

Source: https://github.com/msoap/go-carpet
*/
package main
