/*
Get all colors for 255-colors terminal:
	gommand 'for i := 0; i < 256; i++ {fmt.Println(i, ansi.ColorCode(strconv.Itoa(i)) + "String" + ansi.ColorCode("reset"))}'
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/cover"
)

const usageMessage = `go-carpet - show coverage for Go source files

usage: go-carpet [dirs]`

func getDirsWithTests(roots ...string) []string {
	if len(roots) == 0 {
		roots = []string{"."}
	}

	dirs := map[string]struct{}{}
	for _, root := range roots {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(path, "_test.go") {
				dirs[filepath.Dir(path)] = struct{}{}
			}
			return nil
		})
	}

	result := make([]string, 0, len(dirs))
	for dir := range dirs {
		result = append(result, "./"+dir)
	}
	return result
}

func readFile(fileName string) (result []byte, err error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer fileReader.Close()

	result, err = ioutil.ReadAll(fileReader)
	return result, err
}

func getCoverForDir(path string, coverFileName string) (result []byte, err error) {
	osExec := exec.Command("go", "test", "-coverprofile="+coverFileName, "-covermode=count", path)
	osExec.Stderr = os.Stderr
	err = osExec.Run()
	if err != nil {
		return result, err
	}

	coverProfile, err := cover.ParseProfiles(coverFileName)
	if err != nil {
		return result, err
	}

	for _, fileProfile := range coverProfile {
		fileName := ""
		if strings.HasPrefix(fileProfile.FileName, "_") {
			// absolute path
			fileName = strings.TrimLeft(fileProfile.FileName, "_")
		} else {
			// file in GOPATH
			fileName = os.Getenv("GOPATH") + "/src/" + fileProfile.FileName
		}
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			fmt.Printf("File '%s' is not exists\n", fileName)
			continue
		}

		fileBytes, err := readFile(fileName)
		if err != nil {
			return result, err
		}

		fileNameDisplay := fileProfile.FileName

		result = append(result, []byte(ansi.ColorCode("yellow")+fileNameDisplay+ansi.ColorCode("reset")+"\n"+
			ansi.ColorCode("black+h")+strings.Repeat("~", len(fileNameDisplay))+ansi.ColorCode("reset")+"\n")...)

		boundaries := fileProfile.Boundaries(fileBytes)
		curOffset := 0
		for _, boundary := range boundaries {
			if boundary.Offset > curOffset {
				result = append(result, fileBytes[curOffset:boundary.Offset]...)
			}
			switch {
			case boundary.Start && boundary.Count > 0:
				result = append(result, []byte(ansi.ColorCode("green"))...)
			case boundary.Start && boundary.Count == 0:
				result = append(result, []byte(ansi.ColorCode("red"))...)
			case !boundary.Start:
				result = append(result, []byte(ansi.ColorCode("reset"))...)
			}

			curOffset = boundary.Offset
		}
		if curOffset < len(fileBytes) {
			result = append(result, fileBytes[curOffset:len(fileBytes)]...)
		}
		result = append(result, []byte("\n")...)
	}

	return result, err
}

func getTempFileName() string {
	tmpFile, err := ioutil.TempFile(".", "go-carpet-coverage-out-")
	if err != nil {
		log.Fatal(err)
	}
	tmpFile.Close()

	return tmpFile.Name()
}

func main() {
	flag.Usage = func() {
		fmt.Println(usageMessage)
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	testDirs := flag.Args()

	coverFileName := getTempFileName()
	defer os.RemoveAll(coverFileName)
	stdOut := getColorWriter()

	if len(testDirs) > 0 {
		testDirs = getDirsWithTests(testDirs...)
	} else {
		testDirs = getDirsWithTests(".")
	}
	for _, path := range testDirs {
		coverInBytes, err := getCoverForDir(path, coverFileName)
		if err != nil {
			log.Print(err)
		}
		stdOut.Write(coverInBytes)
	}
}
