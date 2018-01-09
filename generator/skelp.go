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
	skelpAliasesFilename      = "aliases.cfg"
	skelpTemplateCacheDirname = "gitcache"
)

type SkelpOptions struct {
	Download          bool
	CheckForUpdates   bool
	DryRun            bool
	QuietMode         bool
	OutputDir         string
	HomeDirOverride   string
	SkelpDirOverride  string
	OverwriteProvider provider.OverwriteProvider
	BasicAuthProvider provider.BasicAuthProvider
	HookProvider      provider.HookProvider
	ExcludesProvider  provider.ExcludesProvider
	ReplayProvider    provider.ReplayProvider
}

func DefaultOptions() SkelpOptions {
	bap := &provider.DefaultBasicAuthProvider{}

	return SkelpOptions{
		Download:          true,
		CheckForUpdates:   true,
		OverwriteProvider: provider.DefaultOverwriteProvider,
		BasicAuthProvider: bap.ProvideAuth,
		HookProvider:      provider.DefaultHookProvider,
		ExcludesProvider:  provider.DefaultExcludesProvider,
		ReplayProvider:    &provider.DefaultReplayProvider{},
	}
}

type SkelpGenerator struct {
	SkelpHome    string
	funcMap      map[string]interface{}
	tOptions     []string
	skelpOptions SkelpOptions
	aliases      aliasRegistry
	mu           sync.Mutex
}

func New(options SkelpOptions) *SkelpGenerator {
	return &SkelpGenerator{
		skelpOptions: options,
		funcMap:      skelputil.FunctionMap(),
		tOptions:     skelputil.TemplateOptions(),
	}
}

func (sg *SkelpGenerator) absCacheDirFromURL(u string) (string, error) {
	var err error
	var skelpHome, cacheDir, templateDir, absDir string

	skelpHome, err = sg.InitSkelpHome()

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

func (sg *SkelpGenerator) InitSkelpHome() (string, error) {
	var err error
	var homeDir string
	var path string

	if len(strings.TrimSpace(sg.SkelpHome)) > 0 && skelputil.PathExists(sg.SkelpHome) {
		return sg.SkelpHome, nil
	}

	skelpDir := defaultSkelpDir
	if len(strings.TrimSpace(sg.skelpOptions.SkelpDirOverride)) > 0 {
		skelpDir = sg.skelpOptions.SkelpDirOverride
	}

	if len(strings.TrimSpace(sg.skelpOptions.HomeDirOverride)) > 0 {
		homeDir = sg.skelpOptions.HomeDirOverride
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
