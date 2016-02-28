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

const usageMessage = `go-carpet - show test coverage for Go source files

usage: go-carpet [options] [paths]`

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

// isStringInSlice - one of the elements of the array contained in the string
func isSliceInString(src string, slice []string) bool {
	for _, dst := range slice {
		if strings.Contains(src, dst) {
			return true
		}
	}
	return false
}

func getShadeOfGreen(normCover float64) string {
	/*
		Get all colors for 255-colors terminal:
			gommand 'for i := 0; i < 256; i++ {fmt.Println(i, ansi.ColorCode(strconv.Itoa(i)) + "String" + ansi.ColorCode("reset"))}'
	*/
	var tenShadesOfGreen = []string{
		"29",
		"30",
		"34",
		"36",
		"40",
		"42",
		"46",
		"48",
		"50",
		"51",
	}
	if normCover < 0 {
		normCover = 0
	}
	if normCover > 1 {
		normCover = 1
	}
	index := int((normCover - 0.00001) * float64(len(tenShadesOfGreen)))
	return tenShadesOfGreen[index]
}

func getCoverForDir(path string, coverFileName string, filesFilter []string, colors256 bool) (result []byte, err error) {
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

		if len(filesFilter) > 0 && !isSliceInString(fileName, filesFilter) {
			continue
		}

		fileBytes, err := readFile(fileName)
		if err != nil {
			return result, err
		}

		fileNameDisplay := fileProfile.FileName

		result = append(result,
			[]byte(ansi.ColorCode("yellow")+
				fileNameDisplay+ansi.ColorCode("reset")+"\n"+
				ansi.ColorCode("black+h")+
				strings.Repeat("~", len(fileNameDisplay))+
				ansi.ColorCode("reset")+"\n",
			)...,
		)

		boundaries := fileProfile.Boundaries(fileBytes)
		curOffset := 0
		for _, boundary := range boundaries {
			if boundary.Offset > curOffset {
				result = append(result, fileBytes[curOffset:boundary.Offset]...)
			}
			switch {
			case boundary.Start && boundary.Count > 0:
				coverColor := ansi.ColorCode("green")
				if colors256 {
					coverColor = ansi.ColorCode(getShadeOfGreen(boundary.Norm))
				}
				result = append(result, []byte(coverColor)...)
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
	filesFilter, colors256 := "", false
	flag.StringVar(&filesFilter, "file", "", "comma separated list of files to test (defualt: all)")
	flag.BoolVar(&colors256, "256colors", false, "use more colors on 256-color terminal")
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
		coverInBytes, err := getCoverForDir(path, coverFileName, strings.Split(filesFilter, ","), colors256)
		if err != nil {
			log.Print(err)
		}
		stdOut.Write(coverInBytes)
	}
}
