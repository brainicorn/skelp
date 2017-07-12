package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/brainicorn/skelp/skelputil"
)

func TestBashMe(t *testing.T) {
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpCompletionDir, _ := ioutil.TempDir("", "custom-completion")
	defer os.RemoveAll(tmpCompletionDir)

	code := Execute([]string{"bashme", "--no-color", "--homedir", tmpHomeDir, "--output", tmpCompletionDir, "--no-sudo"}, nil)

	if code != 0 {
		t.Errorf("bashme should not have errored")
	}

	if !skelputil.PathExists(filepath.Join(tmpCompletionDir, "completion.sh")) {
		t.Errorf("completion file should exist at: %s", filepath.Join(tmpCompletionDir, "completion.sh"))
	}
}
