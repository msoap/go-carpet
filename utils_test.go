package main

import (
	"reflect"
	"testing"
)

func Test_isSliceInString(t *testing.T) {
	testData := []struct {
		src    string
		slice  []string
		result bool
	}{
		{
			src:    "one/file.go",
			slice:  []string{"one.go", "file.go"},
			result: true,
		},
		{
			src:    "path/path/file.go",
			slice:  []string{"one.go", "path/file.go"},
			result: true,
		},
		{
			src:    "one/file.go",
			slice:  []string{"one.go", "two.go"},
			result: false,
		},
		{
			src:    "one/file.go",
			slice:  []string{},
			result: false,
		},
	}

	for i, item := range testData {
		result := isSliceInString(item.src, item.slice)
		if result != item.result {
			t.Errorf("\n%d. isSliceInString()\nexpected: %v\nreal    :%v", i, item.result, result)
		}
	}
}

func Test_isSliceInStringPrefix(t *testing.T) {
	testData := []struct {
		src    string
		slice  []string
		result bool
	}{
		{
			src:    "one/file.go",
			slice:  []string{"vendor", "Godeps"},
			result: false,
		},
		{
			src:    "vendor/path/file.go",
			slice:  []string{"vendor", "Godeps"},
			result: true,
		},
		{
			src:    "Godeps/path/file.go",
			slice:  []string{"vendor", "Godeps"},
			result: true,
		},
		{
			src:    "one/file.go",
			slice:  []string{},
			result: false,
		},
	}

	for i, item := range testData {
		result := isSliceInStringPrefix(item.src, item.slice)
		if result != item.result {
			t.Errorf("\n%d. isSliceInStringPrefix()\nexpected: %v\nreal     :%v", i, item.result, result)
		}
	}
}

func Test_grepEmptyStringSlice(t *testing.T) {
	testData := []struct {
		inSlice []string
		result  []string
	}{
		{
			inSlice: []string{},
			result:  []string{},
		},
		{
			inSlice: nil,
			result:  []string{},
		},
		{
			inSlice: []string{""},
			result:  []string{},
		},
		{
			inSlice: []string{"A", "B"},
			result:  []string{"A", "B"},
		},
		{
			inSlice: []string{"A", "", "B"},
			result:  []string{"A", "B"},
		},
		{
			inSlice: []string{"", "", "B"},
			result:  []string{"B"},
		},
	}

	for i, item := range testData {
		result := grepEmptyStringSlice(item.inSlice)

		if !reflect.DeepEqual(result, item.result) {
			t.Errorf("\n%d. grepEmptyStringSlice()\nexpected: %#v\nreal     :%#v", i, item.result, result)
		}
	}
}
