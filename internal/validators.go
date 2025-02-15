package internal

import (
	"strings"
)

func ValidateAccountLogin(login string, conf ConfigManager) error {
	value := strings.TrimSpace(login)
	if value == "" {
		return EmptyLogin
	}

	if len(conf.GetAccounts()) > 0 {
		for _, account := range conf.GetAccounts() {
			if account.Login == value {
				return AccountAlreadyExists
			}
		}
	}

	return nil
}

func ValidatePassword(password string) error {
	value := strings.TrimSpace(password)
	if value == "" {
		return EmptyPassword
	}

	if len(value) < 6 {
		return PasswordTooShort
	}

	return nil
}
