package cmd

import (
	"os/exec"
	"path/filepath"

	"github.com/brainicorn/skelp/generator"
	"github.com/spf13/cobra"
)

const (
	defaultCompletionDir = "/etc/bash_completion.d/"
)

var (
	completionDir string
	noSudo        bool
)

func newBashmeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bashme",
		Short: "Creates a bash completion file for skelp",
		Long: `Creates a bash completion file for skelp.
By default the completion file is written to /etc/bash_completion.d/ using sudo.`,
		RunE: executeBashme,
	}

	cmd.Flags().BoolVar(&noSudo, "no-sudo", false, "will try to write the completion file without using sudo")
	cmd.Flags().StringVar(&completionDir, "output", defaultCompletionDir, "path to the directory where the completion file will be written")

	return cmd
}

func executeBashme(cmd *cobra.Command, args []string) error {
	var completionFilePath string

	opts := generator.DefaultOptions()
	gen := generator.New(opts)

	skelpHome, err := gen.InitSkelpHome()

	if err == nil {
		completionFilePath = filepath.Join(skelpHome, "completion.sh")
		err = cmd.Parent().GenBashCompletionFile(completionFilePath)
	}

	if err == nil {
		cmd := "sudo"
		args := []string{"cp", completionFilePath, completionDir}

		if noSudo {
			cmd = "cp"
			args = []string{completionFilePath, completionDir}
		}

		err = exec.Command(cmd, args...).Run()
	}

	if err == nil {
		cmd.Println("bash completion successfully installed")
		cmd.Println("please close and reopen your terminal")
	}

	return err
}
