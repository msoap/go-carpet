//go:build !windows
// +build !windows

package main

import (
	"io"
	"os"
)

func getColorWriter() io.Writer {
	return (io.Writer)(os.Stdout)
}
