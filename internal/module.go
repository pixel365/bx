package internal

import (
	"context"
	"errors"
	"fmt"
)

type FileExistsAction string

const (
	Replace        FileExistsAction = "replace"
	Skip           FileExistsAction = "skip"
	ReplaceIfNewer FileExistsAction = "replace_if_newer"
)

type Item struct {
	Name               string           `yaml:"name"`
	To                 string           `yaml:"to"`
	ActionIfFileExists FileExistsAction `yaml:"actionIfFileExists"`
	From               []string         `yaml:"from"`
}

type Module struct {
	Ctx            context.Context
	Name           string   `yaml:"name"`
	Version        string   `yaml:"version"`
	Account        string   `yaml:"account"`
	Repository     string   `yaml:"repository,omitempty"`
	BuildDirectory string   `yaml:"buildDirectory,omitempty"`
	LogDirectory   string   `yaml:"logDirectory,omitempty"`
	Stages         []Item   `yaml:"stages"`
	Ignore         []string `yaml:"ignore"`
}

func (m *Module) IsValid() error {
	if m.Name == "" {
		return errors.New("module name is required")
	}

	if err := ValidateVersion(m.Version); err != nil {
		return err
	}

	if m.Account == "" {
		return errors.New("account is not valid")
	}

	//if m.Repository != "" {
	//	//TODO: check repository
	//}

	if len(m.Stages) == 0 {
		return errors.New("stages is not valid")
	}

	for index, item := range m.Stages {
		if item.Name == "" {
			return fmt.Errorf("stages [%d]: name is required", index)
		}

		if item.To == "" {
			return fmt.Errorf("stages [%d]: to is required", index)
		}

		if item.ActionIfFileExists == "" {
			return fmt.Errorf("stages [%d]: actionIfFileExists is required", index)
		}

		for pathIndex, path := range item.From {
			if path == "" {
				return fmt.Errorf("stages [%s]: path [%d] is required", item.Name, pathIndex)
			}
		}
	}

	return nil
}
