package cmd

import (
	"fmt"
	"strings"
	"time"

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

# Build .last_version
bx build --name my_module --last
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return build(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("version", "v", "", "Version of the module")
	cmd.Flags().StringP("repository", "r", "", "Path to a repository")
	cmd.Flags().BoolP("last", "", false, "Build a module .last_version.zip")

	return cmd
}

// build handles the logic of building a module based on the flags provided by the user.
// It retrieves the module name, file path, and version from the command flags, validates them,
// and triggers the build process for the module. The function supports building modules
// both by name and from a specified YAML file.
//
// Parameters:
//   - cmd (*cobra.Command): The Cobra command that invoked the build function.
//   - args ([]string): A slice of arguments passed to the command (unused here).
//
// Returns:
//   - error: An error if the build process encounters any issues or validation fails.
func build(cmd *cobra.Command, _ []string) error {
	if cmd == nil {
		return internal.NilCmdError
	}

	module, err := internal.ReadModuleFromFlags(cmd)
	if err != nil {
		return err
	}

	version, err := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)
	if err != nil {
		return err
	}

	if version != "" {
		if err := internal.ValidateVersion(version); err != nil {
			return err
		}
		module.Version = version
	}

	repository, err := cmd.Flags().GetString("repository")
	if err != nil {
		return err
	}

	if repository != "" {
		module.Repository = repository
	}

	if err := module.IsValid(); err != nil {
		return err
	}

	last, err := cmd.Flags().GetBool("last")
	if err != nil {
		return err
	}

	if last {
		if err := internal.ValidateLastVersion(module); err != nil {
			return err
		}
	}

	module.LastVersion = last

	logPath := fmt.Sprintf(
		"./%s-%s.%s.log",
		module.Name,
		module.GetVersion(),
		time.Now().UTC().Format(time.RFC3339),
	)
	logger := internal.NewFileZeroLogger(logPath, module.LogDirectory)
	builder := internal.NewModuleBuilder(module, logger)
	defer builder.Cleanup()

	if err := builder.Build(); err != nil {
		return err
	}

	fmt.Printf("Module %s successfully built. Version: %s\n", module.Name, module.Version)

	return nil
}
