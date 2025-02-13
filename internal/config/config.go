package config

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/pixel365/bx/internal/model"
)

const (
	configDirName  = "/.bx"
	configFileName = "/.config.json"
)

type Config struct {
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	Accounts  []model.Account `json:"accounts,omitempty"`
	Modules   []model.Module  `json:"modules,omitempty"`
	mu        sync.RWMutex
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

	return os.WriteFile(filePath, data, 0644)
}

func GetConfig() (*Config, error) {
	var err error

	filePath, err := path()
	if err != nil {
		return nil, err
	}

	config := &Config{}
	content, err := os.ReadFile(filePath)
	if err == nil {
		err = json.Unmarshal(content, config)
	} else {
		if os.IsNotExist(err) {
			now := time.Now().UTC()
			config.CreatedAt = now
			config.UpdatedAt = now
			err = config.Save()
		}
	}

	return config, err
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
			err = os.Mkdir(dirFullPath, os.ModePerm)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return dirFullPath, nil
}
