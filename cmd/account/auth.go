package account

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/fatih/color"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"

	"github.com/spf13/cobra"
)

func authCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with an account",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			login, _ := c.Flags().GetString("login")
			login = strings.TrimSpace(login)
			if login == "" {
				if err = internal.Choose(&conf.Accounts, &login,
					"Select the account you want to log in with:"); err != nil {
					return err
				}
			}

			password := ""
			if err = huh.NewInput().
				Title("Enter password:").
				Prompt("> ").
				EchoMode(1).
				Value(&password).
				Validate(internal.ValidatePassword).
				Run(); err != nil {
				return err
			}

			index, err := internal.AccountIndexByLogin(&conf.Accounts, login)
			if err != nil {
				return err
			}

			if conf.Accounts[index].IsLoggedIn() {
				confirm := false
				if err = internal.Confirmation(&confirm,
					fmt.Sprintf("Are you sure you want to re-login to %s?", login)); err != nil {
					return err
				}

				if !confirm {
					return nil
				}
			}

			cookies, err := postForm(ctx, login, password)
			if err != nil {
				return err
			}

			var newCookies []http.Cookie
			for _, cookie := range cookies {
				newCookies = append(newCookies, *cookie)
			}

			conf.Accounts[index].Cookies = newCookies

			if err := conf.Save(); err != nil {
				return err
			}

			color.Green("Account %s successfully logged in!", login)

			return nil
		},
	}

	cmd.Flags().StringP("login", "l", "", "Login")

	return cmd
}

func postForm(ctx context.Context, login, password string) ([]*http.Cookie, error) {
	ttlCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	body := url.Values{
		"AUTH_FORM":     {"Y"},
		"TYPE":          {"AUTH"},
		"USER_LOGIN":    {login},
		"USER_PASSWORD": {password},
		"USER_REMEMBER": {"Y"},
	}

	encodedBody := []byte(body.Encode())
	req, err := http.NewRequestWithContext(ttlCtx, http.MethodPost,
		"https://partners.1c-bitrix.ru/personal/", bytes.NewReader(encodedBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err = resp.Body.Close(); err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	var cookies []*http.Cookie
	for _, c := range resp.Cookies() {
		if c.Name == "BITRIX_SM_LOGIN" {
			cookies = resp.Cookies()
			break
		}
	}

	if len(cookies) == 0 {
		return nil, errors.New("no cookies found")
	}

	return cookies, nil
}
