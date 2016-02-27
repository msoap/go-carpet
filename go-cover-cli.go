package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/mgutz/ansi"
	"golang.org/x/tools/cover"
)

func readFileByLines(fileName string) (result []string, err error) {
	fileReader, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer fileReader.Close()

	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}

type Marker struct {
	col     int
	isBegin bool
	count   int
}

type MarkerList []Marker

func (list MarkerList) Len() int           { return len(list) }
func (list MarkerList) Swap(a, b int)      { list[a], list[b] = list[b], list[a] }
func (list MarkerList) Less(a, b int) bool { return list[a].col < list[b].col }

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

		profileBlockByLine := map[int]MarkerList{}
		for _, profileBlock := range fileProfile.Blocks {
			profileBlockByLine[profileBlock.StartLine] = append(profileBlockByLine[profileBlock.StartLine], Marker{
				col:     profileBlock.StartCol,
				isBegin: true,
				count:   profileBlock.Count,
			})
			profileBlockByLine[profileBlock.EndLine] = append(profileBlockByLine[profileBlock.EndLine], Marker{
				col:     profileBlock.EndCol,
				isBegin: false,
			})
		}

		fileContent, err := readFileByLines(fileName)
		if err != nil {
			log.Fatal(err)
		}

		for num, line := range fileContent {
			if lineBlocks, ok := profileBlockByLine[num+1]; ok {
				sort.Sort(lineBlocks)

				lineRunes := []rune(line)
				colorChunks := []string{}
				curCol := 0

				for _, block := range lineBlocks {
					colorChunks = append(colorChunks, string(lineRunes[curCol:block.col-1]))
					if block.isBegin {
						colorCode := ""
						if block.count > 0 {
							colorCode = ansi.ColorCode("green")
						} else {
							colorCode = ansi.ColorCode("red")
						}
						colorChunks = append(colorChunks, colorCode)
					} else {
						colorChunks = append(colorChunks, ansi.ColorCode("reset"))
					}

					curCol = block.col - 1
				}
				if curCol < len(lineRunes) {
					colorChunks = append(colorChunks, string(lineRunes[curCol:len(lineRunes)]))
				}

				fmt.Println(strings.Join(colorChunks, ""))
			} else {
				fmt.Println(line)
			}
		}
	}

	err = os.Remove(coverFile)
	if err != nil {
		log.Fatal(err)
	}
}
