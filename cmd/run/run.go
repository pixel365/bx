package run

import (
	"fmt"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/module"

	"github.com/spf13/cobra"
)

var (
	readModuleFromFlagsFunc = module.ReadModuleFromFlags
	handleStagesFunc        = module.HandleStages
)

func NewRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a custom command",
		Example: `
# Run a custom command
bx run --name my_module --cmd custom_command
`,
		RunE: run,
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("cmd", "c", "", "Command to run")

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	mod, err := readModuleFromFlagsFunc(cmd)
	if err != nil {
		return err
	}

	command, _ := cmd.Flags().GetString("cmd")
	if command == "" {
		return errors.ErrNoCommandSpecified
	}

	if err = module.ValidateRun(mod); err != nil {
		return err
	}

	stages, ok := mod.Run[command]
	if !ok {
		return fmt.Errorf("command %s not found in mod %s", command, mod.Name)
	}

	err = handleStagesFunc(cmd.Context(), stages, mod, nil, true)

	return err
}
