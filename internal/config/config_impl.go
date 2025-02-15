package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pixel365/bx/internal/model"
)

const (
	configDirName  = "/.bx"
	configFileName = "/.config.json"
)

func (o *Config) PrintSummary(verbose bool) {
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

func (o *Config) Save() error {
	o.mu.Lock()
	defer o.mu.Unlock()

	filePath, err := path()
	if err != nil {
		return err
	}

	o.UpdatedAt = time.Now().UTC()
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}

func (o *Config) Reset() error {
	o.mu.Lock()

	now := time.Now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	o.Accounts = nil
	o.Modules = nil

	o.mu.Unlock()

	return o.Save()
}

func (o *Config) GetAccounts() []model.Account {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.Accounts
}

func (o *Config) GetModules() []model.Module {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.Modules
}

func (o *Config) SetAccounts(accounts ...model.Account) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Accounts = accounts
}

func (o *Config) SetModules(modules ...model.Module) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Modules = modules
}

func (o *Config) AddAccounts(accounts ...model.Account) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Accounts = append(o.Accounts, accounts...)
}

func (o *Config) AddModules(modules ...model.Module) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.Modules = append(o.Modules, modules...)
}

func NewConfig() (*Config, error) {
	var err error

	filePath, err := path()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	cleanPath := filepath.Clean(absPath)

	cfg := &Config{}
	content, err := os.ReadFile(cleanPath)
	if err == nil {
		err = json.Unmarshal(content, cfg)
	} else {
		if os.IsNotExist(err) {
			now := time.Now().UTC()
			cfg.CreatedAt = now
			cfg.UpdatedAt = now
			err = cfg.Save()
		}
	}

	return cfg, err
}

func path() (string, error) {
	dir, err := dirPath()
	if err != nil {
		return "", err
	}

	return dir + configFileName, nil
}

func dirPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dirFullPath := dir + configDirName
	if _, err = os.Stat(dirFullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir(dirFullPath, 0750)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return dirFullPath, nil
}
