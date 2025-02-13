package account

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/config"
	"github.com/pixel365/bx/internal/model"
)

func rmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove an account",
		RunE: func(_ *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			if len(conf.Accounts) == 0 {
				fmt.Println("No accounts found")
				return nil
			}

			login := ""
			confirm := false

			var options []huh.Option[string]
			for _, a := range conf.Accounts {
				options = append(options, huh.NewOption(a.Login, a.Login))
			}

			if err = huh.NewSelect[string]().
				Title("Select the account you want to delete:").
				Options(options...).
				Value(&login).
				Run(); err != nil {
				return err
			}

			if err = huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete %s?", login)).
				Affirmative("Yes").
				Negative("No").
				Value(&confirm).
				Run(); err != nil {
				return err
			}

			if confirm {
				var accounts []model.Account
				for _, a := range conf.Accounts {
					if a.Login != login {
						accounts = append(accounts, a)
					}
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

	return cmd
}
