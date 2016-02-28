/*
go-carpet - show test coverage for Go source files

It works not only in the directory GOPATH. And it works recursively for multiple packages.

Install/update:

	go get -u github.com/msoap/go-carpet
	ln -s $GOPATH/bin/go-carpet ~/bin/go-carpet

Usage:

	go-carpet [-options] [paths]
	options:
		-256colors - use more colors on 256-color terminal
		-file string - comma separated list of files to test (defualt: all)

Source: https://github.com/msoap/go-carpet

*/
package main
