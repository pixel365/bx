package account

import (
	"fmt"
	"strings"

	"github.com/pixel365/bx/internal"

	"github.com/pixel365/bx/internal/config"

	"github.com/spf13/cobra"
)

func moduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "modules",
		Short:   "List available modules",
		Aliases: []string{"mods"},
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			login, _ := c.Flags().GetString("login")
			login = strings.TrimSpace(login)
			if login == "" {
				if err = internal.Choose(&conf.Accounts, &login,
					"OptionProvider the account whose modules you want to show:"); err != nil {
					return err
				}
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
	cmd.Flags().StringP("login", "l", "", "Login")

	return cmd
}
