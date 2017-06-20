package generator

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brainicorn/skelp/skelputil"
)

const (
	ErrAliasNotFound        = "Alias '%s' not found in registry"
	ErrInvalidAlias         = "Invalid alias '%s'"
	ErrInvalidAliasTemplate = "Invalid alias template: '%s' must be a filepath or url"
)

type aliasRegistry map[string]string

func (sg *SkelpGenerator) IDForAlias(alias string) (string, error) {
	var err error
	var templateID string
	var found bool

	if sg.aliases == nil {
		_, err = sg.initAliasRegistry()
	}

	if err == nil {
		templateID, found = sg.aliases[alias]

		if !found {
			err = fmt.Errorf(ErrAliasNotFound, alias)
		}
	}

	return templateID, err
}

func (sg *SkelpGenerator) AddAlias(alias, fileOrUrl string) error {
	var err error
	var aliasPath string

	if !IsAlias(alias) {
		return fmt.Errorf(ErrInvalidAlias, alias)
	}

	if !IsFilePath(fileOrUrl) && !IsRepoURL(fileOrUrl) {
		return fmt.Errorf(ErrInvalidAliasTemplate, fileOrUrl)
	}

	aliasPath, err = sg.initAliasRegistry()

	if err == nil {
		sg.aliases[alias] = fileOrUrl
	}

	return sg.saveAliases(aliasPath)
}

func (sg *SkelpGenerator) saveAliases(path string) error {
	var err error
	var file *os.File
	aliases := aliasRegistry{}

	if sg.aliases != nil {
		aliases = sg.aliases
	}

	file, err = os.Create(path)

	if err == nil {
		encoder := gob.NewEncoder(file)
		sg.mu.Lock()
		encoder.Encode(aliases)
		sg.mu.Unlock()
	}

	file.Close()
	return err
}

func (sg *SkelpGenerator) loadAliases(path string) error {
	var err error
	var file *os.File

	if sg.aliases == nil {
		file, err = os.Open(path)
		if err == nil {
			decoder := gob.NewDecoder(file)
			sg.mu.Lock()
			err = decoder.Decode(&sg.aliases)
			sg.mu.Unlock()
		}
		file.Close()
	}

	return err
}

func (sg *SkelpGenerator) initAliasRegistry() (string, error) {
	var err error
	var skelpHome, aliasesPath string

	skelpHome, err = sg.initSkelpHome()

	if err == nil {
		aliasesPath = filepath.Join(skelpHome, skelpAliasesFilename)

		if !skelputil.PathExists(aliasesPath) {
			err = sg.saveAliases(aliasesPath)
		}

		if err == nil {
			err = sg.loadAliases(aliasesPath)
		}

	}
	return aliasesPath, err

}
