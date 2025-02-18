package internal

import (
	"errors"
	"fmt"
)

type FileExistsMode string

const (
	Replace FileExistsMode = "replace"
	Skip    FileExistsMode = "skip"
	CopyNew FileExistsMode = "copy-new"
)

type Item struct {
	Name         string         `yaml:"name"`
	RelativePath string         `yaml:"relativePath"`
	IfFileExists FileExistsMode `yaml:"ifFileExists"`
	Paths        []string       `yaml:"paths"`
}

type Module struct {
	Name           string   `yaml:"name"`
	Version        string   `yaml:"version"`
	Account        string   `yaml:"account"`
	Repository     string   `yaml:"repository,omitempty"`
	BuildDirectory string   `yaml:"buildDirectory,omitempty"`
	LogDirectory   string   `yaml:"logDirectory,omitempty"`
	Mapping        []Item   `yaml:"mapping"`
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

	if len(m.Mapping) == 0 {
		return errors.New("mapping is not valid")
	}

	for index, item := range m.Mapping {
		if item.Name == "" {
			return fmt.Errorf("mapping [%d]: name is required", index)
		}

		if item.RelativePath == "" {
			return fmt.Errorf("mapping [%d]: relativePath is required", index)
		}

		if item.IfFileExists == "" {
			return fmt.Errorf("mapping [%d]: ifFileExists is required", index)
		}

		for pathIndex, path := range item.Paths {
			if path == "" {
				return fmt.Errorf("mapping [%s]: path [%d] is required", item.Name, pathIndex)
			}
		}
	}

	return nil
}
