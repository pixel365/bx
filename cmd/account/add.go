package account

import (
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"
	"github.com/pixel365/bx/internal/model"
)

func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new account",
		RunE: func(_ *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			login := ""
			if err = huh.NewInput().
				Title("Enter Login:").
				Prompt("> ").
				Value(&login).
				Validate(func(input string) error {
					return internal.ValidateAccountLogin(input, conf)
				}).
				Run(); err != nil {
				return err
			}

			now := time.Now().UTC()
			account := model.Account{
				CreatedAt: now,
				UpdatedAt: now,
				Login:     login,
			}

			conf.Accounts = append(conf.Accounts, account)

			return conf.Save()
		},
	}

	return cmd
}
