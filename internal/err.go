package internal

import "errors"

var (
	NoConfigError        = errors.New("no config found in context")
	NoAccountFound       = errors.New("account not found")
	NoAccountsFound      = errors.New("no accounts found")
	NoItemsFound         = errors.New("item not found")
	EmptyLogin           = errors.New("login is empty")
	AccountAlreadyExists = errors.New("account already exists")
	EmptyPassword        = errors.New("password is empty")
	PasswordTooShort     = errors.New("password is too short")
	NoModulesFound       = errors.New("no modules found")
)
