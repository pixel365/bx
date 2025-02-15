package account

import (
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/fatih/color"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/model"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new account",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, ok := c.Context().Value(internal.CfgContextKey).(internal.ConfigManager)
			if !ok {
				return internal.NoConfigError
			}

			login, _ := c.Flags().GetString("login")
			login = strings.TrimSpace(login)
			if login == "" {
				if err := huh.NewInput().
					Title("Enter Login:").
					Prompt("> ").
					Value(&login).
					Validate(func(input string) error {
						return internal.ValidateAccountLogin(input, conf)
					}).
					Run(); err != nil {
					return err
				}
			}

			now := time.Now().UTC()
			account := model.Account{
				CreatedAt: now,
				UpdatedAt: now,
				Login:     login,
			}

			conf.AddAccounts(account)

			if err := conf.Save(); err != nil {
				return err
			}

			color.Green("Account created")

			confirm := false
			if err := internal.Confirmation(&confirm,
				"Do you want to log into this account right away?"); err != nil {
				return err
			}

			if confirm {
				c.Root().SetArgs([]string{c.Parent().Use, "auth", "--login", login})
				return c.Root().Execute()
			}

			return nil
		},
	}

	cmd.Flags().StringP("login", "l", "", "Login")

	return cmd
}
