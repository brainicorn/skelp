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

var skelpCmd = &cobra.Command{
	Use:   "skelp",
	Short: "A commandline tool for generating skeleton projects",
	Long: `skelp is a commandline tool for applying templates to a directory.

Skelp can be used to generate full project skeletons and/or apply templates to
an existing project.`,
}

func init() {
	skelpCmd.PersistentFlags().BoolVar(&quiet, "quiet", false, "run in 'quiet mode'")

	addCommands()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if c, err := skelpCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println("")
			c.Println(c.UsageString())
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		os.Exit(-1)
	}
}

func addCommands() {
	skelpCmd.AddCommand(applyCmd)
}
