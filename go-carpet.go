package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/cover"
)

const usageMessage = `go-carpet - show test coverage for Go source files

usage: go-carpet [options] [paths]`

var (
	reNewLine = regexp.MustCompile("\n")
	// vendors directories for skip
	vendorDirs = []string{"Godeps", "vendor", ".vendor", "_vendor"}

	// directories for skip
	skipDirs = []string{"testdata"}
)

func getDirsWithTests(includeVendor bool, roots ...string) []string {
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
		if !includeVendor && isSliceInStringPrefix(dir, vendorDirs) || isSliceInStringPrefix(dir, skipDirs) {
			continue
		}
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

// isSliceInString - one of the elements of the array contained in the string
func isSliceInString(src string, slice []string) bool {
	for _, dst := range slice {
		if strings.Contains(src, dst) {
			return true
		}
	}
	return false
}

// isSliceInStringPrefix - one of the elements of the array is are prefix in the string
func isSliceInStringPrefix(src string, slice []string) bool {
	for _, dst := range slice {
		if strings.HasPrefix(src, dst) {
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

func runGoTest(path string, coverFileName string, hideStderr bool) error {
	osExec := exec.Command("go", "test", "-coverprofile="+coverFileName, "-covermode=count", path)
	if !hideStderr {
		osExec.Stderr = os.Stderr
	}
	output, err := osExec.Output()
	if err != nil {
		fmt.Print(string(output))
		return err
	}

	return nil
}

func guessAbsPathInGOPATH(GOPATH, relPath string) (absPath string, err error) {
	if GOPATH == "" {
		return "", fmt.Errorf("GOPATH is not set")
	}

	gopathChunks := strings.Split(GOPATH, string(os.PathListSeparator))
	for _, gopathChunk := range gopathChunks {
		guessAbsPath := filepath.Join(gopathChunk, "src", relPath)
		if _, err = os.Stat(guessAbsPath); err == nil {
			absPath = guessAbsPath
			break
		}
	}

	if absPath == "" {
		return "", fmt.Errorf("File '%s' not found in GOPATH", relPath)
	}
	return absPath, err
}

func getCoverForDir(path string, coverFileName string, filesFilter []string, colors256 bool) (result []byte, profileBlocks []cover.ProfileBlock, err error) {
	coverProfile, err := cover.ParseProfiles(coverFileName)
	if err != nil {
		return result, profileBlocks, err
	}

	for _, fileProfile := range coverProfile {
		fileName := ""
		if strings.HasPrefix(fileProfile.FileName, "_") {
			// absolute path (or relative in tests)
			if runtime.GOOS != "windows" {
				fileName = strings.TrimLeft(fileProfile.FileName, "_")
			} else {
				// "_\C_\Users\..." -> "C:\Users\..."
				fileName = regexp.MustCompile(`^_\\([A-Z])_`).ReplaceAllString(fileProfile.FileName, "$1:")
			}
		} else {
			// file in one dir in GOPATH
			fileName, err = guessAbsPathInGOPATH(os.Getenv("GOPATH"), fileProfile.FileName)
			if err != nil {
				return result, profileBlocks, err
			}
		}

		if len(filesFilter) > 0 && !isSliceInString(fileName, filesFilter) {
			continue
		}

		fileBytes, err := readFile(fileName)
		if err != nil {
			return result, profileBlocks, err
		}

		result = append(result, getCoverForFile(fileProfile, fileBytes, colors256)...)
		profileBlocks = append(profileBlocks, fileProfile.Blocks...)
	}

	return result, profileBlocks, err
}

func getColorHeader(header string, addUnderiline bool) string {
	result := ansi.ColorCode("yellow") +
		header + ansi.ColorCode("reset") + "\n"

	if addUnderiline {
		result += ansi.ColorCode("black+h") +
			strings.Repeat("~", len(header)) +
			ansi.ColorCode("reset") + "\n"
	}

	return result
}

// algorithms from Go-sources:
//   src/cmd/cover/html.go::percentCovered()
//   src/testing/cover.go::coverReport()
func getStatForProfileBlocks(fileProfileBlocks []cover.ProfileBlock) (stat float64) {
	var total, covered int64
	for _, profileBlock := range fileProfileBlocks {
		total += int64(profileBlock.NumStmt)
		if profileBlock.Count > 0 {
			covered += int64(profileBlock.NumStmt)
		}
	}
	if total > 0 {
		stat = float64(covered) / float64(total) * 100.0
	}

	return stat
}

func getCoverForFile(fileProfile *cover.Profile, fileBytes []byte, colors256 bool) (result []byte) {
	stat := getStatForProfileBlocks(fileProfile.Blocks)
	fileNameDisplay := fmt.Sprintf("%s - %.1f%%", strings.TrimLeft(fileProfile.FileName, "_"), stat)
	result = append(result, []byte(getColorHeader(fileNameDisplay, true))...)

	boundaries := fileProfile.Boundaries(fileBytes)
	curOffset := 0
	coverColor := ""

	for _, boundary := range boundaries {
		if boundary.Offset > curOffset {
			nextChunk := fileBytes[curOffset:boundary.Offset]
			// Add ansi color code in begin of each line (this fixed view in "less -R")
			if coverColor != "" && coverColor != ansi.ColorCode("reset") {
				nextChunk = reNewLine.ReplaceAllLiteral(nextChunk, []byte(ansi.ColorCode("reset")+"\n"+coverColor))
			}
			result = append(result, nextChunk...)
		}

		switch {
		case boundary.Start && boundary.Count > 0:
			coverColor = ansi.ColorCode("green")
			if colors256 {
				coverColor = ansi.ColorCode(getShadeOfGreen(boundary.Norm))
			}
		case boundary.Start && boundary.Count == 0:
			coverColor = ansi.ColorCode("red")
		case !boundary.Start:
			coverColor = ansi.ColorCode("reset")
		}
		result = append(result, []byte(coverColor)...)

		curOffset = boundary.Offset
	}
	if curOffset < len(fileBytes) {
		result = append(result, fileBytes[curOffset:]...)
	}

	result = append(result, []byte("\n")...)
	return result
}

func getTempFileName() (string, error) {
	tmpFile, err := ioutil.TempFile(".", "go-carpet-coverage-out-")
	if err != nil {
		return "", err
	}
	tmpFile.Close()

	return tmpFile.Name(), nil
}

var (
	config struct {
		filesFilter   string
		funcFilter    string
		colors256     bool
		includeVendor bool
	}
)

func init() {
	flag.StringVar(&config.filesFilter, "file", "", "comma-separated list of files to test (default: all)")
	flag.StringVar(&config.funcFilter, "func", "", "comma-separated functions list (default: all functions)")
	flag.BoolVar(&config.colors256, "256colors", false, "use more colors on 256-color terminal (indicate the level of coverage)")
	flag.BoolVar(&config.includeVendor, "include-vendor", false, "include vendor directories for show coverage (Godeps, vendor)")
	flag.Usage = func() {
		fmt.Println(usageMessage)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	flag.Parse()
	testDirs := flag.Args()

	coverFileName, err := getTempFileName()
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(coverFileName)
	stdOut := getColorWriter()
	allProfileBlocks := []cover.ProfileBlock{}

	if len(testDirs) > 0 {
		testDirs = getDirsWithTests(config.includeVendor, testDirs...)
	} else {
		testDirs = getDirsWithTests(config.includeVendor, ".")
	}
	for _, path := range testDirs {
		err := runGoTest(path, coverFileName, false)
		if err != nil {
			log.Print(err)
			continue
		}

		coverInBytes, profileBlocks, err := getCoverForDir(path, coverFileName, strings.Split(config.filesFilter, ","), config.colors256)
		if err != nil {
			log.Print(err)
			continue
		}
		stdOut.Write(coverInBytes)
		allProfileBlocks = append(allProfileBlocks, profileBlocks...)
	}

	if len(allProfileBlocks) > 0 {
		stat := getStatForProfileBlocks(allProfileBlocks)
		totalCoverage := fmt.Sprintf("Coverage: %.1f%% of statements", stat)
		stdOut.Write([]byte(getColorHeader(totalCoverage, false)))
	}
}
