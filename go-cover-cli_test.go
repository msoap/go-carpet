package main

import "testing"

func Test_readFileByLines(t *testing.T) {
	file, err := readFileByLines("go-cover-cli_test.go")
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	if len(file) == 0 {
		t.Errorf("File empty")
	}

	_, err = readFileByLines("dont exists file")
	if err == nil {
		t.Errorf("File exists error:")
	}
}
