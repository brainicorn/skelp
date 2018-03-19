package skelputil

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/sprig"
)

const (
	missingKeyOption = "missingkey=zero"
)

func FunctionMap() map[string]interface{} {
	fmap := sprig.FuncMap()
	fmap["cwd"] = func() string {
		s, _ := os.Getwd()
		return s
	}
	return fmap
}

func TemplateOptions() []string {
	return []string{missingKeyOption}
}

// Check if a file or directory exists.
func PathExists(path string) bool {
	// note: the err is either IsNotExist or something else
	// if it's something else, you have bigger issues...
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func DirIsEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

func MkdirAll(path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	return nil

}

func SafeFilenameFromPath(p string) string {
	filename := filepath.Clean(filepath.ToSlash(p))

	if strings.HasPrefix(filename, "/") || strings.HasPrefix(filename, ".") {
		filename = filename[1:]
	}

	if strings.HasSuffix(filename, "/") || strings.HasSuffix(filename, ".") {
		filename = filename[:len(filename)-1]
	}

	filename = strings.Replace(filename, " ", "_", -1)
	filename = strings.Replace(filename, "/", "-", -1)

	return filename
}

func GetFileMode(path string) (os.FileMode, error) {
	var err error
	var source *os.File
	var sourceInfo os.FileInfo
	var mode os.FileMode

	mode = os.ModePerm

	source, err = os.Open(path)

	if err == nil {
		defer source.Close()
		sourceInfo, err = source.Stat()
	}

	if err == nil {
		mode = sourceInfo.Mode()
	}

	return mode, err
}

func ListFilesByExtension(dir, ext string) []string {
	files := []string{}

	filepath.Walk(dir, func(curPath string, fi os.FileInfo, _ error) error {
		if !fi.IsDir() {
			if filepath.Ext(curPath) == ext {
				files = append(files, fi.Name())
			}
		}

		return nil
	})

	return files
}

func IsBlank(s string) bool {
	return len(strings.TrimSpace(s)) < 1
}
