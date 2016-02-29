package main

import (
	"os"
	"testing"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/cover"
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
		t.Errorf("readFile(): got error: %s", err)
	}
	if len(file) == 0 {
		t.Errorf("readFile(): file empty")
	}
	if string(file[:12]) != "package main" {
		t.Errorf("readFile(): failed read first line")
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
	tmpFileName, err := getTempFileName()
	if err != nil {
		t.Errorf("getTempFileName() got error")
	}
	defer os.RemoveAll(tmpFileName)

	if len(tmpFileName) == 0 {
		t.Errorf("getTempFileName() failed")
	}

	// on RO-dir
	cwd, _ := os.Getwd()
	os.Chdir("/")
	_, err = getTempFileName()
	if err == nil {
		t.Errorf("getTempFileName() not got error")
	}
	os.Chdir(cwd)
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

func Test_getColorHeader(t *testing.T) {
	result := getColorHeader("filename.go")
	expected := ansi.ColorCode("yellow") + "filename.go" + ansi.ColorCode("reset") + "\n" +
		ansi.ColorCode("black+h") + "~~~~~~~~~~~" + ansi.ColorCode("reset") + "\n"

	if result != expected {
		t.Errorf("1. getColorHeader() failed")
	}
}

func Test_getCoverForFile(t *testing.T) {
	fileProfile := &cover.Profile{
		FileName: "filename.go",
		Mode:     "count",
		Blocks: []cover.ProfileBlock{
			{
				StartLine: 2,
				StartCol:  5,
				EndLine:   2,
				EndCol:    10,
				NumStmt:   1,
				Count:     1,
			},
		},
	}
	fileContent := []byte("1 line\n123 green 456\n3 line red and other")

	coloredBytes := getCoverForFile(fileProfile, fileContent, false)
	expectOut := getColorHeader("filename.go") +
		"1 line\n" +
		"123 " + ansi.ColorCode("green") + "green" + ansi.ColorCode("reset") + " 456\n" +
		"3 line red and other\n"
	if string(coloredBytes) != expectOut {
		t.Errorf("1. getCoverForFile() failed")
	}

	// with red blocks
	fileProfile.Blocks = append(fileProfile.Blocks,
		cover.ProfileBlock{
			StartLine: 3,
			StartCol:  8,
			EndLine:   3,
			EndCol:    11,
			NumStmt:   0,
			Count:     0,
		},
	)
	coloredBytes = getCoverForFile(fileProfile, fileContent, false)
	expectOut = getColorHeader("filename.go") +
		"1 line\n" +
		"123 " + ansi.ColorCode("green") + "green" + ansi.ColorCode("reset") + " 456\n" +
		"3 line " + ansi.ColorCode("red") + "red" + ansi.ColorCode("reset") + " and other\n"
	if string(coloredBytes) != expectOut {
		t.Errorf("1. getCoverForFile() failed")
	}

	// 256 colors
	coloredBytes = getCoverForFile(fileProfile, fileContent, true)
	expectOut = getColorHeader("filename.go") +
		"1 line\n" +
		"123 " + ansi.ColorCode("48") + "green" + ansi.ColorCode("reset") + " 456\n" +
		"3 line " + ansi.ColorCode("red") + "red" + ansi.ColorCode("reset") + " and other\n"
	if string(coloredBytes) != expectOut {
		t.Errorf("1. getCoverForFile() failed")
	}
}

func Test_runGoTest(t *testing.T) {
	err := runGoTest("./not exists dir", "", true)
	if err == nil {
		t.Errorf("runGoTest() error failed")
	}
}
