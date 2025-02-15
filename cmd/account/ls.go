package account

import (
	"fmt"

	"github.com/pixel365/bx/internal"

	"github.com/spf13/cobra"
)

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List accounts",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, ok := c.Context().Value(internal.CfgContextKey).(internal.ConfigManager)
			if !ok {
				return internal.NoConfigError
			}

			if len(conf.GetAccounts()) == 0 {
				fmt.Println("No accounts found")
				return nil
			}

			verbose, _ := c.Flags().GetBool("verbose")
			for _, acc := range conf.GetAccounts() {
				acc.PrintSummary(verbose)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")

	return cmd
}
