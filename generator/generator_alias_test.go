package generator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brainicorn/skelp/skelplate"
)

var (
	aliasreadmeFmt  = "README.md"
	aliasprojectFmt = "%s.md"
	aliaspackageFmt = "%s/%s.go"

	aliasprojectName    = "localgen"
	aliasnewProjectName = "newlocalgen"
	aliaspackageName    = "localpack"

	aliasreadmeExpected    = "## " + aliasprojectName + " by brainicorn"
	aliasnewReadmeExpected = "## " + aliasnewProjectName + " by brainicorn"
	aliasprojectExpected   = aliasprojectName + " contains package " + aliaspackageName
	aliaspackageExpected   = "package " + aliaspackageName
)

func TestUnknownAlias(t *testing.T) {
	opts := DefaultOptions()
	opts.OutputDir = "/tmp"

	defData := map[string]interface{}{"projectName": aliasprojectName, "packageName": aliaspackageName}
	dp := skelplate.NewDataProvider(defData)

	gen := New(opts)

	err := gen.Generate("simplealias", dp.DataProviderFunc)

	if err == nil || !strings.HasSuffix(err.Error(), "not found in registry") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "not found in registry")
	}

}

func TestAddInvalidAlias(t *testing.T) {
	opts := DefaultOptions()
	opts.OutputDir = "/tmp"

	gen := New(opts)

	err := gen.AddAlias("http://some.url", "https://github.com/somerepo")

	if err == nil || !strings.HasPrefix(err.Error(), "Invalid alias ") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Invalid alias ")
	}

}

func TestAddInvalidTemplate(t *testing.T) {
	opts := DefaultOptions()
	opts.OutputDir = "/tmp"

	gen := New(opts)

	err := gen.AddAlias("somealias", "justsomerandomtext")

	if err == nil || !strings.HasPrefix(err.Error(), "Invalid alias template:") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Invalid alias template:")
	}

}

func TestAdHocAliasFileCreation(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.HomeDirOverride = tmpDir

	gen := New(opts)

	_, err := gen.initAliasRegistry()

	if err != nil {
		t.Error("should not have gotten error")
	}

}

func TestAliasGenSimple(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-aliasgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)
	defData := map[string]interface{}{"projectName": aliasprojectName, "packageName": aliaspackageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.AddAlias("knownalias", "../testdata/generator/simple")

	if err != nil {
		t.Fatalf("alias error: %s", err)
	}

	err = gen.Generate("knownalias", dp.DataProviderFunc)

	if err != nil {
		t.Fatalf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, aliasreadmeFmt)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(aliasprojectFmt, aliasprojectName))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(aliaspackageFmt, aliaspackageName, aliaspackageName))

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	prjfile, err := ioutil.ReadFile(projectPath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", projectPath, err)
	}

	pkgfile, err := ioutil.ReadFile(packagePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", packagePath, err)
	}

	if string(readme) != aliasreadmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), aliasreadmeExpected)
	}

	if string(prjfile) != aliasprojectExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), aliasprojectExpected)
	}

	if string(pkgfile) != aliaspackageExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), aliaspackageExpected)
	}

}
