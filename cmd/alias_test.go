package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

const (
	aliasListHeader = `------------------
Registered Aliases
------------------
`
)

var (
	listTestAlias = aliasListHeader + fmt.Sprintf("%s -> %s\n", "test.alias", "/tmp")
)

func TestAliasAdd(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "add", "test.alias", "/tmp", "--no-color", "--homedir", tmpDir}, nil)

	if code != 0 {
		t.Errorf("alias add should not have errored")
	}
}

func TestAliasAddMissingArg(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "add", "test.alias", "--no-color", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("alias add should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errAliasAddMissingArgs {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errAliasAddMissingArgs)
		t.Errorf("alias error does not match")
	}
}

func TestAliasAddBadAlias(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "add", "./", "my.alias", "--no-color", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("alias add should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errAliasAddBadAlias {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errAliasAddBadAlias)
		t.Errorf("alias error does not match")
	}
}

func TestAliasAddBadPath(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "add", "my.alias", "my.alias", "--no-color", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("alias add should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errAliasAddBadPath {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errAliasAddBadPath)
		t.Errorf("alias error does not match")
	}
}

func TestAliasRemove(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	Execute([]string{"alias", "add", "test.alias", "/tmp", "--no-color", "--homedir", tmpDir}, nil)

	code := Execute([]string{"alias", "remove", "test.alias", "--no-color", "--homedir", tmpDir}, out)

	if code != 0 {
		fmt.Println(out)
		t.Errorf("alias remove should not have errored")
	}
}

func TestAliasRemoveMissingName(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "remove", "--no-color", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("alias remove should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errAliasRemoveMissingName {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errAliasRemoveMissingName)
		t.Errorf("alias error does not match")
	}
}

func TestAliasRemoveBadName(t *testing.T) {
	out := new(bytes.Buffer)
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "remove", "./", "--no-color", "--homedir", tmpDir}, out)

	if code == 0 {
		t.Errorf("alias remove should have errored")
	}

	lines := strings.Split(out.String(), "\n")
	if lines[0] != errAliasRemoveInvalidName {
		fmt.Println(lines[0])
		fmt.Println("........................")
		fmt.Println(errAliasRemoveInvalidName)
		t.Errorf("alias error does not match")
	}
}

func TestAliasListEmpty(t *testing.T) {
	out := new(bytes.Buffer)

	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "list", "--no-color", "--homedir", tmpDir}, out)

	if code != 0 {
		fmt.Println(out)
		t.Errorf("alias list should not have errored")
	}

	if out.String() != aliasListHeader {
		fmt.Println(out)
		t.Errorf("alias list should return an empty list")
	}
}

func TestAliasList(t *testing.T) {
	out := new(bytes.Buffer)

	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	// create an alias for testing
	Execute([]string{"alias", "add", "test.alias", "/tmp", "--no-color", "--homedir", tmpDir}, nil)
	code := Execute([]string{"alias", "list", "--homedir", tmpDir}, out)

	if code != 0 {
		fmt.Println(out)
		t.Errorf("alias list should not have errored")
	}

	if out.String() != listTestAlias {
		fmt.Println(out)
		t.Errorf("alias list should have returned the test.alias")
	}
}
