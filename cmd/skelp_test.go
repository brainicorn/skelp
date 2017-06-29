package cmd

import (
	"os"
	"testing"
)

func TestSkelpCmdError(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{"skelp", "badcommand"}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("execute should have errored ")
		}
	}()

	Execute()

}

func TestSkelpCmdUserError(t *testing.T) {
	oldArgs := os.Args
	defer func() {
		os.Args = oldArgs
	}()

	os.Args = []string{"skelp", "apply"}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("execute should have errored ")
		}
	}()

	Execute()

}
