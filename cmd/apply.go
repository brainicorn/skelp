package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply [git-url|file-path|alias]",
	Short: "Apply a template to the current directory",
	Long:  `Apply a template to the current directory`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return newUserError("apply command requires an argument")
		}

		fmt.Println("skelphome: ", skelper.SkelpHome)
		return nil
	},
}
