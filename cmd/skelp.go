package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Flags that are to be added to commands.
var (
	quiet bool
)

func newSkelpCommand() *cobra.Command {
	skelpCmd := &cobra.Command{
		Use:   "skelp",
		Short: "A commandline tool for generating skeleton projects",
		Long: `skelp is a commandline tool for applying templates to a directory.

Skelp can be used to generate full project skeletons and/or apply templates to
an existing project.`,
		SilenceErrors: true,
	}

	skelpCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "run in 'quiet mode'")

	addCommandsToRoot(skelpCmd)

	return skelpCmd
}

func addCommandsToRoot(cmd *cobra.Command) {
	cmd.AddCommand(newApplyCommand())
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	skelpCmd := newSkelpCommand()

	if c, err := skelpCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println("error: ", err.Error())
			c.Println(c.UsageString())
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		panic("-1")
	}
}
