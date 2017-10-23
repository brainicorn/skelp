package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/brainicorn/skelp/executor"
	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelputil"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	ErrBlankTemplateID           = "Template ID not provided"
	ErrTemplateRootNotFound      = "Template root not found %s"
	ErrSkelpTemplatesDirNotFound = "Skelp templates dir not found %s"
	ErrCacheNotFoundNoDownload   = "Cached template not found and downloads are turned off: %s"
)

func (sg *SkelpGenerator) Generate(templateID string, dataProvider provider.DataProvider) error {
	var err error

	if skelputil.IsBlank(templateID) {
		return fmt.Errorf(ErrBlankTemplateID)
	}

	switch TypeForTemplateID(templateID) {
	case TIDTypeAlias:
		err = sg.aliasGeneration(templateID, dataProvider)
	case TIDTypeFile:
		err = sg.pathGeneration(templateID, dataProvider)
	case TIDTypeRepo:
		err = sg.repoGeneration(templateID, dataProvider)
	}

	return err
}

func (sg *SkelpGenerator) pathGeneration(rootTemplateDir string, dataProvider provider.DataProvider) error {
	var err error
	var absRootTemplateDir string
	var skelpTemplatespath string
	var out string
	var tmplData interface{}

	absRootTemplateDir, err = filepath.Abs(rootTemplateDir)

	if err == nil {
		out = sg.skelpOptions.OutputDir
		skelpTemplatespath = filepath.Join(absRootTemplateDir, skelpTemplatesDirname)

		if !skelputil.PathExists(absRootTemplateDir) {
			err = fmt.Errorf(ErrTemplateRootNotFound, absRootTemplateDir)
		}
	}

	if err == nil && !skelputil.PathExists(skelpTemplatespath) {
		err = fmt.Errorf(ErrSkelpTemplatesDirNotFound, skelpTemplatespath)
	}

	if err == nil && skelputil.IsBlank(out) {
		out, err = os.Getwd()
	}

	if err == nil {
		tmplData, err = dataProvider(absRootTemplateDir)
	}

	if err == nil {
		if sg.skelpOptions.DryRun {
			var jsn []byte
			jsn, err = json.MarshalIndent(tmplData, "", "    ")

			if err == nil {
				fmt.Println("--------------")
				fmt.Println("Data Gathered:")
				fmt.Println("--------------")
				fmt.Println(string(jsn))
				fmt.Println("--------------")
			}

			return err
		}
	}

	if err == nil {
		skelpExec := executor.New(sg.funcMap, sg.tOptions)
		err = skelpExec.Execute(skelpTemplatespath, out, tmplData, sg.skelpOptions.OverwriteProvider)
	}

	return err
}

func (sg *SkelpGenerator) repoGeneration(templateID string, dataProvider provider.DataProvider) error {
	var err error
	var localTemplatePath string

	justDownloaded := false

	localTemplatePath, err = sg.absCacheDirFromURL(templateID)

	if err == nil {
		if !skelputil.PathExists(localTemplatePath) {
			if sg.skelpOptions.Download {
				err = sg.doDownload(templateID, localTemplatePath)
			} else {
				err = fmt.Errorf(ErrCacheNotFoundNoDownload, templateID)
			}
		}
	}

	if err == nil && !justDownloaded && sg.skelpOptions.CheckForUpdates {
		err = sg.checkForUpdates(templateID, localTemplatePath)
	}

	if err == nil {
		err = sg.pathGeneration(localTemplatePath, dataProvider)
	}

	return err
}

func (sg *SkelpGenerator) doDownload(u, path string) error {
	var err error
	am := AuthMethodForURL(u)

	opts := git.CloneOptions{
		URL:      u,
		Auth:     am,
		Progress: os.Stdout,
	}

	_, err = git.PlainClone(path, false, &opts)

	if err != nil && err == transport.ErrAuthenticationRequired {
		os.RemoveAll(path)

		// ask for authentication credentials and try again...
		if sg.skelpOptions.BasicAuthProvider != nil {
			user, pass := sg.skelpOptions.BasicAuthProvider()
			opts.Auth = http.NewBasicAuth(user, pass)
			_, err = git.PlainClone(path, false, &opts)
		}
	}

	return err
}

func (sg *SkelpGenerator) checkForUpdates(u, path string) error {
	var err error
	var repo *git.Repository
	var wt *git.Worktree

	am := AuthMethodForURL(u)
	repo, err = git.PlainOpen(path)

	if err == nil {
		opts := git.PullOptions{
			Auth:     am,
			Progress: os.Stdout,
		}

		wt, err = repo.Worktree()

		if err == nil {
			err = wt.Pull(&opts)
		}

		if err != nil {
			if err == transport.ErrAuthenticationRequired {
				// ask for authentication credentials and try again...
				if sg.skelpOptions.BasicAuthProvider != nil {
					user, pass := sg.skelpOptions.BasicAuthProvider()
					opts.Auth = http.NewBasicAuth(user, pass)
					err = wt.Pull(&opts)
				}
			}

			if err == git.NoErrAlreadyUpToDate {
				err = nil
			}
		}
	}

	return err

}

func (sg *SkelpGenerator) aliasGeneration(templateID string, dataProvider provider.DataProvider) error {
	var err error
	var aliasedTemplateID string

	aliasedTemplateID, err = sg.IDForAlias(templateID)

	if err == nil {
		err = sg.Generate(aliasedTemplateID, dataProvider)
	}

	return err
}
