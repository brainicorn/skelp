package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestValidateSuccess(t *testing.T) {
	out := new(bytes.Buffer)
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpOutputDir, _ := ioutil.TempDir("", "skelp-output")
	defer os.RemoveAll(tmpOutputDir)

	code := Execute([]string{"validate", "../testdata/generator/simple/skelp.json", "--no-color", "--homedir", tmpHomeDir}, out)

	if code != 0 {
		t.Errorf("validate should not have errored")
	}

	if out.String() != "Descriptor OK.\n" {
		t.Errorf("should have gotten 'Descriptor OK.' but got %s", out.String())
	}
}

func TestValidateError(t *testing.T) {
	out := new(bytes.Buffer)
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpOutputDir, _ := ioutil.TempDir("", "skelp-output")
	defer os.RemoveAll(tmpOutputDir)

	code := Execute([]string{"validate", "../testdata/generator/baddescriptor/skelp.json", "--no-color", "--homedir", tmpHomeDir}, out)

	if code == 0 {
		t.Errorf("validate should have errored")
	}

	if !strings.HasPrefix(out.String(), "Error validating skelp descriptor:") {
		t.Errorf("should have gotten 'Error validating skelp descriptor:' but got %s", out.String())
	}
}

func TestValidateMissingArg(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"validate", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("validate should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errValidateMissingArgs {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errValidateMissingArgs)
		t.Errorf("apply error does not match")
	}
}

func TestValidateBadPath(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"validate", "../does.not.exist/skelp.json", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("validate should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if !strings.HasPrefix(lines[0], "descriptor not found:") {
		fmt.Println(lines[0])
		t.Errorf("validate error does not match")
	}
}
