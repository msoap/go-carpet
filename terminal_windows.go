package main

import (
	"io"

	"github.com/mattn/go-colorable"
)

func getColorWriter() io.Writer {
	return colorable.NewColorableStdout()
}
