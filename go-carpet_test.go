package main

import "testing"

func Test_readFileByLines(t *testing.T) {
	file, err := readFile("go-carpet_test.go")
	if err != nil {
		t.Errorf("Got error: %s", err)
	}
	if len(file) == 0 {
		t.Errorf("File empty")
	}

	_, err = readFile("dont exists file")
	if err == nil {
		t.Errorf("File exists error:")
	}
}
