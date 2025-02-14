package internal

import (
	"errors"

	"github.com/charmbracelet/huh"

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

func ChooseAccount(accounts *[]model.Account, login *string, title string) error {
	var options []huh.Option[string]
	for _, a := range *accounts {
		options = append(options, huh.NewOption(a.Login, a.Login))
	}

	if err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(login).
		Run(); err != nil {
		return err
	}

	return nil
}

func Confirmation(flag *bool, title string) error {
	if err := huh.NewConfirm().
		Title(title).
		Affirmative("Yes").
		Negative("No").
		Value(flag).
		Run(); err != nil {
		return err
	}

	return nil
}
