package cmd

import (
	"context"

	"github.com/pixel365/bx/internal"

	"github.com/pixel365/bx/cmd/config"

	"github.com/pixel365/bx/cmd/module"

	"github.com/pixel365/bx/cmd/account"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context, conf internal.ConfigManager) error {
	cmd := rootCmd(ctx, conf)
	return cmd.ExecuteContext(ctx)
}

func rootCmd(ctx context.Context, conf internal.ConfigManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bx",
		Short: "Command-line tool for developers of 1C-Bitrix platform modules.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx = context.WithValue(ctx, internal.CfgContextKey, conf)
			cmd.SetContext(ctx)
		},
	}

	cmd.AddCommand(account.NewAccountCommand(ctx))
	cmd.AddCommand(module.NewModuleCommand())
	cmd.AddCommand(config.NewConfigCmd())

	return cmd
}
