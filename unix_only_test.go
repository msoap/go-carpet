//go:build !windows
// +build !windows

package main

import (
	"reflect"
	"testing"
)

func Test_getCoverForDir(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		_, _, err := getCoverForDir("./testdata/not_exists.out", []string{}, Config{colors256: false})
		if err == nil {
			t.Errorf("1. getCoverForDir() error failed")
		}
	})

	t.Run("cover", func(t *testing.T) {
		bytes, _, err := getCoverForDir("./testdata/cover_00.out", []string{}, Config{colors256: false})
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
	})

	t.Run("cover with 256 colors", func(t *testing.T) {
		bytes, _, err := getCoverForDir("./testdata/cover_00.out", []string{}, Config{colors256: true})
		if err != nil {
			t.Errorf("5. getCoverForDir() failed: %v", err)
		}
		expect, err := readFile("./testdata/colored_01.txt")
		if err != nil {
			t.Errorf("6. getCoverForDir() failed: %v", err)
		}
		if !reflect.DeepEqual(bytes, expect) {
			t.Errorf("7. getCoverForDir() not equal")
		}
	})

	t.Run("cover with 256 colors with error", func(t *testing.T) {
		_, _, err := getCoverForDir("./testdata/cover_01.out", []string{}, Config{colors256: true})
		if err == nil {
			t.Errorf("8. getCoverForDir() not exists go file")
		}
	})

	t.Run("cover 01 without 256 colors", func(t *testing.T) {
		bytes, _, err := getCoverForDir("./testdata/cover_00.out", []string{"file_01.go"}, Config{colors256: false})
		if err != nil {
			t.Errorf("9. getCoverForDir() failed: %v", err)
		}
		expect, err := readFile("./testdata/colored_02.txt")
		if err != nil {
			t.Errorf("10. getCoverForDir() failed: %v", err)
		}
		if !reflect.DeepEqual(bytes, expect) {
			t.Errorf("11. getCoverForDir() not equal")
		}
	})

	t.Run("cover 02 without 256 colors", func(t *testing.T) {
		bytes, _, err := getCoverForDir("./testdata/cover_02.out", []string{}, Config{colors256: false})
		if err != nil {
			t.Errorf("12. getCoverForDir() failed: %v", err)
		}
		expect, err := readFile("./testdata/colored_03.txt")
		if err != nil {
			t.Errorf("13. getCoverForDir() failed: %v", err)
		}
		if !reflect.DeepEqual(bytes, expect) {
			t.Errorf("14. getCoverForDir() not equal\ngot:\n%s\nexpect:\n%s\n", bytes, expect)
		}
	})
}
