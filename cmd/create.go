package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

	return os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
}
