package account

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/fatih/color"

	"github.com/charmbracelet/huh"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"

	"github.com/spf13/cobra"
)

func authCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with an account",
		RunE: func(_ *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			login := ""
			password := ""

			var options []huh.Option[string]
			for _, a := range conf.Accounts {
				options = append(options, huh.NewOption(a.Login, a.Login))
			}

			if err = huh.NewSelect[string]().
				Title("Select the account you want to log in with:").
				Options(options...).
				Value(&login).
				Run(); err != nil {
				return err
			}

			if err = huh.NewInput().
				Title("Enter password:").
				Prompt("> ").
				Value(&password).
				Validate(internal.ValidatePassword).
				Run(); err != nil {
				return err
			}

			index, err := internal.AccountIndexByLogin(&conf.Accounts, login)
			if err != nil {
				return err
			}

			if conf.Accounts[index].IsAuthenticated() {
				confirm := false
				if err = huh.NewConfirm().
					Title(fmt.Sprintf("Are you sure you want to re-login to %s?", login)).
					Affirmative("Yes").
					Negative("No").
					Value(&confirm).
					Run(); err != nil {
					return err
				}

				if !confirm {
					return nil
				}
			}

			cookies, err := postForm(login, password)
			if err != nil {
				return err
			}

			var c []http.Cookie
			for _, cookie := range cookies {
				c = append(c, *cookie)
			}

			conf.Accounts[index].Cookies = c

			if err := conf.Save(); err != nil {
				return err
			}

			color.Green("Account %s successfully logged in!", login)

			return nil
		},
	}

	return cmd
}

func postForm(login, password string) ([]*http.Cookie, error) {
	body := url.Values{
		"AUTH_FORM":     {"Y"},
		"TYPE":          {"AUTH"},
		"USER_LOGIN":    {login},
		"USER_PASSWORD": {password},
		"USER_REMEMBER": {"Y"},
	}

	r, err := http.PostForm("https://partners.1c-bitrix.ru/personal/", body)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			return
		}
	}()

	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}

	var cookies []*http.Cookie
	for _, c := range r.Cookies() {
		if c.Name == "BITRIX_SM_LOGIN" {
			cookies = r.Cookies()
			break
		}
	}

	if len(cookies) == 0 {
		return nil, errors.New("no cookies found")
	}

	return cookies, nil
}
