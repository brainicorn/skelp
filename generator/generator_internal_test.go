package generator

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelputil"
	"github.com/joho/godotenv"
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

type baUser struct {
	in         *os.File
	keystrokes []string
}

func newBAUser(keystrokes []string) *baUser {
	in, _ := ioutil.TempFile("", "")
	os.Stdin = in

	return &baUser{
		in:         in,
		keystrokes: keystrokes,
	}
}

func (f *baUser) nextKeystroke() {
	var keystroke string
	keystroke, f.keystrokes = f.keystrokes[0], f.keystrokes[1:]
	f.in.Truncate(0)
	f.in.Seek(0, os.SEEK_SET)
	io.WriteString(f.in, keystroke+"\n")
	f.in.Seek(0, os.SEEK_SET)
}

func (f *baUser) done() {
	f.in.Close()
}

func TestDoDownloadBasicAuth(t *testing.T) {
	bauser := os.Getenv("BAUSER")
	bapass := os.Getenv("BAPASS")

	if skelputil.IsBlank(bauser) {
		godotenv.Load("../.env")
		bauser = os.Getenv("BAUSER")
		bapass = os.Getenv("BAPASS")
	}

	user := newBAUser([]string{bauser, bapass})
	defer user.done()

	opts := DefaultOptions()
	baprovider := provider.DefaultBasicAuthProvider{BeforePrompt: user.nextKeystroke}
	opts.BasicAuthProvider = baprovider.ProvideAuth

	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New(opts)
	err := gen.doDownload("https://bitbucket.org/brainicorn/skelp-basicauth-test.git", tmpDir)

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

}

func TestDoCheckForUpdatesBasicAuth(t *testing.T) {
	bauser := os.Getenv("BAUSER")
	bapass := os.Getenv("BAPASS")

	if skelputil.IsBlank(bauser) {
		godotenv.Load("../.env")
		bauser = os.Getenv("BAUSER")
		bapass = os.Getenv("BAPASS")
	}

	user := newBAUser([]string{bauser, bapass, bauser, bapass})
	defer user.done()

	opts := DefaultOptions()
	baprovider := provider.DefaultBasicAuthProvider{BeforePrompt: user.nextKeystroke}
	opts.BasicAuthProvider = baprovider.ProvideAuth

	tmpDir, _ := ioutil.TempDir("", "skelp-git-test")
	defer os.RemoveAll(tmpDir)
	gen := New(opts)
	err := gen.doDownload("https://bitbucket.org/brainicorn/skelp-basicauth-test.git", tmpDir)

	if err != nil {
		t.Fatalf("error downloading %s", err.Error())
	}

	dpath := filepath.Join(tmpDir, "README.md")
	if !skelputil.PathExists(dpath) {
		t.Fatalf("download path doesn't exist %s", dpath)
	}

	err = gen.checkForUpdates("https://bitbucket.org/brainicorn/skelp-basicauth-test.git", tmpDir)

	if err != nil {
		t.Fatalf("error updating %s", err.Error())
	}
}
