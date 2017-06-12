package generator

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelputil"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	defaultSkelpDir           = ".skelp"
	skelpTemplatesDirname     = "templates"
	skelpAliasesFilename      = "aliases.gob"
	skelpTemplateCacheDirname = "gitcache"

	ErrAliasNotFound = "Alias '%s' not found in registry"
)

type SkelpOptions struct {
	Download          bool
	CheckForUpdates   bool
	OutputDir         string
	HomeOverride      string
	OverwriteProvider provider.OverwriteProvider
	BasicAuthProvider provider.BasicAuthProvider
}

func DefaultOptions() SkelpOptions {
	return SkelpOptions{
		Download:          true,
		CheckForUpdates:   true,
		OverwriteProvider: provider.DefaultOverwriteProvider,
	}
}

type aliasRegistry map[string]string

type SkelpGenerator struct {
	SkelpHome string
	funcMap   map[string]interface{}
	tOptions  []string
	aliases   aliasRegistry
	mu        sync.Mutex
}

func New() *SkelpGenerator {
	return &SkelpGenerator{
		funcMap:  skelputil.FunctionMap(),
		tOptions: skelputil.TemplateOptions(),
	}
}

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
func (sg *SkelpGenerator) absCacheDirFromURL(u string, options SkelpOptions) (string, error) {
	var err error
	var skelpHome, cacheDir, templateDir, absDir string

	skelpHome, err = sg.initSkelpHome(options)

	if err == nil {
		cacheDir = filepath.Join(skelpHome, skelpTemplateCacheDirname)
		if !skelputil.PathExists(cacheDir) {
			err = os.MkdirAll(cacheDir, os.ModePerm)
		}
	}

	if err == nil {
		templateDir, err = FilepathFromURL(u)
	}

	if err == nil {
		absDir = filepath.Join(cacheDir, templateDir)
	}

	return absDir, err
}

func (sg *SkelpGenerator) initSkelpHome(options SkelpOptions) (string, error) {
	var err error
	var homeDir string
	var path string

	skelpDir := defaultSkelpDir
	if len(strings.TrimSpace(options.HomeOverride)) > 0 {
		skelpDir = options.HomeOverride
	}

	homeDir, err = homedir.Dir()
	if err == nil {
		path = filepath.Join(homeDir, skelpDir)

		if !skelputil.PathExists(path) {
			err = os.MkdirAll(path, os.ModePerm)
		}
	}

	if err == nil {
		sg.mu.Lock()
		sg.SkelpHome = path
		sg.mu.Unlock()
	}

	return path, err
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
