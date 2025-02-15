package config

import (
	"fmt"
	"time"

	"github.com/pixel365/bx/internal/model"
)

type MockConfig struct {
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Accounts  []model.Account `json:"accounts,omitempty"`
	Modules   []model.Module  `json:"modules,omitempty"`
}

func (o *MockConfig) PrintSummary(verbose bool) {
	if verbose {
		fmt.Printf("Created At: %s\n", o.CreatedAt)
		fmt.Printf("Updated At: %s\n", o.UpdatedAt)
		fmt.Printf("Accounts: %d\n", len(o.Accounts))
		fmt.Printf("Modules: %d\n", len(o.Modules))
	} else {
		fmt.Printf("Created At: %s\n", o.CreatedAt)
		fmt.Printf("Updated At: %s\n", o.UpdatedAt)
	}
}

func (o *MockConfig) Save() error {
	return nil
}

func (o *MockConfig) Reset() error {
	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	o.Accounts = nil
	o.Modules = nil

	return o.Save()
}

func (o *MockConfig) GetAccounts() []model.Account {
	return o.Accounts
}

func (o *MockConfig) GetModules() []model.Module {
	return o.Modules
}

func (o *MockConfig) SetAccounts(accounts ...model.Account) {
	o.Accounts = accounts
}

func (o *MockConfig) SetModules(modules ...model.Module) {
	o.Modules = modules
}

func (o *MockConfig) AddAccounts(accounts ...model.Account) {
	o.Accounts = append(o.Accounts, accounts...)
}

func (o *MockConfig) AddModules(modules ...model.Module) {
	o.Modules = append(o.Modules, modules...)
}

func NewMockConfig() (*MockConfig, error) {
	now := time.Now().UTC()
	cfg := &MockConfig{
		CreatedAt: now,
		UpdatedAt: now,
	}

	return cfg, nil
}
