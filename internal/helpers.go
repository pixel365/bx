package internal

import (
	"errors"

	"github.com/pixel365/bx/internal/model"
)

func AccountIndexByLogin(accounts *[]model.Account, login string) (int, error) {
	if len(*accounts) == 0 {
		return 0, errors.New("no accounts found")
	}

	for i, account := range *accounts {
		if account.Login == login {
			return i, nil
		}
	}

	return 0, errors.New("account not found")
}
