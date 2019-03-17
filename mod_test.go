package main

import (
	"os"
	"testing"
)

func Test_guessAbsPathInGoMod(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		if _, err := guessAbsPathInGoMod(""); err == nil {
			t.Errorf("failed to test empty file")
		}
	})

	t.Run("real", func(t *testing.T) {
		gotAbsPath, err := guessAbsPathInGoMod("github.com/msoap/go-carpet/terminal_posix.go")
		if err != nil {
			t.Errorf("failed to test real file")
		}

		if _, err := os.Stat(gotAbsPath); err != nil {
			t.Errorf("failed to test real file")
		}
	})

	t.Run("not exists", func(t *testing.T) {
		_, err := guessAbsPathInGoMod("github.com/msoap/go-carpet/terminal_posix_another_file.go")
		if err == nil {
			t.Errorf("failed to test not exists file")
		}
	})
}
