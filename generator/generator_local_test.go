package generator

// TODO refactor to reuse common code
import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelplate"
)

var (
	readmeFmtLocal  = "README.md"
	projectFmtLocal = "%s.md"
	packageFmtLocal = "%s/%s.go"

	projectNameLocal    = "localgen"
	newProjectNameLocal = "newlocalgen"
	packageNameLocal    = "localpack"

	readmeExpectedLocal    = "## " + projectNameLocal + " by brainicorn"
	newReadmeExpectedLocal = "## " + newProjectNameLocal + " by brainicorn"
	projectExpectedLocal   = projectNameLocal + " contains package " + packageNameLocal
	packageExpectedLocal   = "package " + packageNameLocal
)

func TestLocalGenSimple(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtLocal)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(projectFmtLocal, projectNameLocal))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(packageFmtLocal, packageNameLocal, packageNameLocal))

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

	if string(readme) != readmeExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedLocal)
	}

	if string(prjfile) != projectExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpectedLocal)
	}

	if string(pkgfile) != packageExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpectedLocal)
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

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate(absTemplateDir, dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtLocal)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(projectFmtLocal, projectNameLocal))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(packageFmtLocal, packageNameLocal, packageNameLocal))

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

	if string(readme) != readmeExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedLocal)
	}

	if string(prjfile) != projectExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpectedLocal)
	}

	if string(pkgfile) != packageExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpectedLocal)
	}

}

func TestLocalBlankTemplateID(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-badtmpl-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
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

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
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

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
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

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
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

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
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

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtLocal)

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(readme) != readmeExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedLocal)
	}

	// run again with different data
	newData := map[string]interface{}{"projectName": newProjectNameLocal, "packageName": packageNameLocal}
	newDP := skelplate.NewDataProvider(newData)

	err = gen.Generate("../testdata/generator/simple", newDP.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	newReadme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(newReadme) != readmeExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(newReadme), readmeExpectedLocal)
	}
}

func TestOverwrite(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-overwrite-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtLocal)

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(readme) != readmeExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedLocal)
	}

	// run again with different data
	newData := map[string]interface{}{"projectName": newProjectNameLocal, "packageName": packageNameLocal}
	newDP := skelplate.NewDataProvider(newData)

	opts.OverwriteProvider = provider.AlwaysOverwriteProvider
	gen.skelpOptions = opts
	err = gen.Generate("../testdata/generator/simple", newDP.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	newReadme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(newReadme) != newReadmeExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(newReadme), newReadmeExpectedLocal)
	}
}
