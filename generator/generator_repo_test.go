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
	"github.com/brainicorn/skelp/skelputil"
)

var (
	readmeFmtRepo  = "README.md"
	projectFmtRepo = "%s.md"
	packageFmtRepo = "%s/%s.go"

	projectNameRepo    = "repogen"
	newProjectNameRepo = "newrepogen"
	packageNameRepo    = "repopack"

	readmeExpectedRepo    = "## " + projectNameRepo + " by brainicorn"
	newReadmeExpectedRepo = "## " + newProjectNameRepo + " by brainicorn"
	projectExpectedRepo   = projectNameRepo + " contains package " + packageNameRepo
	packageExpectedRepo   = "package " + packageNameRepo
)

func optionsForRepoTests() SkelpOptions {
	outDir, _ := ioutil.TempDir("", "skelp-repogen-test")
	homeDir, _ := ioutil.TempDir("", "skelp-custom-home")

	opts := DefaultOptions()
	opts.HomeDirOverride = homeDir
	opts.OutputDir = outDir

	return opts
}

func cleanOptions(opts SkelpOptions) {
	os.RemoveAll(opts.HomeDirOverride)
	os.RemoveAll(opts.OutputDir)
}

func TestRepoGenSimple(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	tmpDir := opts.OutputDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template.git", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtRepo)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(projectFmtRepo, projectNameRepo))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(packageFmtRepo, packageNameRepo, packageNameRepo))

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

	if string(readme) != readmeExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedRepo)
	}

	if string(prjfile) != projectExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpectedRepo)
	}

	if string(pkgfile) != packageExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpectedRepo)
	}

}

func TestRepoGenCWD(t *testing.T) {
	origCWD, _ := os.Getwd()
	tmpDir, _ := ioutil.TempDir("", "skelp-repogen-test")
	defer os.RemoveAll(tmpDir)

	homeDir, _ := ioutil.TempDir("", "skelp-custom-home")
	defer os.RemoveAll(homeDir)

	os.Chdir(tmpDir)
	defer os.Chdir(origCWD)

	opts := DefaultOptions()
	opts.HomeDirOverride = homeDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtRepo)
	projectPath := filepath.Join(tmpDir, fmt.Sprintf(projectFmtRepo, projectNameRepo))
	packagePath := filepath.Join(tmpDir, fmt.Sprintf(packageFmtRepo, packageNameRepo, packageNameRepo))

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

	if string(readme) != readmeExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedRepo)
	}

	if string(prjfile) != projectExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(prjfile), projectExpectedRepo)
	}

	if string(pkgfile) != packageExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(pkgfile), packageExpectedRepo)
	}

}

func TestRepoNotFound(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("git@github.com:brainicorn/does-not-exist", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "repository not found") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "repository not found")
	}
}

func TestRepoTemplatesFolderNotFound(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/notmplfolder", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "Skelp templates dir not found") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Skelp templates dir not found")
	}
}

func TestRepoGenBadTmpl(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/badtmpl", dp.DataProviderFunc)

	if err == nil {
		t.Error("expected error but was nil")
	}
}

func TestRepoMissingDescriptor(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("../testdata/generator/nodescriptor", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "skelp.json not found:") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "skelp.json not found:")
	}
}

func TestRepoNoOverwrite(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	tmpDir := opts.OutputDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtRepo)

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(readme) != readmeExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedRepo)
	}

	// run again with different data
	newData := map[string]interface{}{"projectName": newProjectNameRepo, "packageName": packageNameRepo}
	newDP := skelplate.NewDataProvider(newData)

	err = gen.Generate("https://github.com/brainicorn/skelp-test-template", newDP.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	newReadme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(newReadme) != readmeExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(newReadme), readmeExpectedRepo)
	}
}

func TestRepoOverwrite(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	tmpDir := opts.OutputDir

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmePath := filepath.Join(tmpDir, readmeFmtRepo)

	readme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(readme) != readmeExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(readme), readmeExpectedRepo)
	}

	// run again with different data
	newData := map[string]interface{}{"projectName": newProjectNameRepo, "packageName": packageNameRepo}
	newDP := skelplate.NewDataProvider(newData)

	opts.OverwriteProvider = func(rootDir, relFile string) bool { return true }
	gen.skelpOptions = opts
	err = gen.Generate("https://github.com/brainicorn/skelp-test-template", newDP.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	newReadme, err := ioutil.ReadFile(readmePath)

	if err != nil {
		t.Errorf("can't open out file (%s): %s", readmePath, err)
	}

	if string(newReadme) != newReadmeExpectedRepo {
		t.Errorf("contents don't match, have (%s), want (%s)", string(newReadme), newReadmeExpectedRepo)
	}
}

func TestRepoNoDownloadNoCache(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	opts.Download = false
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err == nil || !strings.HasPrefix(err.Error(), "Cached template not found and downloads are turned off:") {
		t.Errorf("wrong error: have (%s), want (%s)", err, "Cached template not found and downloads are turned off:")
	}

}

func TestRepoNoDownload(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	opts.Download = false
	gen2 := New(opts)

	err = gen2.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("cached generation error: %s", err)
	}
}

func TestReoSkelpDirOverride(t *testing.T) {
	opts := optionsForRepoTests()
	defer cleanOptions(opts)

	skelpDir, _ := ioutil.TempDir("", "custom-skelp-dir")
	defer os.RemoveAll(skelpDir)

	opts.SkelpDirOverride = skelpDir
	gen := New(opts)

	defData := map[string]interface{}{"projectName": projectNameRepo, "packageName": packageNameRepo}
	dp := skelplate.NewDataProvider(defData)

	err := gen.Generate("https://github.com/brainicorn/skelp-test-template", dp.DataProviderFunc)

	if err != nil {
		t.Errorf("generation error: %s", err)
	}

	readmeCachePath := filepath.Join(opts.HomeDirOverride, skelpDir, "gitcache", "github.com", "brainicorn", "skelp-test-template", readmeFmtRepo)

	if !skelputil.PathExists(readmeCachePath) {
		t.Errorf("cached readme not found: %s", readmeCachePath)
	}
}
