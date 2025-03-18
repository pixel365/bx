package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func newCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a new module",
		Example: `
# Create a new module
bx create --name my_module
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")

	return cmd
}

// create handles the logic of creating a new module by generating a YAML file for the module configuration.
// It validates the module name and uses default values to create the module's YAML configuration file.
//
// Parameters:
//   - cmd (*cobra.Command): The Cobra command that invoked the create function.
//   - args ([]string): A slice of arguments passed to the command (unused here).
//
// Returns:
//   - error: An error if the module name is invalid or any other error occurs during the creation process.
func create(cmd *cobra.Command, _ []string) error {
	name, err := cmd.Flags().GetString("name")
	name = strings.TrimSpace(name)
	if err != nil {
		return err
	}

	directory := cmd.Context().Value(internal.RootDir).(string)

	if name == "" {
		if err := huh.NewInput().
			Title("Enter Module Name:").
			Prompt("> ").
			Value(&name).
			Validate(func(input string) error {
				return internal.ValidateModuleName(input, directory)
			}).
			Run(); err != nil {
			return err
		}
	} else {
		if err := internal.ValidateModuleName(name, directory); err != nil {
			return err
		}
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s.yaml", directory, name))
	if err != nil {
		return err
	}

	var module internal.Module
	def := []byte(internal.DefaultYAML())
	if err = yaml.Unmarshal(def, &module); err != nil {
		return err
	}

	module.Name = name
	module.Account = ""

	out, err := module.ToYAML()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, out, 0600)
}
