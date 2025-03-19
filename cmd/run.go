package cmd

import (
	"fmt"
	"sync"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func newRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a custom command",
		Example: `
# Run a custom command
bx run --name my_module --cmd custom_command
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("cmd", "c", "", "Command to run")

	return cmd
}

func run(cmd *cobra.Command, _ []string) error {
	command, err := cmd.Flags().GetString("cmd")
	if err != nil {
		return err
	}

	if command == "" {
		return internal.NoCommandSpecifiedError
	}

	module, err := internal.ReadModuleFromFlags(cmd)
	if err != nil {
		return err
	}

	if err := module.IsValid(); err != nil {
		return err
	}

	if err := internal.ValidateRun(module); err != nil {
		return err
	}

	stages, ok := module.Run[command]
	if !ok {
		return fmt.Errorf("command %s not found in module %s", command, module.Name)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(stages))

	err = internal.HandleStages(stages, module, &wg, errCh, nil, true)

	wg.Wait()
	close(errCh)

	return err
}
