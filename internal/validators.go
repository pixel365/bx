package internal

import (
	"errors"
	"strings"

	"github.com/pixel365/bx/internal/config"
)

func ValidateAccountLogin(login string, conf *config.Config) error {
	value := strings.TrimSpace(login)
	if value == "" {
		return errors.New("login is empty")
	}

	if len(conf.Accounts) > 0 {
		for _, account := range conf.Accounts {
			if account.Login == value {
				return errors.New("an account with this login already exists")
			}
		}
	}

	return nil
}
