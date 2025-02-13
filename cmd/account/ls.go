package account

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/config"
)

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List accounts",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			if len(conf.Accounts) == 0 {
				fmt.Println("No accounts found")
				return nil
			}

			verbose, _ := c.Flags().GetBool("verbose")
			for _, acc := range conf.Accounts {
				acc.PrintSummary(verbose)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")

	return cmd
}
