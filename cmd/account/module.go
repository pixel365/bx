package account

import (
	"fmt"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"

	"github.com/spf13/cobra"
)

func moduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "module",
		Short:   "List available modules",
		Aliases: []string{"mod"},
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			login := ""
			if err = internal.ChooseAccount(&conf.Accounts, &login,
				"Select the account whose modules you want to show:"); err != nil {
				return err
			}

			j := 0
			verbose, _ := c.Flags().GetBool("verbose")
			for _, module := range conf.Modules {
				if module.Login == login {
					j++
					module.PrintSummary(verbose)
				}
			}

			if j == 0 {
				fmt.Println("No modules found")
				return nil
			}

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")

	return cmd
}
