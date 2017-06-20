package generator

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
)

func (sg *SkelpGenerator) IDForAlias(alias string, options SkelpOptions) (string, error) {
	var err error
	var skelpHome string
	var templateID string
	var found bool

	if sg.aliases == nil {
		skelpHome, err = sg.initSkelpHome(options)

		if err == nil {
			aliasesPath := filepath.Join(skelpHome, skelpAliasesFilename)
			err = sg.loadAliases(aliasesPath)
		}
	}

	if err == nil {
		templateID, found = sg.aliases[alias]

		if !found {
			err = fmt.Errorf(ErrAliasNotFound, alias)
		}
	}

	return templateID, err
}

func (sg *SkelpGenerator) saveAliases(path string) error {
	var err error
	var file *os.File

	file, err = os.Create(path)

	if err == nil {
		encoder := gob.NewEncoder(file)
		sg.mu.Lock()
		encoder.Encode(sg.aliases)
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
