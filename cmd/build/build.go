package build

import (
	"fmt"
	"strings"
	"time"

	"github.com/pixel365/bx/internal/module"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/logger"
	"github.com/pixel365/bx/internal/validators"

	"github.com/spf13/cobra"
)

var (
	builderFunc             = module.NewModuleBuilder
	validateLastVersionFunc = module.ValidateLastVersion
	readModuleFromFlagsFunc = module.ReadModuleFromFlags
)

func NewBuildCommand() *cobra.Command {
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
	cmd.Flags().StringP("description", "d", "", "Version description")
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
		return errors.NilCmdError
	}

	mod, err := readModuleFromFlagsFunc(cmd)
	if err != nil {
		return err
	}

	version, _ := cmd.Flags().GetString("version")

	if version != "" {
		version = strings.TrimSpace(version)
		if err := validators.ValidateVersion(version); err != nil {
			return err
		}
		mod.Version = version
	}

	repository, _ := cmd.Flags().GetString("repository")

	if repository != "" {
		mod.Repository = repository
	}

	description, _ := cmd.Flags().GetString("description")

	if description != "" {
		mod.Description = description
	}

	if err := mod.IsValid(); err != nil {
		return err
	}

	last, _ := cmd.Flags().GetBool("last")

	if last {
		if err := validateLastVersionFunc(mod.Builds.LastVersion, mod.FindStage); err != nil {
			return err
		}
	}

	mod.LastVersion = last

	logPath := fmt.Sprintf(
		"./%s-%s.%s.log",
		mod.Name,
		mod.GetVersion(),
		time.Now().UTC().Format(time.RFC3339),
	)
	loggerInstance := logger.NewFileZeroLogger(logPath, mod.LogDirectory)
	builder := builderFunc(mod, loggerInstance)
	defer builder.Cleanup()

	if err := builder.Build(); err != nil {
		return err
	}

	fmt.Printf("Module %s successfully built. Version: %s\n", mod.Name, mod.Version)

	return nil
}
