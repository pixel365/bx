package internal

import (
	"errors"

	"github.com/charmbracelet/huh"

	"github.com/pixel365/bx/internal/model"
)

type Cfg string

var NoConfigError = errors.New("no config found in context")

const (
	CfgContextKey Cfg = "config"
)

type Printer interface {
	PrintSummary(verbose bool)
}

type OptionProvider interface {
	Option() string
}

type ConfigManager interface {
	Save() error
	Reset() error
	GetAccounts() []model.Account
	GetModules() []model.Module
	SetAccounts(...model.Account)
	SetModules(...model.Module)
	AddAccounts(...model.Account)
	AddModules(...model.Module)
}

func AccountIndexByLogin(accounts []model.Account, login string) (int, error) {
	if len(accounts) == 0 {
		return 0, errors.New("no accounts found")
	}

	for i, account := range accounts {
		if account.Login == login {
			return i, nil
		}
	}

	return 0, errors.New("account not found")
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

func Choose[T OptionProvider](items []T, value *string, title string) error {
	if len(items) == 0 {
		return errors.New("no items found")
	}

	var options []huh.Option[string]
	for _, item := range items {
		options = append(options, huh.NewOption(item.Option(), item.Option()))
	}

	if err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value).
		Run(); err != nil {
		return err
	}

	return nil
}
