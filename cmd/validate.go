package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/brainicorn/skelp/skelplate"
	"github.com/brainicorn/skelp/skelputil"
	"github.com/spf13/cobra"
)

const (
	errValidateMissingArgs    = "descriptor path is required"
	errDescriptorFileNotFound = "descriptor not found: %s"
)

func newValidateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "validate [descriptor-path]",
		Short:   "Validates a skelp descriptor",
		Long:    `Validates a skelp descriptor`,
		PreRunE: validateValidateFlags,
		RunE:    executeValidate,
	}

	return cmd
}

func validateValidateFlags(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return newUserError(errValidateMissingArgs)
	}

	return nil
}

func executeValidate(cmd *cobra.Command, args []string) error {
	var err error
	var jsonPath string
	var descriptorBytes []byte

	jsonPath, err = filepath.Abs(args[0])

	if err == nil {
		if !skelputil.PathExists(jsonPath) {
			err = fmt.Errorf(errDescriptorFileNotFound, jsonPath)
		}
	}

	if err == nil {
		descriptorBytes, err = ioutil.ReadFile(jsonPath)
	}

	if err == nil {
		_, err = skelplate.ValidateDescriptor(descriptorBytes)
	}

	if err == nil {
		cmd.Println("Descriptor OK.")
	}

	return err
}
