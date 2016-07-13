package main

import "strings"

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
func grepEmptyStringSlice(inSlice []string) (result []string) {
	for _, item := range inSlice {
		if len(item) > 0 {
			result = append(result, item)
		}
	}
	return result
}
