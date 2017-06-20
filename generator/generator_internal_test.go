package generator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/brainicorn/skelp/skelputil"
)

func TestDoDownloadSSH(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New(DefaultOptions())
	err := gen.doDownload("git@github.com:brainicorn/skelp.git", tmpDir)

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

}

func TestDoCheckForUpdatesSSH(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New(DefaultOptions())
	err := gen.doDownload("git@github.com:brainicorn/skelp.git", tmpDir)

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

	err = gen.checkForUpdates("git@github.com:brainicorn/skelp.git", tmpDir)

	if err != nil {
		t.Fatalf("error updating %s", err.Error())
	}
}

func TestDoDownloadHTTP(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New(DefaultOptions())
	err := gen.doDownload("https://github.com/brainicorn/skelp.git", tmpDir)

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

}

func TestDoCheckForUpdatesHTTP(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New(DefaultOptions())
	err := gen.doDownload("https://github.com/brainicorn/skelp.git", tmpDir)

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

	err = gen.checkForUpdates("https://github.com/brainicorn/skelp.git", tmpDir)

	if err != nil {
		t.Fatalf("error updating %s", err.Error())
	}
}
