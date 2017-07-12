package cmd

import (
	"fmt"
	"io"

	"github.com/brainicorn/skelp/generator"
	"github.com/brainicorn/skelp/skelputil"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

// Flags that are to be added to commands.
var (
	quietFlag    bool
	noColorFlag  bool
	homedirFlag  string
	skelpdirFlag string
)

func NewSkelpCommand() *cobra.Command {

	skelpCmd := &cobra.Command{
		Use:   "skelp",
		Short: "A commandline tool for generating skeleton projects",
		Long: `skelp is a commandline tool for applying templates to a directory.

Skelp can be used to generate full project skeletons and/or apply templates to
an existing project.`,
		SilenceErrors:     true,
		SilenceUsage:      true,
		PersistentPreRunE: validateRootFlags,
	}

	skelpCmd.PersistentFlags().BoolVar(&quietFlag, "quiet", false, "run in 'quiet mode'")
	skelpCmd.PersistentFlags().BoolVar(&noColorFlag, "no-color", false, "turn off terminal colors")
	skelpCmd.PersistentFlags().StringVar(&homedirFlag, "homedir", "", "path to override user's home directory where skelp stores data")
	skelpCmd.PersistentFlags().StringVar(&skelpdirFlag, "skelpdir", "", "override name of skelp folder within the user's home directory")

	addCommandsToRoot(skelpCmd)

	return skelpCmd
}

func validateRootFlags(cmd *cobra.Command, args []string) error {

	if noColorFlag {
		ansi.DisableColors(true)
	}

	if !skelputil.IsBlank(homedirFlag) && !skelputil.PathExists(homedirFlag) {
		return newUserError(fmt.Sprintf("%s is not a valid path for --homedir flag", homedirFlag))
	}

	return nil
}

func addCommandsToRoot(cmd *cobra.Command) {
	cmd.AddCommand(newApplyCommand())
	cmd.AddCommand(newAliasCommand())
	cmd.AddCommand(newBashmeCommand())
}

// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(args []string, out io.Writer) int {

	var cmd *cobra.Command
	var err error
	exitcode := 0
	skelpCmd := NewSkelpCommand()
	skelpCmd.SetArgs(args)

	if out != nil {
		skelpCmd.SetOutput(out)
	}

	if cmd, err = skelpCmd.ExecuteC(); err != nil {
		exitcode = 1
		if isUserError(err) {
			cmd.Println(colorError(err.Error()))
			cmd.Println(cmd.UsageString())
		} else {
			cmd.Println(colorError(err.Error()))
		}
	}

	return exitcode
}

func getBaseOptions() generator.SkelpOptions {
	opts := generator.DefaultOptions()

	if !skelputil.IsBlank(homedirFlag) {
		opts.HomeDirOverride = homedirFlag
	}

	if !skelputil.IsBlank(skelpdirFlag) {
		opts.SkelpDirOverride = skelpdirFlag
	}

	return opts
}

func colorError(s string) string {
	return ansi.Color(s, "red+b")
}
