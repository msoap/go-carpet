package main

import (
	"flag"
	"fmt"
	"go/build"
	"io"
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

const (
	usageMessage = `go-carpet - show test coverage for Go source files

usage: go-carpet [options] [paths]`

	version = "1.9.0"

	// predefined go test options
	goTestCoverProfile = "-coverprofile"
	goTestCoverMode    = "-covermode"
)

var (
	reNewLine        = regexp.MustCompile("\n")
	reWindowsPathFix = regexp.MustCompile(`^_\\([A-Z])_`)

	// vendors directories for skip
	vendorDirs = []string{"Godeps", "vendor", ".vendor", "_vendor"}

	// directories for skip
	skipDirs = []string{"testdata"}

	errIsNotInGoMod = fmt.Errorf("is not in go modules")
)

func getDirsWithTests(includeVendor bool, roots ...string) (result []string, err error) {
	if len(roots) == 0 {
		roots = []string{"."}
	}

	dirs := map[string]struct{}{}
	for _, root := range roots {
		err = filepath.Walk(root, func(path string, _ os.FileInfo, _ error) error {
			if strings.HasSuffix(path, "_test.go") {
				dirs[filepath.Dir(path)] = struct{}{}
			}
			return nil
		})
		if err != nil {
			return result, err
		}
	}

	result = make([]string, 0, len(dirs))
	for dir := range dirs {
		if !includeVendor && isSliceInStringPrefix(dir, vendorDirs) || isSliceInStringPrefix(dir, skipDirs) {
			continue
		}
		result = append(result, "./"+dir)
	}

	return result, nil
}

func readFile(fileName string) (result []byte, err error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return result, err
	}

	result, err = io.ReadAll(fileReader)
	if err == nil {
		err = fileReader.Close()
	}

	return result, err
}

func getShadeOfGreen(normCover float64) string {
	/*
		Get all colors for 255-colors terminal:
			gommand 'for i := 0; i < 256; i++ {fmt.Println(i, ansi.ColorCode(strconv.Itoa(i)) + "String" + ansi.ColorCode("reset"))}'
	*/
	var tenShadesOfGreen = [...]string{
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

func runGoTest(path string, coverFileName string, goTestArgs []string, hideStderr bool) error {
	args := []string{"test", goTestCoverProfile + "=" + coverFileName, goTestCoverMode + "=count"}
	args = append(args, goTestArgs...)
	args = append(args, path)
	osExec := exec.Command("go", args...) // #nosec
	if !hideStderr {
		osExec.Stderr = os.Stderr
	}

	if output, err := osExec.Output(); err != nil {
		fmt.Print(string(output))
		return err
	}

	return nil
}

func guessAbsPathInGOPATH(GOPATH, relPath string) (absPath string, err error) {
	if GOPATH == "" {
		GOPATH = build.Default.GOPATH
		if GOPATH == "" {
			return "", fmt.Errorf("GOPATH is not set")
		}
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
		return "", fmt.Errorf("file '%s' not found in GOPATH", relPath)
	}

	return absPath, err
}

func getCoverForDir(coverFileName string, filesFilter []string, config Config) (result []byte, profileBlocks []cover.ProfileBlock, err error) {
	coverProfile, err := cover.ParseProfiles(coverFileName)
	if err != nil {
		return result, profileBlocks, err
	}

	for _, fileProfile := range coverProfile {
		// Skip files if minimal coverage is set and is covered more than minimal coverage
		if config.minCoverage > 0 && config.minCoverage < 100.0 && getStatForProfileBlocks(fileProfile.Blocks) > config.minCoverage {
			continue
		}

		var fileName string
		if strings.HasPrefix(fileProfile.FileName, "/") {
			// TODO: what about windows?
			fileName = fileProfile.FileName
		} else if strings.HasPrefix(fileProfile.FileName, "_") {
			// absolute path (or relative in tests)
			if runtime.GOOS != "windows" {
				fileName = strings.TrimLeft(fileProfile.FileName, "_")
			} else {
				// "_\C_\Users\..." -> "C:\Users\..."
				fileName = reWindowsPathFix.ReplaceAllString(fileProfile.FileName, "$1:")
			}
		} else if fileName, err = guessAbsPathInGoMod(fileProfile.FileName); err != errIsNotInGoMod {
			if err != nil {
				return result, profileBlocks, err
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

		var fileBytes []byte
		fileBytes, err = readFile(fileName)
		if err != nil {
			return result, profileBlocks, err
		}

		result = append(result, getCoverForFile(fileProfile, fileBytes, config)...)
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
//
//	src/cmd/cover/html.go::percentCovered()
//	src/testing/cover.go::coverReport()
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

func getCoverForFile(fileProfile *cover.Profile, fileBytes []byte, config Config) (result []byte) {
	stat := getStatForProfileBlocks(fileProfile.Blocks)

	// Avoid if flag.Parse() was not called yet.
	if config.minCoverage == 0 {
		config.minCoverage = stat
	}

	// Retun empty to skip if minimal coverage is set and is covered more than minimal coverage
	if int(stat) > int(config.minCoverage) {
		return []byte{}
	}

	textRanges, err := getFileFuncRanges(fileBytes, config.funcFilter)
	if err != nil {
		return result
	}

	var fileNameDisplay string
	if len(config.funcFilter) == 0 {
		fileNameDisplay = fmt.Sprintf("%s - %.1f%%", strings.TrimLeft(fileProfile.FileName, "_"), stat)
	} else {
		fileNameDisplay = strings.TrimLeft(fileProfile.FileName, "_")
	}

	if config.summary {
		return []byte(fileNameDisplay + "\n")
	}

	result = append(result, []byte(getColorHeader(fileNameDisplay, true))...)

	boundaries := fileProfile.Boundaries(fileBytes)

	for _, textRange := range textRanges {
		fileBytesPart := fileBytes[textRange.begin:textRange.end]
		curOffset := 0
		coverColor := ""

		for _, boundary := range boundaries {
			if boundary.Offset < textRange.begin || boundary.Offset > textRange.end {
				// skip boundary which is not in filter function
				continue
			}

			boundaryOffset := boundary.Offset - textRange.begin

			if boundaryOffset > curOffset {
				nextChunk := fileBytesPart[curOffset:boundaryOffset]
				// Add ansi color code in begin of each line (this fixed view in "less -R")
				if coverColor != "" && coverColor != ansi.ColorCode("reset") {
					nextChunk = reNewLine.ReplaceAllLiteral(nextChunk, []byte(ansi.ColorCode("reset")+"\n"+coverColor))
				}
				result = append(result, nextChunk...)
			}

			switch {
			case boundary.Start && boundary.Count > 0:
				coverColor = ansi.ColorCode("green")
				if config.colors256 {
					coverColor = ansi.ColorCode(getShadeOfGreen(boundary.Norm))
				}
			case boundary.Start && boundary.Count == 0:
				coverColor = ansi.ColorCode("red")
			case !boundary.Start:
				coverColor = ansi.ColorCode("reset")
			}
			result = append(result, []byte(coverColor)...)

			curOffset = boundaryOffset
		}
		if curOffset < len(fileBytesPart) {
			result = append(result, fileBytesPart[curOffset:]...)
		}

		result = append(result, []byte("\n")...)
	}

	return result
}

type textRange struct {
	begin, end int
}

func getFileFuncRanges(fileBytes []byte, funcs []string) (result []textRange, err error) {
	if len(funcs) == 0 {
		return []textRange{{
			begin: 0,
			end:   len(fileBytes),
		}}, nil
	}

	golangFuncs, err := getGolangFuncs(fileBytes)
	if err != nil {
		return nil, err
	}

	for _, existsFunc := range golangFuncs {
		for _, filterFuncName := range funcs {
			if existsFunc.Name == filterFuncName {
				result = append(result, textRange{begin: existsFunc.Begin - 1, end: existsFunc.End - 1})
			}
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("filter by functions: %v - not found", funcs)
	}

	return result, nil
}

func getTempFileName() (string, error) {
	tmpFile, err := os.CreateTemp(".", "go-carpet-coverage-out-")
	if err != nil {
		return "", err
	}
	err = tmpFile.Close()
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

// Config - application config
type Config struct {
	filesFilterRaw string
	filesFilter    []string
	funcFilterRaw  string
	funcFilter     []string
	argsRaw        string
	minCoverage    float64
	colors256      bool
	includeVendor  bool
	summary        bool
}

var config Config

func init() {
	flag.StringVar(&config.filesFilterRaw, "file", "", "comma-separated list of `files` to test (default: all)")
	flag.StringVar(&config.funcFilterRaw, "func", "", "comma-separated `functions` list (default: all functions)")
	flag.BoolVar(&config.colors256, "256colors", false, "use more colors on 256-color terminal (indicate the level of coverage)")
	flag.BoolVar(&config.summary, "summary", false, "only show summary for each file")
	flag.BoolVar(&config.includeVendor, "include-vendor", false, "include vendor directories for show coverage (Godeps, vendor)")
	flag.StringVar(&config.argsRaw, "args", "", "pass additional `arguments` for go test")
	flag.Float64Var(&config.minCoverage, "mincov", 100.0, "coverage threshold of the file to be displayed (in percent)")
	flag.Usage = func() {
		fmt.Println(usageMessage)
		flag.PrintDefaults()
		os.Exit(0)
	}
}

func main() {
	versionFl := flag.Bool("version", false, "get version")
	flag.Parse()

	if *versionFl {
		fmt.Println(version)
		os.Exit(0)
	}

	config.filesFilter = grepEmptyStringSlice(strings.Split(config.filesFilterRaw, ","))
	config.funcFilter = grepEmptyStringSlice(strings.Split(config.funcFilterRaw, ","))
	additionalArgs, err := parseAdditionalArgs(config.argsRaw, []string{goTestCoverProfile, goTestCoverMode})
	if err != nil {
		log.Fatal(err)
	}

	testDirs := flag.Args()

	coverFileName, err := getTempFileName()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = os.RemoveAll(coverFileName)
		if err != nil {
			log.Fatal(err)
		}
	}()

	stdOut := getColorWriter()
	allProfileBlocks := []cover.ProfileBlock{}

	if len(testDirs) > 0 {
		testDirs, err = getDirsWithTests(config.includeVendor, testDirs...)
	} else {
		testDirs, err = getDirsWithTests(config.includeVendor, ".")
	}
	if err != nil {
		log.Fatal(err)
	}

	for _, path := range testDirs {
		if err = runGoTest(path, coverFileName, additionalArgs, false); err != nil {
			log.Print(err)
			continue
		}

		coverInBytes, profileBlocks, errCover := getCoverForDir(coverFileName, config.filesFilter, config)
		if errCover != nil {
			log.Print(errCover)
			continue
		}
		_, err = stdOut.Write(coverInBytes)
		if err != nil {
			log.Fatal(err)
		}

		allProfileBlocks = append(allProfileBlocks, profileBlocks...)
	}

	if len(allProfileBlocks) > 0 && len(config.funcFilter) == 0 {
		stat := getStatForProfileBlocks(allProfileBlocks)
		totalCoverage := fmt.Sprintf("Coverage: %.1f%% of statements", stat)
		_, err = stdOut.Write([]byte(getColorHeader(totalCoverage, false)))
		if err != nil {
			log.Fatal(err)
		}
	}
}
