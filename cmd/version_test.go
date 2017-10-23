package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	out := new(bytes.Buffer)
	tmpHomeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpHomeDir)

	tmpOutputDir, _ := ioutil.TempDir("", "skelp-output")
	defer os.RemoveAll(tmpOutputDir)

	code := Execute([]string{"version", "--homedir", tmpHomeDir}, out)

	if code != 0 {
		fmt.Println(out)
		t.Errorf("version should not have errored")
	}
}
