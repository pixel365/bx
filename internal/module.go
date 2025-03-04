package internal

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileExistsAction string

const (
	Replace        FileExistsAction = "replace"
	Skip           FileExistsAction = "skip"
	ReplaceIfNewer FileExistsAction = "replace_if_newer"
)

type Stage struct {
	Name               string           `yaml:"name"`
	To                 string           `yaml:"to"`
	ActionIfFileExists FileExistsAction `yaml:"actionIfFileExists"`
	From               []string         `yaml:"from"`
	ConvertTo1251      bool             `yaml:"convertTo1251,omitempty"`
}

type Module struct {
	Ctx            context.Context   `yaml:"-"`
	Variables      map[string]string `yaml:"variables,omitempty"`
	Name           string            `yaml:"name"`
	Version        string            `yaml:"version"`
	Account        string            `yaml:"account"`
	Repository     string            `yaml:"repository,omitempty"`
	BuildDirectory string            `yaml:"buildDirectory,omitempty"`
	LogDirectory   string            `yaml:"logDirectory,omitempty"`
	Stages         []Stage           `yaml:"stages"`
	Ignore         []string          `yaml:"ignore"`
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

	if m.Variables != nil {
		i := 0
		for key, value := range m.Variables {
			i++
			if key == "" {
				return fmt.Errorf("variable [#%d]: key is required", i)
			}

			if value == "" {
				return fmt.Errorf("variable [%s]: value is required", key)
			}
		}
	}

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

	if len(m.Ignore) > 0 {
		for index, rule := range m.Ignore {
			if rule == "" {
				return fmt.Errorf("ignore [%d]: rule is required", index)
			}
		}
	}

	if err := m.NormalizeStages(); err != nil {
		return err
	}

	return nil
}

func (m *Module) ToYAML() ([]byte, error) {
	return yaml.Marshal(m)
}

func (m *Module) NormalizeStages() error {
	if m.Variables != nil {
		var err error
		for i, stage := range m.Stages {
			m.Stages[i].Name, err = ReplaceVariables(stage.Name, m.Variables, 0)
			if err != nil {
				return err
			}

			m.Stages[i].To, err = ReplaceVariables(stage.To, m.Variables, 0)
			if err != nil {
				return err
			}

			for j, from := range stage.From {
				m.Stages[i].From[j], err = ReplaceVariables(from, m.Variables, 0)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Module) ZipPath() (string, error) {
	path, err := filepath.Abs(fmt.Sprintf("%s/%s.zip", m.BuildDirectory, m.Version))
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)

	if err = CheckPath(path); err != nil {
		return path, err
	}

	return path, nil
}

// PasswordEnv returns the environment variable name
// that stores the password for the module.
//
// The variable name is generated based on the module's name:
// - Converted to uppercase
// - All dots (".") are replaced with underscores ("_")
// - The suffix "_PASSWORD" is appended
//
// For example, for a module named "my.module", the function will return "MY_MODULE_PASSWORD".
func (m *Module) PasswordEnv() string {
	name := strings.ToUpper(m.Name)
	name = strings.ReplaceAll(name, ".", "_")
	return fmt.Sprintf("%s_PASSWORD", name)
}
