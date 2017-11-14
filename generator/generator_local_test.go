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
	readmeFmtLocal   = "README.md"
	projectFmtLocal  = "%s.md"
	packageFmtLocal  = "%s/%s.go"
	preInputFmtLocal = "preinput.txt"
	preGenFmtLocal   = "pregen.txt"
	postGenFmtLocal  = "postgen.txt"

	projectNameLocal    = "localgen"
	newProjectNameLocal = "newlocalgen"
	packageNameLocal    = "localpack"

	readmeExpectedLocal        = "## " + projectNameLocal + " by brainicorn"
	newReadmeExpectedLocal     = "## " + newProjectNameLocal + " by brainicorn"
	readmeComplexExpectedLocal = `## %s by brainicorn
This project uses the following database:

mongo namespace: myspace in regions: east,ap
`

	projectExpectedLocal = projectNameLocal + " contains package " + packageNameLocal
	packageExpectedLocal = "package " + packageNameLocal

	preinputExpectedLocal = "Greetings yo yo yo"
	pregenExpectedLocal   = "Greetings preBob"
	postgenExpectedLocal  = "Greetings postBob"
)

func TestLocalGenSimple(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal}
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	newDP := skelplate.NewDataProvider(newData, 0)

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
	dp := skelplate.NewDataProvider(defData, 0)

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
	newDP := skelplate.NewDataProvider(newData, 0)

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

func TestLocalGenComplex(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "skelp-localcomplex-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)
	defData := map[string]interface{}{"projectName": projectNameLocal, "packageName": packageNameLocal, "multiComplex": []map[string]interface{}{{"varone": "foo", "vartwo": "bar"}, {"varone": "foo2", "vartwo": "bar2"}}, "database": map[string]interface{}{"db": "mongo", "namespace": "myspace", "regions": []interface{}{"east", "ap"}}}
	dp := skelplate.NewDataProvider(defData, skelplate.SkipMulti)

	err := gen.Generate("../testdata/generator/complex", dp.DataProviderFunc)

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

	expected := fmt.Sprintf(readmeComplexExpectedLocal, projectNameLocal)
	if string(readme) != expected {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), expected)
	}

	if string(prjfile) != projectExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpectedLocal)
	}

	if string(pkgfile) != packageExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpectedLocal)
	}

}

func TestLocalHooks(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"preGreeting": "preBob", "postGreeting": "postBob"}
	dp := skelplate.NewDataProvider(defData, 0)

	err := gen.GenerateWithHooks("../testdata/generator/simplehooks", dp.DataProviderFunc, dp.HookProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	preinputPath := filepath.Join(tmpDir, preInputFmtLocal)
	pregenPath := filepath.Join(tmpDir, preGenFmtLocal)
	postgenPath := filepath.Join(tmpDir, postGenFmtLocal)

	preinput, err := ioutil.ReadFile(preinputPath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", preinputPath, err)
	}

	pregen, err := ioutil.ReadFile(pregenPath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", pregenPath, err)
	}

	postgen, err := ioutil.ReadFile(postgenPath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", postgenPath, err)
	}

	if strings.TrimSpace(string(preinput)) != preinputExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(preinput), preinputExpectedLocal)
	}

	if strings.TrimSpace(string(pregen)) != pregenExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pregen), pregenExpectedLocal)
	}

	if strings.TrimSpace(string(postgen)) != postgenExpectedLocal {
		t.Errorf("contents don't match, have (%s), want (%s)", string(postgen), postgenExpectedLocal)
	}

}

func TestPreInHookErr(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"preGreeting": "preBob", "postGreeting": "postBob"}
	dp := skelplate.NewDataProvider(defData, 0)

	err := gen.GenerateWithHooks("../testdata/generator/preinhookserr", dp.DataProviderFunc, dp.HookProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "error executing preInput hook") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "error executing preInput hook")
	}

}

func TestPreGenHookErr(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"preGreeting": "preBob", "postGreeting": "postBob"}
	dp := skelplate.NewDataProvider(defData, 0)

	err := gen.GenerateWithHooks("../testdata/generator/pregenhookserr", dp.DataProviderFunc, dp.HookProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "error executing preGen hook") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "error executing preGen hook")
	}

}

func TestPostGenHookErr(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
	defer os.RemoveAll(tmpDir)

	opts := DefaultOptions()
	opts.OutputDir = tmpDir

	gen := New(opts)

	defData := map[string]interface{}{"preGreeting": "preBob", "postGreeting": "postBob"}
	dp := skelplate.NewDataProvider(defData, 0)

	err := gen.GenerateWithHooks("../testdata/generator/postgenhookserr", dp.DataProviderFunc, dp.HookProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "error executing postGen hook") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "error executing postGen hook")
	}

}
