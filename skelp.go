package skelp

import (
	"os"
	"path/filepath"

	"github.com/Masterminds/sprig"
	"github.com/brainicorn/skelp/util"
	homedir "github.com/mitchellh/go-homedir"
)

//go:generate go-bindata -pkg template -o ./template/bindata.go ./template/schema
const (
	skelpDir         = ".skelp"
	missingKeyOption = "missingkey=zero"
)

type Skelp struct {
	SkelpHome string
	fMap      map[string]interface{}
}

func NewSkelp() (*Skelp, error) {
	skelpHome, err := initSkelpHome()
	if err == nil {
		return &Skelp{
			SkelpHome: skelpHome,
			fMap:      initFunctionMap(),
		}, nil
	}

	return nil, err
}

func initSkelpHome() (string, error) {
	var err error
	var homeDir string
	var path string

	homeDir, err = homedir.Dir()
	if err == nil {
		path = filepath.Join(homeDir, skelpDir)
		if !util.PathExists(path) {
			err = os.MkdirAll(path, os.ModePerm)
		}
	}

	return path, err
}

func initFunctionMap() map[string]interface{} {
	fmap := sprig.FuncMap()

	return fmap
}
