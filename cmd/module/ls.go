package module

import (
	"fmt"
	"strings"

	"github.com/pixel365/bx/internal"

	"github.com/pixel365/bx/internal/config"

	"github.com/spf13/cobra"
)

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List available modules",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			if len(conf.Modules) == 0 {
				fmt.Println("No modules found")
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
					if err = internal.Choose(&conf.Accounts, &login,
						"Select the account whose modules you want to see:"); err != nil {
						return err
					}
				}
			}

			if login != "" {
				for _, mod := range conf.Modules {
					mod.PrintSummary(verbose)
				}
			} else {
				j := 0
				for _, mod := range conf.Modules {
					if mod.Login != login {
						continue
					}

					j++
					mod.PrintSummary(verbose)
				}

				if j == 0 {
					fmt.Println("No modules found")
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
