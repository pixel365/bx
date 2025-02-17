package internal

import "errors"

var (
	NoConfigError             = errors.New("no config found in context")
	NoAccountFoundError       = errors.New("account not found")
	NoAccountsFoundError      = errors.New("no accounts found")
	NoItemsFoundError         = errors.New("item not found")
	EmptyLoginError           = errors.New("login is empty")
	AccountAlreadyExistsError = errors.New("account already exists")
	EmptyPasswordError        = errors.New("password is empty")
	PasswordTooShortError     = errors.New("password is too short")
	NoModulesFoundError       = errors.New("no modules found")
)
