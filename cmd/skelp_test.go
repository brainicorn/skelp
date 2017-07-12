package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brainicorn/skelp/skelputil"
)

func TestSkelpCmdError(t *testing.T) {

	code := Execute([]string{"badcommand"}, nil)

	if code == 0 {
		t.Errorf("execute should have errored ")
	}

}

func TestSkelpCmdUserError(t *testing.T) {

	code := Execute([]string{"apply"}, nil)

	if code == 0 {
		t.Errorf("execute should have errored ")
	}

}

func TestSkelpCmdBadHomedir(t *testing.T) {
	out := new(bytes.Buffer)
	code := Execute([]string{"alias", "list", "--no-color", "--homedir", "/does.not.exist"}, out)

	if code == 0 {
		t.Errorf("execute should have errored ")
	}

	lines := strings.Split(out.String(), "\n")
	if !strings.HasSuffix(lines[0], "is not a valid path for --homedir flag") {
		fmt.Println(lines[0])
		t.Errorf("alias error does not match")
	}

}

func TestSkelpCmdSkelpdirOVerride(t *testing.T) {
	out := new(bytes.Buffer)

	tmpDir, _ := ioutil.TempDir("", "custom-home")
	defer os.RemoveAll(tmpDir)

	code := Execute([]string{"alias", "list", "--no-color", "--homedir", tmpDir, "--skelpdir", "customdir"}, out)

	if code != 0 {
		t.Errorf("execute should not have errored ")
	}

	if !skelputil.PathExists(filepath.Join(tmpDir, "customdir")) {
		t.Errorf("custom skepdir should exist: %s", filepath.Join(tmpDir, "customdir"))
	}

}
