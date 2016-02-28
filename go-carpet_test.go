package main

import "testing"

func Test_readFile(t *testing.T) {
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

func Test_getDirsWithTests(t *testing.T) {
	dirs := getDirsWithTests(".")
	if len(dirs) == 0 {
		t.Errorf("Dir list is empty")
	}
	dirs = getDirsWithTests()
	if len(dirs) == 0 {
		t.Errorf("Dir list is empty")
	}
	dirs = getDirsWithTests(".", ".")
	if len(dirs) != 1 {
		t.Errorf("The same directory failed")
	}
}
