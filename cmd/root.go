package cmd

import (
	"context"

	"github.com/pixel365/bx/internal"

	"github.com/pixel365/bx/cmd/config"

	"github.com/pixel365/bx/cmd/module"

	"github.com/pixel365/bx/cmd/account"

	"github.com/spf13/cobra"
)

var confirm bool

func NewRootCmd(ctx context.Context, conf internal.ConfigManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bx",
		Short: "Command-line tool for developers of 1C-Bitrix platform modules.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx = context.WithValue(ctx, internal.CfgContextKey, conf)
			cmd.SetContext(ctx)
		},
	}

	cmd.PersistentFlags().
		BoolVarP(&confirm, "confirm", "", false, "Automatically confirms all yes/no prompts")

	cmd.AddCommand(account.NewAccountCommand(ctx))
	cmd.AddCommand(module.NewModuleCommand())
	cmd.AddCommand(config.NewConfigCmd())

	return cmd
}
