package cmd

import (
	"fmt"

	"github.com/brainicorn/skelp/generator"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

const (
	errAliasAddMissingArgs    = "alias name and a template url/path are required"
	errAliasAddBadAlias       = "first argument must be a valid alias name"
	errAliasAddBadPath        = "second argument must be a valid template url or filepath"
	errAliasRemoveMissingName = "an alias name is required"
	errAliasRemoveInvalidName = "argument must be a valid alias name"
)

func newAliasCommand() *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "manage aliases for urls / filepaths",
		Long:  `manage aliases for urls / filepaths`,
	}

	aliasCmd.AddCommand(newAliasAddCommand())
	aliasCmd.AddCommand(newAliasListCommand())
	aliasCmd.AddCommand(newAliasRemoveCommand())
	return aliasCmd
}

func newAliasAddCommand() *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:     "add [alias name] [git-url|file-path]",
		Short:   "Create a short alias name for a template url/path",
		Long:    `Create a short alias name for a template url/path`,
		PreRunE: validateAliasAddFlags,
		RunE:    executeAliasAdd,
	}

	return aliasCmd
}

func newAliasListCommand() *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:   "list",
		Short: "list the registered aliases",
		Long:  `list the registered aliases`,
		RunE:  executeAliasList,
	}

	return aliasCmd
}

func newAliasRemoveCommand() *cobra.Command {
	aliasCmd := &cobra.Command{
		Use:     "remove [alias name]",
		Short:   "remove a registered alias",
		Long:    `remove a registered alias`,
		PreRunE: validateAliasRemoveFlags,
		RunE:    executeAliasRemove,
	}

	return aliasCmd
}

func validateAliasAddFlags(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return newUserError(errAliasAddMissingArgs)
	}

	if !generator.IsAlias(args[0]) {
		return newUserError(errAliasAddBadAlias)
	}

	if generator.IsAlias(args[1]) {
		return newUserError(errAliasAddBadPath)
	}

	return nil
}

func validateAliasRemoveFlags(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError(errAliasRemoveMissingName)
	}

	if !generator.IsAlias(args[0]) {
		return newUserError(errAliasRemoveInvalidName)
	}

	return nil
}

func executeAliasAdd(cmd *cobra.Command, args []string) error {

	gen := generator.New(getBaseOptions())

	return gen.AddAlias(args[0], args[1])
}

func executeAliasList(cmd *cobra.Command, args []string) error {
	gen := generator.New(getBaseOptions())

	entries, err := gen.AliasEntries()

	if err == nil {
		cmd.Println("------------------")
		cmd.Println("Registered Aliases")
		cmd.Println("------------------")

		for _, v := range entries {
			cmd.Println(fmt.Sprintf("%s -> %s", ansi.Color(v.Name, "green+b"), ansi.Color(v.Path, "blue+h")))
		}
	}

	return err
}

func executeAliasRemove(cmd *cobra.Command, args []string) error {

	gen := generator.New(getBaseOptions())

	return gen.RemoveAlias(args[0])
}
