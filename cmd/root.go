package cmd

import (
	"fmt"
	"os"

	"github.com/brainicorn/skelp"

	"github.com/spf13/cobra"
)

var skelper *skelp.Skelp

var rootCmd = &cobra.Command{
	Use:   "skelp",
	Short: "A commandline tool for generating skeleton projects",
	Long:  `skelp is a commandline tool for generating skeleton projects`,
}

func init() {
	rootCmd.AddCommand(applyCmd)
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	s, err := skelp.NewSkelp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initializing skelp: %v\n", err)
		os.Exit(1)
	}

	skelper = s

	if c, err := rootCmd.ExecuteC(); err != nil {
		if isUserError(err) {
			c.Println("")
			c.Println(c.UsageString())
		} else {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}

		os.Exit(-1)
	}
}
