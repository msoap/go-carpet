package main

import (
	"io/ioutil"
	"log"
	"os"
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

func Test_parseAdditionalArgs(t *testing.T) {
	tests := []struct {
		name           string
		argsRaw        string
		excludes       []string
		wantResultArgs []string
		wantErr        bool
	}{
		{
			name:           "empty args",
			argsRaw:        "",
			excludes:       []string{},
			wantResultArgs: nil,
			wantErr:        false,
		},
		{
			name:           "with error",
			argsRaw:        "str '...",
			excludes:       []string{},
			wantResultArgs: nil,
			wantErr:        true,
		},
		{
			name:           "one args",
			argsRaw:        `-short`,
			excludes:       []string{},
			wantResultArgs: []string{"-short"},
			wantErr:        false,
		},
		{
			name:           "all args",
			argsRaw:        `-short    -option`,
			excludes:       []string{},
			wantResultArgs: []string{"-short", "-option"},
			wantErr:        false,
		},
		{
			name:           "with excludes",
			argsRaw:        `-short -ex=23 -new '-ex2 "ex word"' -option`,
			excludes:       []string{"-ex", "-ex2"},
			wantResultArgs: []string{"-short", "-new", "-option"},
			wantErr:        false,
		},
	}

	log.SetOutput(ioutil.Discard)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResultArgs, err := parseAdditionalArgs(tt.argsRaw, tt.excludes)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseAdditionalArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResultArgs, tt.wantResultArgs) {
				t.Errorf("parseAdditionalArgs() = %v, want %v", gotResultArgs, tt.wantResultArgs)
			}
		})
	}

	log.SetOutput(os.Stdout)
}
