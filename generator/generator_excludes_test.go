package generator

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/brainicorn/skelp/skelplate"
	"github.com/brainicorn/skelp/skelputil"
)

var (
	projectFmtExclude  = "%s.md"
	packageFmtExclude  = "%s/%s.go"
	projectNameExclude = "localgen"
	packageNameExclude = "localpack"
)

var exldTests = []struct {
	descriptor    string
	expected      map[string]bool
	expectedError error
}{
	{
		"exclude-single-file.json",
		map[string]bool{
			"README.md": false,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     true,
			packageNameExclude:                                                     true,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): true,
		},
		nil,
	},
	{
		"exclude-multi-file.json",
		map[string]bool{
			"README.md": false,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     true,
			packageNameExclude:                                                     true,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): false,
		},
		nil,
	},
	{
		"exclude-dir.json",
		map[string]bool{
			"README.md": true,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     true,
			packageNameExclude:                                                     false,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): false,
		},
		nil,
	},
	{
		"exclude-bad-condition.json",
		map[string]bool{
			"README.md": false,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     false,
			packageNameExclude:                                                     false,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): false,
		},
		errors.New("whoops!"),
	},
	{
		"exclude-negative-condition.json",
		map[string]bool{
			"README.md": true,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     true,
			packageNameExclude:                                                     true,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): true,
		},
		nil,
	},
	{
		"exclude-multiple.json",
		map[string]bool{
			"README.md": false,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     true,
			packageNameExclude:                                                     true,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): false,
		},
		nil,
	},
	{
		"exclude-template-condition.json",
		map[string]bool{
			"README.md": false,
			fmt.Sprintf(projectFmtExclude, projectNameExclude):                     true,
			packageNameExclude:                                                     true,
			fmt.Sprintf(packageFmtExclude, packageNameExclude, packageNameExclude): true,
		},
		nil,
	},
}

func TestExcludeExcludesSimple(t *testing.T) {
	for _, et := range exldTests {
		fmt.Println("")
		fmt.Println("---------------------------")
		fmt.Println(et.descriptor)
		fmt.Println("---------------------------")

		tmpDir, _ := ioutil.TempDir("", "skelp-localgen-test")
		defer os.RemoveAll(tmpDir)

		defData := map[string]interface{}{"projectName": projectNameExclude, "packageName": packageNameExclude}
		dp := skelplate.NewDataProvider(defData, 0)
		dp.OverrideSkelpFilename(et.descriptor)

		opts := DefaultOptions()
		opts.OutputDir = tmpDir
		opts.ExcludesProvider = dp.ExcludesProviderFunc

		gen := New(opts)

		err := gen.Generate("../testdata/generator/simple", dp.DataProviderFunc)

		if et.expectedError != nil && err == nil {
			t.Fatalf("[%s] expected error: %s", et.descriptor, et.expectedError)
		}

		if et.expectedError == nil && err != nil {
			t.Fatalf("[%s] generation error: %s", et.descriptor, err)
		}

		for expectedPath, shouldExist := range et.expected {
			actualPath := filepath.Join(tmpDir, expectedPath)
			actualExists := skelputil.PathExists(actualPath)

			fmt.Println(fmt.Sprintf("expectedPath: %s, shouldExist: %t", expectedPath, shouldExist))
			fmt.Println(fmt.Sprintf("actualPath: %s, exists: %t", actualPath, actualExists))

			if shouldExist != actualExists {
				modifier := "should"
				if !shouldExist {
					modifier = "should not"
				}

				t.Fatalf("[%s] path %s %s exist!", et.descriptor, expectedPath, modifier)
			}
		}

	}

}
