package generator

import (
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
	ErrSkelpFileNotFound         = "Skelp descriptor not found %s"
	ErrSkelpTemplatesDirNotFound = "Skelp templates dir not found %s"
	ErrCacheNotFoundNoDownload   = "Remote template not found and downloads are turned off: %s"
)

func (sg *SkelpGenerator) Generate(templateID string, dataProvider provider.DataProvider, options SkelpOptions) error {
	var err error

	if skelputil.IsBlank(templateID) {
		return fmt.Errorf(ErrBlankTemplateID)
	}

	switch TypeForTemplateID(templateID) {
	case TIDTypeAlias:
		err = sg.aliasGeneration(templateID, dataProvider, options)
	case TIDTypeFile:
		err = sg.pathGeneration(templateID, dataProvider, options)
	case TIDTypeRepo:
		err = sg.repoGeneration(templateID, dataProvider, options)
	}

	return err
}

func (sg *SkelpGenerator) pathGeneration(rootTemplateDir string, dataProvider provider.DataProvider, options SkelpOptions) error {
	var err error
	var tmplData interface{}

	out := options.OutputDir
	skelpTemplatespath := filepath.Join(rootTemplateDir, skelpTemplatesDirname)

	if !skelputil.PathExists(rootTemplateDir) {
		err = fmt.Errorf(ErrTemplateRootNotFound, rootTemplateDir)
	}

	if err == nil && !skelputil.PathExists(skelpTemplatespath) {
		err = fmt.Errorf(ErrSkelpTemplatesDirNotFound, skelpTemplatespath)
	}

	if err == nil && skelputil.IsBlank(out) {
		out, err = os.Getwd()
	}

	if err == nil {
		tmplData, err = dataProvider(rootTemplateDir)
	}

	if err == nil {
		skelpExec := executor.New(sg.funcMap, sg.tOptions)
		err = skelpExec.Execute(skelpTemplatespath, out, tmplData, options.OverwriteProvider)
	}

	return err
}

func (sg *SkelpGenerator) repoGeneration(templateID string, dataProvider provider.DataProvider, options SkelpOptions) error {
	var err error
	var localTemplatePath string

	justDownloaded := false

	localTemplatePath, err = sg.absCacheDirFromURL(templateID, options)

	if err == nil {
		if !skelputil.PathExists(localTemplatePath) {
			if options.Download {
				err = sg.doDownload(templateID, localTemplatePath, options)
			} else {
				err = fmt.Errorf(ErrCacheNotFoundNoDownload, templateID)
			}
		}
	}

	if err == nil && !justDownloaded && options.CheckForUpdates {
		err = sg.checkForUpdates(templateID, localTemplatePath, options)
	}

	if err == nil {
		err = sg.pathGeneration(localTemplatePath, dataProvider, options)
	}

	return err
}

func (sg *SkelpGenerator) doDownload(u, path string, options SkelpOptions) error {
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
		if options.BasicAuthProvider != nil {
			user, pass := options.BasicAuthProvider()
			opts.Auth = http.NewBasicAuth(user, pass)
			_, err = git.PlainClone(path, false, &opts)
		}
	}

	return err
}

func (sg *SkelpGenerator) checkForUpdates(u, path string, options SkelpOptions) error {
	var err error
	var repo *git.Repository

	am := AuthMethodForURL(u)
	repo, err = git.PlainOpen(path)

	if err == nil {
		opts := git.PullOptions{
			Auth:     am,
			Progress: os.Stdout,
		}
		err = repo.Pull(&opts)

		if err != nil {
			if err == transport.ErrAuthenticationRequired {
				// ask for authentication credentials and try again...
				if options.BasicAuthProvider != nil {
					user, pass := options.BasicAuthProvider()
					opts.Auth = http.NewBasicAuth(user, pass)
					err = repo.Pull(&opts)
				}
			}

			if err == git.NoErrAlreadyUpToDate {
				err = nil
			}
		}
	}

	return err

}

func (sg *SkelpGenerator) aliasGeneration(templateID string, dataProvider provider.DataProvider, options SkelpOptions) error {
	var err error
	var aliasedTemplateID string

	aliasedTemplateID, err = sg.IDForAlias(templateID, options)

	if err == nil {
		err = sg.Generate(aliasedTemplateID, dataProvider, options)
	}

	return err
}
