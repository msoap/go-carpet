package main

import (
	"fmt"
	"log"
	"strings"

	shellwords "github.com/mattn/go-shellwords"
)

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

// grepEmptyStringSlice - return slice with non-empty strings
func grepEmptyStringSlice(inSlice []string) []string {
	result := []string{}
	for _, item := range inSlice {
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	return result
}

// parse additional args for pass to go test
func parseAdditionalArgs(argsRaw string, excludes []string) (resultArgs []string, err error) {
	if argsRaw != "" {
		args, err := shellwords.Parse(argsRaw)
		if err != nil {
			return resultArgs, fmt.Errorf("args %q parse failed: %s", argsRaw, err)
		}

	NEXTARG:
		for _, arg := range args {
			for _, excludeArg := range excludes {
				if excludeArg != "" && strings.HasPrefix(arg, excludeArg) {
					log.Printf("arg: %q is not allowed, skip", arg)
					continue NEXTARG
				}
			}
			resultArgs = append(resultArgs, arg)
		}
	}

	return resultArgs, nil
}
