package module

import (
	"strings"

	"github.com/pixel365/bx/internal"

	"github.com/spf13/cobra"
)

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List available modules",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, ok := c.Context().Value(internal.CfgContextKey).(internal.ConfigManager)
			if !ok {
				return internal.NoConfigError
			}

			if len(conf.GetModules()) == 0 {
				internal.ResultMessage(internal.NoModulesFound.Error())
				return nil
			}

			verbose, _ := c.Flags().GetBool("verbose")
			all, _ := c.Flags().GetBool("all")
			login, _ := c.Flags().GetString("login")
			login = strings.TrimSpace(login)

			if all {
				login = ""
			} else {
				if login == "" {
					if err := internal.Choose(conf.GetAccounts(), &login,
						"Select the account whose modules you want to see:"); err != nil {
						return err
					}
				}
			}

			if login != "" {
				for _, mod := range conf.GetModules() {
					mod.PrintSummary(verbose)
				}
			} else {
				j := 0
				for _, mod := range conf.GetModules() {
					if mod.Login != login {
						continue
					}

					j++
					mod.PrintSummary(verbose)
				}

				if j == 0 {
					internal.ResultMessage(internal.NoModulesFound.Error())
				}
			}

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")
	cmd.Flags().BoolP("all", "a", false, "Show all available modules")
	cmd.Flags().BoolP("login", "l", false, "Login whose modules need to be shown")

	return cmd
}
