package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/brainicorn/skelp/skelputil"
)

func TestDoDownload(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New()
	err := gen.doDownload("git@github.com:brainicorn/skelp.git", tmpDir, DefaultOptions())

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

}

func TestDoCheckForUpdates(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New()
	err := gen.doDownload("git@github.com:brainicorn/skelp.git", tmpDir, DefaultOptions())

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

	err = gen.checkForUpdates("git@github.com:brainicorn/skelp.git", tmpDir, DefaultOptions())

	if err != nil {
		t.Fatalf("error updating %s", err.Error())
	}
}