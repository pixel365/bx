package account

import "github.com/spf13/cobra"

func NewAccountCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Aliases: []string{"acc"},
		Short:   "Manage accounts",
	}

	return cmd
}
