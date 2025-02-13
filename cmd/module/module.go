package module

import "github.com/spf13/cobra"

func NewModuleCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "module",
		Aliases: []string{"mod"},
		Short:   "Manage Modules",
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(lsCmd())
	cmd.AddCommand(rmCmd())

	return cmd
}
