package main

import (
	"os"
	"testing"
)

func assertDontPanic(t *testing.T, fn func(), name string) {
	defer func() {
		if recoverInfo := recover(); recoverInfo != nil {
			t.Errorf("The code panic: %s\npanic: %s", name, recoverInfo)
		}
	}()
	fn()
}

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

func Test_getTempFileName(t *testing.T) {
	tmpFileName := getTempFileName()
	defer os.RemoveAll(tmpFileName)

	if len(tmpFileName) == 0 {
		t.Errorf("getTempFileName() failed")
	}
}

func Test_isSliceInString(t *testing.T) {
	testData := []struct {
		src    string
		slice  []string
		result bool
	}{
		{
			src:    "one/file.go",
			slice:  []string{"one.go", "file.go"},
			result: true,
		},
		{
			src:    "path/path/file.go",
			slice:  []string{"one.go", "path/file.go"},
			result: true,
		},
		{
			src:    "one/file.go",
			slice:  []string{"one.go", "two.go"},
			result: false,
		},
		{
			src:    "one/file.go",
			slice:  []string{},
			result: false,
		},
	}

	for i, item := range testData {
		result := isSliceInString(item.src, item.slice)
		if result != item.result {
			t.Errorf("\n%d.\nexpected: %v\nreal    :%v", i, item.result, result)
		}
	}
}

func Test_getShadeOfGreen(t *testing.T) {
	testData := []struct {
		normCover float64
		result    string
	}{
		{
			normCover: 0,
			result:    "29",
		},
		{
			normCover: 1,
			result:    "51",
		},
		{
			normCover: 0.99999,
			result:    "51",
		},
		{
			normCover: 0.5,
			result:    "40",
		},
		{
			normCover: -1,
			result:    "29",
		},
		{
			normCover: 11,
			result:    "51",
		},
		{
			normCover: 100500,
			result:    "51",
		},
	}

	for i, item := range testData {
		result := getShadeOfGreen(item.normCover)
		if result != item.result {
			t.Errorf("\n%d.\nexpected: %v\nreal    : %v", i, item.result, result)
		}
	}
}

func Test_getColorWriter(t *testing.T) {
	assertDontPanic(t, func() { getColorWriter() }, "getColorWriter()")
}
