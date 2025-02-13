package account

import "github.com/spf13/cobra"

func NewAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Aliases: []string{"acc"},
		Short:   "Manage accounts",
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(lsCmd())
	cmd.AddCommand(authCmd())
	cmd.AddCommand(rmCmd())

	return cmd
}
