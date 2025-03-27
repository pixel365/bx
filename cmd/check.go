package cmd

import (
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

var checkStagesFunc = internal.CheckStages

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check the configuration of a module",
		Example: `
# Check the configuration of a module by name
bx check --name my_module


# Check the configuration of a module by file path
bx check -f module-path/config.yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return check(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("repository", "r", "", "Path to a repository")

	return cmd
}

// check handles the logic of checking the configuration of a module based on the flags provided by the user.
// It retrieves the module name, file path, and validates the module configuration, including its stages.
// The function supports checking modules by name or by the specified YAML file.
//
// Parameters:
//   - cmd (*cobra.Command): The Cobra command that invoked the check function.
//   - args ([]string): A slice of arguments passed to the command (unused here).
//
// Returns:
//   - error: An error if the module configuration is invalid or any other error occurs.
func check(cmd *cobra.Command, _ []string) error {
	if cmd == nil {
		return internal.NilCmdError
	}

	module, err := readModuleFromFlags(cmd)
	if err != nil {
		return err
	}

	repository, _ := cmd.Flags().GetString("repository")

	if repository != "" {
		module.Repository = repository
	}

	if err := module.IsValid(); err != nil {
		return err
	}

	if err := checkStagesFunc(module); err != nil {
		return err
	}

	internal.ResultMessage("ok")

	return nil
}
