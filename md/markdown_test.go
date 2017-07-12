package main

import (
	"os"
	"testing"
)

func TestMD(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{"markdown", "./"}
	main()
}
