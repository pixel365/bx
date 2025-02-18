package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func newBuildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"b"},
		Short:   "Build a module",
		Example: `
# Build a module by name
bx build --name my_module

# Build a module by file path
bx build -f config.yaml

# Override version
bx build --name my_module --version 1.2.3
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return build(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("version", "v", "", "Version of the module")

	return cmd
}

func build(cmd *cobra.Command, _ []string) error {
	path := cmd.Context().Value(internal.RootDir).(string)
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	file, err := cmd.Flags().GetString("file")
	file = strings.TrimSpace(file)
	if err != nil {
		return err
	}

	version, err := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)
	if err != nil {
		return err
	}

	isFile := len(file) > 0

	if !isFile && name == "" {
		err := internal.Choose(internal.AllModules(path), &name, "")
		if err != nil {
			return err
		}
	}

	if isFile {
		path = file
	}

	module, err := internal.ReadModule(path, name, isFile)
	if err != nil {
		return err
	}

	if version != "" {
		module.Version = version
	}

	if err := module.IsValid(); err != nil {
		return err
	}

	module.Ctx = cmd.Context()

	if err := module.Build(); err != nil {
		return err
	}

	fmt.Printf("Module %s successfully built. Version: %s\n", module.Name, module.Version)

	return nil
}
