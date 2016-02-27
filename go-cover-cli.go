package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/tools/cover"
)

func main() {
	coverFile := "coverage.out"

	execTestCover := exec.Command("go", "test", "-coverprofile="+coverFile)
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

		fmt.Println(fileName + "\n---\n")
	}

	err = os.Remove(coverFile)
	if err != nil {
		log.Fatal(err)
	}
}
