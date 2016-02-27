/*
Get all colors for 255-colors terminal:
	gommand 'for i := 0; i < 256; i++ {fmt.Println(i, ansi.ColorCode(strconv.Itoa(i)) + "String" + ansi.ColorCode("reset"))}'
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/cover"
)

func readFile(fileName string) (result []byte, err error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer fileReader.Close()

	result, err = ioutil.ReadAll(fileReader)
	return result, err
}

func main() {
	coverFile := "coverage.out"

	execTestCover := exec.Command("go", "test", "-coverprofile="+coverFile, "-covermode=count")
	err := execTestCover.Run()
	if err != nil {
		log.Fatal(err)
	}

	coverProfile, err := cover.ParseProfiles(coverFile)
	if err != nil {
		log.Fatal(err)
	}

	for _, fileProfile := range coverProfile {
		fileName := ""
		if strings.HasPrefix(fileProfile.FileName, "_") {
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
				fmt.Print(string(fileBytes[curOffset:boundary.Offset]))
			}
			switch {
			case boundary.Start && boundary.Count > 0:
				fmt.Print(ansi.ColorCode("green"))
			case boundary.Start && boundary.Count == 0:
				fmt.Print(ansi.ColorCode("red"))
			case !boundary.Start:
				fmt.Print(ansi.ColorCode("reset"))
			}

			curOffset = boundary.Offset
		}
		if curOffset < len(fileBytes) {
			fmt.Print(string(fileBytes[curOffset:len(fileBytes)]))
		}
	}

	err = os.Remove(coverFile)
	if err != nil {
		log.Fatal(err)
	}
}
