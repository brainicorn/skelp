package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMainOK(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("main should not have errored %s", r)
		}
	}()

	if filepath.Base(os.Args[0]) == "skelp.test" || filepath.Base(os.Args[0]) == "main.test" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()

		os.Args = []string{"skelp"}
	}

	main()
}

func TestMainErr(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("main should have errored %s", r)
		}
	}()

	if filepath.Base(os.Args[0]) == "skelp.test" || filepath.Base(os.Args[0]) == "main.test" {
		oldArgs := os.Args
		defer func() {
			os.Args = oldArgs
		}()

		os.Args = []string{"skelp", "apply"}
	}

	main()
}
