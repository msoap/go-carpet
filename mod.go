package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/msoap/byline"
)

var goModFilename *string

func getGoModFilename() string {
	if goModFilename != nil {
		return *goModFilename
	}

	file := ""
	out, err := exec.Command("go", "env", "GOMOD").Output()
	if err != nil {
		log.Printf("failed to load 'go env GOMOD' content: %s", err)
		goModFilename = &file
		return ""
	}

	file = strings.TrimSpace(string(out))
	goModFilename = &file

	return file
}

func guessAbsPathInGoMod(relPath string) (string, error) {
	modFilename := getGoModFilename()
	if modFilename == "" {
		return "", errIsNotInGoMod
	}

	modFile, err := os.Open(modFilename)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := modFile.Close(); err != nil {
			log.Printf("failed to close %s file: %s", modFilename, err)
		}
	}()

	moduleName := ""
	if err := byline.NewReader(modFile).AWKMode(func(line string, fields []string, vars byline.AWKVars) (string, error) {
		if vars.NF == 2 && fields[0] == "module" && fields[1] != "" {
			moduleName = fields[1]
			return "", io.EOF
		}

		return "", nil
	}).Discard(); err != nil {
		return "", err
	}
	if moduleName == "" {
		return "", errIsNotInGoMod
	}

	absPath := path.Dir(modFilename) + strings.TrimPrefix(relPath, moduleName)
	if stat, err := os.Stat(absPath); err != nil {
		return "", err
	} else if !stat.Mode().IsRegular() {
		return "", fmt.Errorf("%s is not regular file", absPath)
	}

	return absPath, nil
}
