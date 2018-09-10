package executor

import (
	"os"
	"strings"
	"testing"

	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelputil"
)

func TestEmptyTemplatesDir(t *testing.T) {
	exec := New(nil, []string{})
	err := exec.Execute("/tmp/noexist", "/tmp", nil, provider.DefaultOverwriteProvider, provider.DefaultExcludesProvider)

	if err == nil || !strings.HasPrefix(err.Error(), "No templates found in ") {
		t.Errorf("invalid err, want (%s), have (%s)", "No templates found in ", err.Error())
	}
}

func TestEmptyOutputDir(t *testing.T) {
	exec := New(nil, []string{})
	err := exec.Execute("../testdata/generator/simple", "", nil, provider.DefaultOverwriteProvider, provider.DefaultExcludesProvider)

	if err == nil || !strings.HasPrefix(err.Error(), "Output directory not provided") {
		t.Errorf("invalid err, want (%s), have (%s)", "Output directory not provided", err.Error())
	}
}

func TestOutputDirCreated(t *testing.T) {
	outputPath := "/tmp/skelpme"

	if skelputil.PathExists(outputPath) {
		t.Errorf("ouput path (%s) should not exist yet", outputPath)
	}

	defer os.RemoveAll(outputPath)

	testdata := map[string]string{"projectName": "testprj", "packageName": "testpkg"}

	exec := New(nil, []string{})
	exec.Execute("../testdata/generator/simple", outputPath, testdata, provider.DefaultOverwriteProvider, provider.DefaultExcludesProvider)

	if !skelputil.PathExists(outputPath) {
		t.Errorf("ouput path (%s) should have been created", outputPath)
	}
}
