package account

import (
	"fmt"
	"strings"

	"github.com/pixel365/bx/internal"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/config"
	"github.com/pixel365/bx/internal/model"
)

func rmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove an account",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			login, _ := c.Flags().GetString("login")
			login = strings.TrimSpace(login)
			if login == "" {
				if err = internal.ChooseAccount(&conf.Accounts, &login,
					"Select the account you want to delete:"); err != nil {
					return err
				}
			}

			confirm := false
			if err = internal.Confirmation(&confirm,
				fmt.Sprintf("Are you sure you want to delete %s?", login)); err != nil {
				return err
			}

			if confirm {
				deleted := false
				var accounts []model.Account
				for _, a := range conf.Accounts {
					if a.Login == login {
						deleted = true
						continue
					}

					accounts = append(accounts, a)
				}

				if !deleted {
					return fmt.Errorf("account %s not found", login)
				}

				conf.Accounts = accounts

				if err := conf.Save(); err != nil {
					return err
				}

				fmt.Printf("Account %s was deleted.\n", login)
			}

			return nil
		},
	}

	cmd.Flags().StringP("login", "l", "", "Login")

	return cmd
}
