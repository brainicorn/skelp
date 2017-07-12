package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/brainicorn/skelp/generator"
	"github.com/brainicorn/skelp/provider"
	"github.com/brainicorn/skelp/skelplate"
	"github.com/brainicorn/skelp/skelputil"
	"github.com/spf13/cobra"
)

const (
	currentDirectory    = "current directory"
	errApplyMissingArgs = "template url or path is required"
)

var (
	outputDir string
	dataFile  string
	offline   bool
	force     bool
)

func newApplyCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:     "apply [git-url|file-path|alias]",
		Short:   "Apply a template to the current directory",
		Long:    `Apply a template to the current directory`,
		PreRunE: validateApplyFlags,
		RunE:    executeApply,
	}

	applyCmd.Flags().StringVarP(&outputDir, "output", "o", currentDirectory, "path to the directory where the template should be applied")
	applyCmd.Flags().StringVarP(&dataFile, "data", "d", "", "path to a json data file for filling in template data")
	applyCmd.Flags().BoolVar(&offline, "offline", false, "turns off auto-downloading/updating of templates")
	applyCmd.Flags().BoolVarP(&force, "force", "f", false, "force overwriting of files without asking")

	return applyCmd
}

func validateApplyFlags(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError(errApplyMissingArgs)
	}

	if outputDir != currentDirectory && !skelputil.PathExists(outputDir) {
		return newUserError(fmt.Sprintf("%s is not a valid path for --output flag", outputDir))
	}

	if !skelputil.IsBlank(dataFile) && !skelputil.PathExists(dataFile) {
		return newUserError(fmt.Sprintf("%s is not a valid path for --data flag", dataFile))
	}

	return nil
}

func executeApply(cmd *cobra.Command, args []string) error {
	var err error
	var defData map[string]interface{}
	var rawData []byte

	opts := getBaseOptions()

	if outputDir != currentDirectory {
		opts.OutputDir = outputDir
	}

	if offline {
		opts.CheckForUpdates = false
		opts.Download = false
	}

	if force {
		opts.OverwriteProvider = provider.AlwaysOverwriteProvider
	}

	if !skelputil.IsBlank(dataFile) {
		rawData, err = ioutil.ReadFile(dataFile)

		if err == nil {
			err = json.Unmarshal(rawData, &defData)
		}
	}

	if err == nil {
		gen := generator.New(opts)
		dp := skelplate.NewDataProvider(defData)
		err = gen.Generate(args[0], dp.DataProviderFunc)
	}

	return err
}
