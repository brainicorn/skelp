package util

import (
	"io"
	"os"
)

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
