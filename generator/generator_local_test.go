package generator

// TODO refactor to reuse common code
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
	readmeFmt  = "README.md"
	projectFmt = "%s.md"
	packageFmt = "%s/%s.go"

	projectName    = "localgen"
	newProjectName = "newlocalgen"
	packageName    = "localpack"

	readmeExpected    = "## " + projectName + " by brainicorn"
	newReadmeExpected = "## " + newProjectName + " by brainicorn"
	projectExpected   = projectName + " contains package " + packageName
	packageExpected   = "package " + packageName
)

func TestLocalGenSimple(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmt)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(projectFmt, projectName))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(packageFmt, packageName, packageName))

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

	if string(readme) != readmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpected)
	}

	if string(prjfile) != projectExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpected)
	}

	if string(pkgfile) != packageExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpected)
	}

}

func TestLocalGenCWD(t *testing.T) {
	origCWD, _ := os.Getwd()
	absTemplateDir, _ := filepath.Abs("../testdata/generator/simple")
	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	os.Chdir(tmpDir)
	defer os.Chdir(origCWD)

	opts := DefaultOptions()
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate(absTemplateDir, dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmt)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(projectFmt, projectName))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(packageFmt, packageName, packageName))

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

	if string(readme) != readmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpected)
	}

	if string(prjfile) != projectExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpected)
	}

	if string(pkgfile) != packageExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpected)
	}

}

func TestLocalBlankTemplateID(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-badtmpl-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("", dp.DataProviderFunc)

	if err == nil || err.Error() != "Template ID not provided" {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Template ID not provided")
	}
}

func TestLocalTemplateRootNotFound(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-badtmpl-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("/does/not/exist", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "Template root not found") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Template root not found")
	}
}

func TestLocalTemplatesFolderNotFound(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-badtmpl-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/notmplfolder", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "Skelp templates dir not found") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Skelp templates dir not found")
	}
}

func TestLocalGenBadTmpl(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-badtmpl-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/badtmpl", dp.DataProviderFunc)

	if err == nil {
		t.Error("expected error but was nil")
	}
}

func TestLocalMissingDescriptor(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-missing-descriptor-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/nodescriptor", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "skelp.json not found:") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "skelp.json not found:")
	}
}

func TestNoOverwrite(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-nooverwrite-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmt)

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(readme) != readmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpected)
	}

	// run again with different data
	newData := map[string]interface{}{"projectName": newProjectName, "packageName": packageName}
	newDP := skelplate.NewDataProvider(newData)

	err = gen.Generate("../testdata/generator/simple", newDP.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	newReadme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(newReadme) != readmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(newReadme), readmeExpected)
	}
}

func TestOverwrite(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-overwrite-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectName, "packageName": packageName}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmt)

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(readme) != readmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpected)
	}

	// run again with different data
	newData := map[string]interface{}{"projectName": newProjectName, "packageName": packageName}
	newDP := skelplate.NewDataProvider(newData)

	opts.OverwriteProvider = func(rootDir, relFile string) bool { return true }
	gen.skelpOptions = opts
	err = gen.Generate("../testdata/generator/simple", newDP.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	newReadme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(newReadme) != newReadmeExpected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(newReadme), newReadmeExpected)
	}
}
