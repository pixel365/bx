package internal

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"

	"gopkg.in/yaml.v3"
)

type FileExistsAction string
type ChangelogType string
type ChangelogConditionType string

const (
	Replace        FileExistsAction = "replace"
	Skip           FileExistsAction = "skip"
	ReplaceIfNewer FileExistsAction = "replace_if_newer"

	Commit ChangelogType = "commit"
	Tag    ChangelogType = "tag"

	Include ChangelogConditionType = "include"
	Exclude ChangelogConditionType = "exclude"
)

type TypeValue[T1 any, T2 any] struct {
	Type  T1 `yaml:"type"`
	Value T2 `yaml:"value"`
}

type Changelog struct {
	From      TypeValue[ChangelogType, string]            `yaml:"from"`
	To        TypeValue[ChangelogType, string]            `yaml:"to"`
	Condition TypeValue[ChangelogConditionType, []string] `yaml:"condition,omitempty"`
}

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
	Changelog      Changelog         `yaml:"changelog,omitempty"`
	Name           string            `yaml:"name"`
	Version        string            `yaml:"version"`
	Account        string            `yaml:"account"`
	BuildDirectory string            `yaml:"buildDirectory,omitempty"`
	LogDirectory   string            `yaml:"logDirectory,omitempty"`
	Repository     string            `yaml:"repository,omitempty"`
	Stages         []Stage           `yaml:"stages"`
	Ignore         []string          `yaml:"ignore"`
	Callbacks      []Callback        `yaml:"callbacks,omitempty"`
}

// IsValid validates the fields of the Module struct.
//
// It checks the following conditions:
//
//  1. The `Name` field must not be an empty string and must not contain spaces.
//  2. The `Version` field must be a valid version, validated by the `ValidateVersion` function.
//  3. The `Account` field must not be empty.
//  4. If the `Variables` map is not nil, it checks that each key and value in the map is non-empty.
//  5. The `Stages` field must contain at least one stage.
//     Each stage must have a valid `Name`, `To` field, and an `ActionIfFileExists` field.
//     Additionally, each `From` path in a stage must be non-empty.
//  6. If the `Ignore` field is not empty, each rule must be non-empty.
//  7. The `NormalizeStages` function is called to ensure the validity of the stages after other checks.
//
// If any of these conditions are violated, the method returns an error with a detailed message.
// If all validations pass, it returns nil.
func (m *Module) IsValid() error {
	if m.Name == "" {
		return errors.New("module name is required")
	}

	if strings.Contains(m.Name, " ") {
		return errors.New("module name must not contain spaces")
	}

	if err := ValidateVersion(m.Version); err != nil {
		return err
	}

	if m.Account == "" {
		return errors.New("account is not valid")
	}

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

	if len(m.Callbacks) > 0 {
		for index, callback := range m.Callbacks {
			if err := callback.IsValid(); err != nil {
				return fmt.Errorf("callback [%d]: %w", index, err)
			}
		}
	}

	if m.Repository != "" {
		_, err := git.PlainOpen(m.Repository)
		if err != nil {
			return fmt.Errorf("repository [%s]: %w", m.Repository, err)
		}
	}

	if err := m.ValidateChangelog(); err != nil {
		return err
	}

	return nil
}

// ToYAML converts the Module struct to its YAML representation.
//
// It uses the `yaml.Marshal` function to serialize the `Module` struct into a YAML format.
// If the conversion is successful, it returns the resulting YAML as a byte slice.
// If an error occurs during marshaling, it returns the error.
//
// Returns:
// - []byte: The YAML representation of the Module struct.
// - error: Any error that occurred during the marshaling process.
func (m *Module) ToYAML() ([]byte, error) {
	return yaml.Marshal(m)
}

// NormalizeStages processes and normalizes the stages in the Module by replacing any variables
// within the stage fields (Name, To, From) with values from the Module's Variables map.
//
// The method iterates over each stage in the Module's Stages slice, and for each field (Name, To, From),
// it uses the `ReplaceVariables` function to replace any placeholders with corresponding variable values.
//
// If any error occurs while replacing variables or processing the stages, it returns the error.
// If no errors are encountered, it returns nil.
//
// Returns:
// - error: Any error that occurred during the variable replacement process. If successful, returns nil.
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

// ZipPath generates the absolute path for the ZIP file associated with the Module.
//
// The method constructs a path by combining the Module's BuildDirectory and Version fields,
// appending the ".zip" extension. It then checks if the path exists and is valid using the `CheckPath` function.
//
// If the path is valid, it returns the cleaned absolute path of the ZIP file. If any error occurs
// during path creation or validation, it returns an empty string along with the error.
//
// Returns:
// - string: The absolute path of the ZIP file.
// - error: Any error encountered during path creation or validation, otherwise nil.
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

// StageCallback returns the callback associated with the given stage.
// If no matching callback is found, an error is returned.
//
// stageName - the name of the stage to find the callback for.
//
// Returns:
// - Runnable - the found callback if it exists.
// - error - an error if the callback is not found.
func (m *Module) StageCallback(stageName string) (Runnable, error) {
	for _, callback := range m.Callbacks {
		if callback.Stage == stageName {
			return callback, nil
		}
	}

	return Callback{}, errors.New("stage callback not found")
}

// ValidateChangelog validates the changelog configuration of the module.
// It checks for the presence and correctness of required fields:
// - Ensures 'repository' is specified if 'from' or 'to' types are defined.
// - Validates that 'from' and 'to' types are either 'commit' or 'tag'.
// - Confirms 'from' and 'to' values are non-empty.
// - If 'condition' is specified, checks that:
//   - a Condition type is either 'include' or 'exclude'.
//   - Condition values are non-empty and valid regular expressions.
//
// Returns an error detailing the first encountered validation issue, or nil if the configuration is valid.
func (m *Module) ValidateChangelog() error {
	if (m.Changelog.From.Type != "" || m.Changelog.To.Type != "") && m.Repository == "" {
		return errors.New("repository is required")
	}

	if m.Repository != "" {
		if m.Changelog.From.Type != Commit && m.Changelog.From.Type != Tag {
			return fmt.Errorf("changelog from: type must be %s or %s", Commit, Tag)
		}

		if m.Changelog.To.Type != Commit && m.Changelog.To.Type != Tag {
			return fmt.Errorf("changelog to: type must be %s or %s", Commit, Tag)
		}

		if m.Changelog.From.Value == "" {
			return errors.New("changelog from: value is required")
		}

		if m.Changelog.To.Value == "" {
			return errors.New("changelog to: value is required")
		}

		if m.Changelog.Condition.Type != "" {
			if m.Changelog.Condition.Type != Include && m.Changelog.Condition.Type != Exclude {
				return fmt.Errorf(
					"changelog [%s] condition: type must be %s or %s",
					m.Name,
					Include,
					Exclude,
				)
			}

			if len(m.Changelog.Condition.Value) == 0 {
				return errors.New("condition value is required")
			}

			for i, condition := range m.Changelog.Condition.Value {
				if condition == "" {
					return fmt.Errorf("condition [%d]: value is required", i)
				}

				_, err := regexp.Compile(condition)
				if err != nil {
					return fmt.Errorf("invalid condition [%d]: %w", i, err)
				}
			}
		}
	}

	return nil
}
