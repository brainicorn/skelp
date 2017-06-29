package cmd

import (
	"github.com/spf13/cobra"
)

var (
	dataFile string
)

func newApplyCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply [git-url|file-path|alias]",
		Short: "Apply a template to the current directory",
		Long:  `Apply a template to the current directory`,
		RunE:  executeApply,
	}

	applyCmd.Flags().StringVarP(&dataFile, "data", "d", "", "path to a json data file for filling in template data")

	return applyCmd
}

func executeApply(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError("apply command requires an argument")
	}

	// TODO check data file and use a file data provider instead of a prompter

	//		gen := generator.New()
	//		gen.Generate(args[0])
	return nil
}
