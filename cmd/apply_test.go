package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestApply(t *testing.T) {
	out := new(bytes.Buffer)
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpOutputDir, _ := ioutil.TempDir("", "skelp-output")
	defer os.RemoveAll(tmpOutputDir)

	code := Execute([]string{"apply", "../testdata/generator/simple", "--no-color", "--force", "--offline", "--homedir", tmpHomeDir, "-o", tmpOutputDir, "-d", "../testdata/generator/simple-data.json"}, out)

	if code != 0 {
		fmt.Println(out)
		t.Errorf("apply should not have errored")
	}
}

func TestApplyDryRun(t *testing.T) {
	out := new(bytes.Buffer)
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpOutputDir, _ := ioutil.TempDir("", "skelp-output")
	defer os.RemoveAll(tmpOutputDir)

	code := Execute([]string{"apply", "../testdata/generator/simple", "--no-color", "--force", "--offline", "--dry-run", "-d", "../testdata/generator/simple-data.json"}, out)

	if code != 0 {
		fmt.Println(out)
		t.Errorf("apply should not have errored")
	}
}

func TestApplyBadDescriptor(t *testing.T) {
	out := new(bytes.Buffer)
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpOutputDir, _ := ioutil.TempDir("", "skelp-output")
	defer os.RemoveAll(tmpOutputDir)

	code := Execute([]string{"apply", "../testdata/generator/baddescriptor", "--no-color", "--force", "--offline", "--homedir", tmpHomeDir, "-o", tmpOutputDir, "-d", "../testdata/generator/simple-data.json"}, out)

	if code != 1 {
		t.Errorf("apply should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	fmt.Println("err lines: ", lines)
	if lines[0] != "Error validating skelp descriptor:" {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println("Error validating skelp descriptor:")
		t.Errorf("apply error does not match")
	}
}

func TestApplyMissingArg(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"apply", "--no-color", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("apply should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errApplyMissingArgs {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errApplyMissingArgs)
		t.Errorf("apply error does not match")
	}
}

func TestApplyBadOutputDir(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"apply", "test.alias", "--no-color", "--homedir", tmpDir, "-o", "./does.not.exist"}, out)

	if code == 0 {
		t.Errorf("apply should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if !strings.HasSuffix(lines[0], "is not a valid path for --output flag") {
		fmt.Println(lines[0])
		t.Errorf("apply error does not match")
	}
}

func TestApplyBadDataFile(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"apply", "test.alias", "--no-color", "--homedir", tmpDir, "-d", "./does.not.exist.json"}, out)

	if code == 0 {
		t.Errorf("apply should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if !strings.HasSuffix(lines[0], "is not a valid path for --data flag") {
		fmt.Println(lines[0])
		t.Errorf("apply error does not match")
	}
}
