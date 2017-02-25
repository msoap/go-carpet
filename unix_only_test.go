// +build !windows

package main

import (
	"reflect"
	"testing"
)

func Test_getCoverForDir(t *testing.T) {
	_, _, err := getCoverForDir("./testdata", "./testdata/not_exists.out", []string{}, Config{colors256: false})
	if err == nil {
		t.Errorf("1. getCoverForDir() error failed")
	}

	// ---
	bytes, _, err := getCoverForDir("./testdata", "./testdata/cover_00.out", []string{}, Config{colors256: false})
	if err != nil {
		t.Errorf("2. getCoverForDir() failed: %v", err)
	}
	expect, err := readFile("./testdata/colored_00.txt")
	if err != nil {
		t.Errorf("3. getCoverForDir() failed: %v", err)
	}
	if !reflect.DeepEqual(bytes, expect) {
		t.Errorf("4. getCoverForDir() not equal")
	}

	// ---
	bytes, _, err = getCoverForDir("./testdata", "./testdata/cover_00.out", []string{}, Config{colors256: true})
	if err != nil {
		t.Errorf("5. getCoverForDir() failed: %v", err)
	}
	expect, err = readFile("./testdata/colored_01.txt")
	if err != nil {
		t.Errorf("6. getCoverForDir() failed: %v", err)
	}
	if !reflect.DeepEqual(bytes, expect) {
		t.Errorf("7. getCoverForDir() not equal")
	}

	// ---
	_, _, err = getCoverForDir("./testdata", "./testdata/cover_01.out", []string{}, Config{colors256: true})
	if err == nil {
		t.Errorf("8. getCoverForDir() not exists go file")
	}

	// ---
	bytes, _, err = getCoverForDir("./testdata", "./testdata/cover_00.out", []string{"file_01.go"}, Config{colors256: false})
	if err != nil {
		t.Errorf("9. getCoverForDir() failed: %v", err)
	}
	expect, err = readFile("./testdata/colored_02.txt")
	if err != nil {
		t.Errorf("10. getCoverForDir() failed: %v", err)
	}
	if !reflect.DeepEqual(bytes, expect) {
		t.Errorf("11. getCoverForDir() not equal")
	}

	// ---
	bytes, _, err = getCoverForDir("./testdata", "./testdata/cover_02.out", []string{}, Config{colors256: false})
	if err != nil {
		t.Errorf("12. getCoverForDir() failed: %v", err)
	}
	expect, err = readFile("./testdata/colored_03.txt")
	if err != nil {
		t.Errorf("13. getCoverForDir() failed: %v", err)
	}
	if !reflect.DeepEqual(bytes, expect) {
		t.Errorf("14. getCoverForDir() not equal")
	}
}
