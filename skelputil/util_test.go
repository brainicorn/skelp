package skelputil

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestDirIsEmptyErr(t *testing.T) {
	b := DirIsEmpty("doesn't exist")

	if b {
		t.Error("non-existent dir should return false for empty")
	}
}

func TestDirIsEmptyTrue(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-empty-dir")
	defer os.RemoveAll(tmpDir)

	b := DirIsEmpty(tmpDir)

	if !b {
		t.Error("empty dir should return true for empty")
	}
}

func TestMkdirAllErr(t *testing.T) {
	err := MkdirAll("")

	if err == nil {
		t.Error("bad path should have thrown an error")
	}
}

func TestMkdirAllExists(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-mkdir")
	defer os.RemoveAll(tmpDir)

	err := MkdirAll(tmpDir)

	if err != nil {
		t.Error("existing path should not throw an error")
	}
}
