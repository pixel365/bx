package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func Execute(ctx context.Context) error {
	cmd := rootCmd(ctx)
	return cmd.ExecuteContext(ctx)
}

func rootCmd(_ context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use: "bx",
	}

	return cmd
}
