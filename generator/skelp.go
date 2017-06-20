package generator

import (
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
	HomeDirOverride   string
	SkelpDirOverride  string
	OverwriteProvider provider.OverwriteProvider
	BasicAuthProvider provider.BasicAuthProvider
}

func DefaultOptions() SkelpOptions {
	return SkelpOptions{
		Download:          true,
		CheckForUpdates:   true,
		OverwriteProvider: provider.DefaultOverwriteProvider,
		BasicAuthProvider: provider.DefaultBasicAuthProvider,
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
	if len(strings.TrimSpace(options.SkelpDirOverride)) > 0 {
		skelpDir = options.SkelpDirOverride
	}

	if len(strings.TrimSpace(options.HomeDirOverride)) > 0 {
		homeDir = options.HomeDirOverride
	} else {
		homeDir, err = homedir.Dir()
	}

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
