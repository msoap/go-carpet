package main

import (
	"os"
	"path/filepath"
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

// usage:
//     defer testChdir(t, "/path")()
func testChdir(t *testing.T, dir string) func() {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	return func() {
		err := os.Chdir(cwd)
		if err != nil {
			t.Fatal(err)
		}
	}
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
	dirs, err := getDirsWithTests(false, ".")
	if len(dirs) == 0 || err != nil {
		t.Errorf("getDirsWithTests(): dir list is empty")
	}
	dirs, err = getDirsWithTests(false)
	if len(dirs) == 0 || err != nil {
		t.Errorf("getDirsWithTests(): dir list is empty")
	}
	dirs, err = getDirsWithTests(false, ".", ".")
	if len(dirs) != 1 || err != nil {
		t.Errorf("getDirsWithTests(): the same directory failed")
	}

	defer testChdir(t, "./testdata")()
	dirs, err = getDirsWithTests(false, ".")
	if len(dirs) != 1 || err != nil {
		t.Errorf("getDirsWithTests(): without vendor dirs")
	}

	dirs, err = getDirsWithTests(true, ".")
	if len(dirs) != 4 || err != nil {
		t.Errorf("getDirsWithTests(): with vendor dirs")
	}
}

func Test_getTempFileName(t *testing.T) {
	tmpFileName, err := getTempFileName()
	if err != nil {
		t.Errorf("getTempFileName() got error")
	}
	defer func() {
		err = os.RemoveAll(tmpFileName)
		if err != nil {
			t.Errorf("getTempFileName() RemoveAll failed")
		}
	}()

	if len(tmpFileName) == 0 {
		t.Errorf("getTempFileName() failed")
	}

	// on RO-dir
	defer testChdir(t, "/")()
	_, err = getTempFileName()
	if err == nil {
		t.Errorf("getTempFileName() not got error")
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
	result := getColorHeader("filename.go", true)
	expected := ansi.ColorCode("yellow") + "filename.go" + ansi.ColorCode("reset") + "\n" +
		ansi.ColorCode("black+h") + "~~~~~~~~~~~" + ansi.ColorCode("reset") + "\n"

	if result != expected {
		t.Errorf("1. getColorHeader() failed")
	}

	result = getColorHeader("filename.go", false)
	expected = ansi.ColorCode("yellow") + "filename.go" + ansi.ColorCode("reset") + "\n"

	if result != expected {
		t.Errorf("2. getColorHeader() failed")
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

	coloredBytes := getCoverForFile(fileProfile, fileContent, Config{colors256: false})
	expectOut := getColorHeader("filename.go - 100.0%", true) +
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
	coloredBytes = getCoverForFile(fileProfile, fileContent, Config{colors256: false})
	expectOut = getColorHeader("filename.go - 100.0%", true) +
		"1 line\n" +
		"123 " + ansi.ColorCode("green") + "green" + ansi.ColorCode("reset") + " 456\n" +
		"3 line " + ansi.ColorCode("red") + "red" + ansi.ColorCode("reset") + " and other\n"
	if string(coloredBytes) != expectOut {
		t.Errorf("2. getCoverForFile() failed")
	}

	// 256 colors
	coloredBytes = getCoverForFile(fileProfile, fileContent, Config{colors256: true})
	expectOut = getColorHeader("filename.go - 100.0%", true) +
		"1 line\n" +
		"123 " + ansi.ColorCode("48") + "green" + ansi.ColorCode("reset") + " 456\n" +
		"3 line " + ansi.ColorCode("red") + "red" + ansi.ColorCode("reset") + " and other\n"
	if string(coloredBytes) != expectOut {
		t.Errorf("3. getCoverForFile() failed")
	}
}

func Test_runGoTest(t *testing.T) {
	err := runGoTest("./not exists dir", "", []string{}, true)
	if err == nil {
		t.Errorf("runGoTest() error failed")
	}
}

func Test_guessAbsPathInGOPATH(t *testing.T) {
	GOPATH := ""
	absPath, err := guessAbsPathInGOPATH(GOPATH, "file.golang")
	if absPath != "" || err == nil {
		t.Errorf("1. guessAbsPathInGOPATH() empty GOPATH")
	}

	cwd, _ := os.Getwd()

	GOPATH = filepath.Join(cwd, "testdata")
	absPath, err = guessAbsPathInGOPATH(GOPATH, "file.golang")
	if err != nil {
		t.Errorf("2. guessAbsPathInGOPATH() error: %s", err)
	}
	if absPath != filepath.Join(cwd, "testdata", "src", "file.golang") {
		t.Errorf("3. guessAbsPathInGOPATH() empty GOPATH")
	}

	GOPATH = filepath.Join(cwd, "testdata") + string(os.PathListSeparator) + "/tmp"
	absPath, err = guessAbsPathInGOPATH(GOPATH, "file.golang")
	if err != nil {
		t.Errorf("4. guessAbsPathInGOPATH() error: %s", err)
	}
	if absPath != filepath.Join(cwd, "testdata", "src", "file.golang") {
		t.Errorf("5. guessAbsPathInGOPATH() empty GOPATH")
	}

	GOPATH = "/tmp" + string(os.PathListSeparator) + "/"
	absPath, err = guessAbsPathInGOPATH(GOPATH, "file.golang")
	if absPath != "" || err == nil {
		t.Errorf("6. guessAbsPathInGOPATH() file not in GOPATH")
	}
}

func Test_getStatForProfileBlocks(t *testing.T) {
	profileBlocks := []cover.ProfileBlock{
		{
			StartLine: 2,
			StartCol:  5,
			EndLine:   2,
			EndCol:    10,
			NumStmt:   1,
			Count:     1,
		},
	}

	stat := getStatForProfileBlocks(profileBlocks)
	if stat != 100.0 {
		t.Errorf("1. getStatForProfileBlocks() failed")
	}

	profileBlocks = append(profileBlocks,
		cover.ProfileBlock{
			StartLine: 3,
			StartCol:  5,
			EndLine:   3,
			EndCol:    10,
			NumStmt:   1,
			Count:     0,
		},
	)
	stat = getStatForProfileBlocks(profileBlocks)
	if stat != 50.0 {
		t.Errorf("2. getStatForProfileBlocks() failed")
	}

	profileBlocks = append(profileBlocks,
		cover.ProfileBlock{
			StartLine: 4,
			StartCol:  5,
			EndLine:   4,
			EndCol:    10,
			NumStmt:   1,
			Count:     0,
		},
		cover.ProfileBlock{
			StartLine: 4,
			StartCol:  5,
			EndLine:   4,
			EndCol:    10,
			NumStmt:   1,
			Count:     0,
		},
	)
	stat = getStatForProfileBlocks(profileBlocks)
	if stat != 25.0 {
		t.Errorf("3. getStatForProfileBlocks() failed")
	}
}
