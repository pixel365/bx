package account

import (
	"context"

	"github.com/spf13/cobra"
)

func NewAccountCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account",
		Aliases: []string{"acc"},
		Short:   "Manage accounts",
	}

	cmd.AddCommand(addCmd())
	cmd.AddCommand(lsCmd())
	cmd.AddCommand(authCmd(ctx))
	cmd.AddCommand(rmCmd())
	cmd.AddCommand(moduleCmd())

	return cmd
}
