package internal

import (
	"github.com/charmbracelet/huh"

	"github.com/pixel365/bx/internal/model"
)

type Cfg string

const (
	CfgContextKey Cfg = "config"

	Yes = "Yes"
	No  = "No"
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
		return 0, NoAccountsFound
	}

	for i, account := range accounts {
		if account.Login == login {
			return i, nil
		}
	}

	return 0, NoAccountFound
}

func Confirmation(flag *bool, title string) error {
	if err := huh.NewConfirm().
		Title(title).
		Affirmative(Yes).
		Negative(No).
		Value(flag).
		Run(); err != nil {
		return err
	}

	return nil
}

func Choose[T OptionProvider](items []T, value *string, title string) error {
	if len(items) == 0 {
		switch any(items).(type) {
		case []model.Account:
			return NoAccountsFound
		case []model.Module:
			return NoModulesFound
		default:
			return NoItemsFound
		}
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
