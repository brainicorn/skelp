package generator

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/brainicorn/skelp/skelputil"
)

const (
	ErrAliasNotFound        = "Alias '%s' not found in registry"
	ErrInvalidAlias         = "Invalid alias '%s'"
	ErrInvalidAliasTemplate = "Invalid alias template: '%s' must be a filepath or url"
)

type aliasRegistry map[string]aliasEntry

type aliasEntry struct {
	Name string
	Path string
}

func (sg *SkelpGenerator) IDForAlias(alias string) (string, error) {
	var err error
	var templateID string
	var entry aliasEntry
	var found bool

	if sg.aliases == nil {
		_, err = sg.initAliasRegistry()
	}

	if err == nil {
		entry, found = sg.aliases[alias]

		if !found {
			err = fmt.Errorf(ErrAliasNotFound, alias)
		} else {
			templateID = entry.Path
		}
	}

	return templateID, err
}

func (sg *SkelpGenerator) AddAlias(alias, fileOrUrl string) error {
	var err error
	var aliasFile string
	var aliasPath string

	if !IsAlias(alias) {
		return fmt.Errorf(ErrInvalidAlias, alias)
	}

	if !IsFilePath(fileOrUrl) && !IsRepoURL(fileOrUrl) {
		return fmt.Errorf(ErrInvalidAliasTemplate, fileOrUrl)
	}

	aliasFile, err = sg.initAliasRegistry()

	if err == nil {
		aliasPath = fileOrUrl

		if IsFilePath(aliasPath) && !strings.HasPrefix(aliasPath, "file://") {
			aliasPath, err = filepath.Abs(aliasPath)
		}

		if err == nil {
			sg.aliases[alias] = aliasEntry{Name: alias, Path: aliasPath}
		}
	}

	return sg.saveAliases(aliasFile, err)
}

func (sg *SkelpGenerator) RemoveAlias(alias string) error {
	var err error
	var aliasPath string

	if !IsAlias(alias) {
		return fmt.Errorf(ErrInvalidAlias, alias)
	}

	aliasPath, err = sg.initAliasRegistry()

	if err == nil {
		delete(sg.aliases, alias)
	}

	return sg.saveAliases(aliasPath, err)
}

func (sg *SkelpGenerator) AliasEntries() ([]aliasEntry, error) {
	var err error
	entries := []aliasEntry{}

	_, err = sg.initAliasRegistry()

	for _, entry := range sg.aliases {
		entries = append(entries, entry)
	}
	return entries, err
}

func (sg *SkelpGenerator) saveAliases(path string, initialErr error) error {
	var err error
	var file *os.File
	records := [][]string{{}}

	aliases := aliasRegistry{}
	err = initialErr

	if err == nil {
		if sg.aliases != nil {
			aliases = sg.aliases
		}

		for _, alias := range aliases {
			records = append(records, []string{alias.Name, alias.Path})
		}

		file, err = os.Create(path)
	}

	if err == nil {
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		sg.mu.Lock()
		err = writer.WriteAll(records)
		sg.mu.Unlock()
	}

	return err
}

func (sg *SkelpGenerator) loadAliases(path string) error {
	var err error
	var file *os.File
	var records [][]string

	if sg.aliases == nil {
		file, err = os.Open(path)
		if err == nil {
			defer file.Close()
			reader := csv.NewReader(file)
			sg.mu.Lock()
			records, err = reader.ReadAll()
			sg.mu.Unlock()

			if err == nil {
				sg.aliases = make(map[string]aliasEntry)

				for _, record := range records {
					entry := aliasEntry{
						Name: record[0],
						Path: record[1],
					}

					sg.aliases[record[0]] = entry
				}
			}
		}
	}

	return err
}

func (sg *SkelpGenerator) initAliasRegistry() (string, error) {
	var err error
	var skelpHome, aliasesPath string

	skelpHome, err = sg.InitSkelpHome()

	if err == nil {
		aliasesPath = filepath.Join(skelpHome, skelpAliasesFilename)

		if !skelputil.PathExists(aliasesPath) {
			err = sg.saveAliases(aliasesPath, err)
		}

		if err == nil {
			err = sg.loadAliases(aliasesPath)
		}

	}
	return aliasesPath, err

}
