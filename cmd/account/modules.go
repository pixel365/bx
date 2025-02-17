package account

import (
	"strings"

	"github.com/pixel365/bx/internal"

	"github.com/spf13/cobra"
)

func moduleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "modules",
		Short:   "List available modules",
		Aliases: []string{"mods"},
		RunE: func(c *cobra.Command, _ []string) error {
			conf, ok := c.Context().Value(internal.CfgContextKey).(internal.ConfigManager)
			if !ok {
				return internal.NoConfigError
			}

			login, _ := c.Flags().GetString("login")
			login = strings.TrimSpace(login)
			if login == "" {
				if err := internal.Choose(conf.GetAccounts(), &login,
					"Select the account whose modules you want to show:"); err != nil {
					return err
				}
			} else {
				_, err := internal.AccountIndexByLogin(conf.GetAccounts(), login)
				if err != nil {
					return err
				}
			}

			j := 0
			verbose, _ := c.Flags().GetBool("verbose")
			for _, module := range conf.GetModules() {
				if module.Login == login {
					j++
					module.PrintSummary(verbose)
				}
			}

			if j == 0 {
				internal.ResultMessage(internal.NoModulesFoundError.Error())
				return nil
			}

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")
	cmd.Flags().StringP("login", "l", "", "Login")

	return cmd
}
