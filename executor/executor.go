package executor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelputil"
)

const (
	ErrNoTemplatesFound = "No templates found in %s"
	ErrBlankOutputDir   = "Output directory not provided"
)

type WalkingExecutor struct {
	funcMap  map[string]interface{}
	tOptions []string
}

func New(funcMap map[string]interface{}, options []string) *WalkingExecutor {
	return &WalkingExecutor{
		funcMap:  funcMap,
		tOptions: options,
	}
}

func (we *WalkingExecutor) Execute(tmplDir, outputDir string, tmplData interface{}, owProvider provider.OverwriteProvider) error {
	var err error

	if skelputil.IsBlank(tmplDir) || !skelputil.PathExists(tmplDir) || skelputil.DirIsEmpty(tmplDir) {
		err = fmt.Errorf(ErrNoTemplatesFound, tmplDir)
	}

	if skelputil.IsBlank(outputDir) {
		err = fmt.Errorf(ErrBlankOutputDir, tmplDir)
	}

	if err == nil && !skelputil.PathExists(outputDir) {
		err = os.MkdirAll(outputDir, os.ModePerm)
	}

	if err == nil {
		err = filepath.Walk(tmplDir, func(curPath string, fi os.FileInfo, werr error) error {
			var terr error
			var relTarget string

			terr = werr

			if terr == nil {
				relTarget, terr = we.calculateRelativeTarget(tmplDir, curPath, tmplData)
			}

			if terr == nil {
				if fi.IsDir() {
					return skelputil.MkdirAll(filepath.Join(outputDir, relTarget))
				}

				terr = we.processFileTemplate(outputDir, relTarget, curPath, tmplData, owProvider)
			}

			return terr
		})
	}

	return err
}

func (we *WalkingExecutor) processFileTemplate(outputDir, relTarget, templatePath string, tmplData interface{}, owProvider provider.OverwriteProvider) error {
	var err error
	var fileTemplate *template.Template
	var destFile *os.File
	var srcMode os.FileMode

	absTarget := filepath.Join(outputDir, relTarget)

	// if the file exists, see if we should overwrite it
	if skelputil.PathExists(absTarget) && !owProvider(outputDir, relTarget) {
		return nil
	}

	fileTemplate, err = template.ParseFiles(templatePath)

	if err == nil {
		fileTemplate.Option(we.tOptions...).Funcs(we.funcMap)
		destFile, err = os.Create(absTarget)
	}

	if err == nil {
		err = fileTemplate.Execute(destFile, tmplData)
	}

	if err == nil {
		srcMode, err = skelputil.GetFileMode(templatePath)
	}

	if err == nil {
		os.Chmod(absTarget, srcMode)
		err = destFile.Close()
	}

	return err
}

func (we *WalkingExecutor) calculateRelativeTarget(tmplDir, curPath string, tmplData interface{}) (string, error) {
	var err error
	var relTmplFile string
	var target string
	var fnameTmpl *template.Template
	var b bytes.Buffer

	relTmplFile, err = filepath.Rel(tmplDir, curPath)

	if err == nil && !strings.Contains(relTmplFile, "{{") {
		return relTmplFile, nil
	}

	if err == nil {
		fnameTmpl, err = template.New("filename template").Option(we.tOptions...).Funcs(we.funcMap).Parse(relTmplFile)
	}

	if err == nil {
		err = fnameTmpl.Execute(&b, tmplData)
	}

	if err == nil {
		target = b.String()
	}

	return target, err
}
