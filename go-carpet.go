/*
Get all colors for 255-colors terminal:
	gommand 'for i := 0; i < 256; i++ {fmt.Println(i, ansi.ColorCode(strconv.Itoa(i)) + "String" + ansi.ColorCode("reset"))}'
*/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/cover"
)

func getDirsWithTests(root string) []string {
	dirs := map[string]struct{}{}
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "_test.go") {
			dirs["./"+filepath.Dir(path)] = struct{}{}
		}
		return nil
	})

	result := make([]string, 0, len(dirs))
	for dir := range dirs {
		result = append(result, dir)
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

func printCoverForDir(path, coverFileName string, stdOut io.Writer) {
	err := exec.Command("go", "test", "-coverprofile="+coverFileName, "-covermode=count", path).Run()
	if err != nil {
		log.Fatal(err)
	}

	coverProfile, err := cover.ParseProfiles(coverFileName)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}

		boundaries := fileProfile.Boundaries(fileBytes)
		curOffset := 0
		for _, boundary := range boundaries {
			if boundary.Offset > curOffset {
				stdOut.Write(fileBytes[curOffset:boundary.Offset])
			}
			switch {
			case boundary.Start && boundary.Count > 0:
				stdOut.Write([]byte(ansi.ColorCode("green")))
			case boundary.Start && boundary.Count == 0:
				stdOut.Write([]byte(ansi.ColorCode("red")))
			case !boundary.Start:
				stdOut.Write([]byte(ansi.ColorCode("reset")))
			}

			curOffset = boundary.Offset
		}
		if curOffset < len(fileBytes) {
			stdOut.Write(fileBytes[curOffset:len(fileBytes)])
		}
	}
}

func main() {
	tmpDir, err := ioutil.TempDir("", "go-carpet-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	coverFileName := filepath.Join(tmpDir, "coverage.out")
	stdOut := getColorWriter()

	for _, path := range getDirsWithTests(".") {
		printCoverForDir(path, coverFileName, stdOut)
	}
}
