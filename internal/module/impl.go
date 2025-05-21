package module

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"

	"gopkg.in/yaml.v3"

	"github.com/pixel365/bx/internal/callback"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/repo"
)

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
	if err := ValidateMainFields(m); err != nil {
		return err
	}

	if err := ValidateVariables(m); err != nil {
		return err
	}

	if err := ValidateStages(m.Stages); err != nil {
		return err
	}

	if err := ValidateRules(m.Ignore, "ignore"); err != nil {
		return err
	}

	if err := m.NormalizeStages(); err != nil {
		return err
	}

	if err := callback.ValidateCallbacks(m.Callbacks); err != nil {
		return err
	}

	if m.Repository != "" {
		if _, err := repo.OpenRepository(m.Repository); err != nil {
			return err
		}
	}

	if err := m.ValidateChangelog(); err != nil {
		return err
	}

	if err := ValidateRelease(m.Builds.Release, m.FindStage); err != nil {
		return err
	}

	if len(m.Builds.LastVersion) > 0 {
		if err := ValidateRelease(m.Builds.LastVersion, m.FindStage); err != nil {
			return err
		}
	}

	if err := ValidateRun(m); err != nil {
		return err
	}

	if err := ValidateLog(m); err != nil {
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
//   - []byte: The YAML representation of the Module struct.
//   - error: Any error that occurred during the marshaling process.
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
//   - error: Any error that occurred during the variable replacement process. If successful, returns nil.
func (m *Module) NormalizeStages() error {
	if m.Variables != nil {
		var err error
		for i, stage := range m.Stages {
			m.Stages[i].Name, err = helpers.ReplaceVariables(stage.Name, m.Variables, 0)
			if err != nil {
				return err
			}

			m.Stages[i].To, err = helpers.ReplaceVariables(stage.To, m.Variables, 0)
			if err != nil {
				return err
			}

			for j, from := range stage.From {
				m.Stages[i].From[j], err = helpers.ReplaceVariables(from, m.Variables, 0)
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
//   - string: The absolute path of the ZIP file.
//   - error: Any error encountered during path creation or validation, otherwise nil.
func (m *Module) ZipPath() (string, error) {
	path, _ := filepath.Abs(fmt.Sprintf("%s/%s.zip", m.BuildDirectory, m.Version))
	path = filepath.Clean(path)

	if err := helpers.CheckPath(path); err != nil {
		return path, err
	}

	return path, nil
}

// PasswordEnv returns the environment variable name
// that stores the password for the module.
//
// The variable name is generated based on the module's name:
//   - Converted to uppercase
//   - All dots (".") are replaced with underscores ("_")
//   - The suffix "_PASSWORD" is appended
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
//   - Runnable - the found callback if it exists.
//   - error - an error if the callback is not found.
func (m *Module) StageCallback(stageName string) (interfaces.Runnable, error) {
	for i := range m.Callbacks {
		cb := m.Callbacks[i]
		if cb.Stage == stageName {
			return cb, nil
		}
	}

	return nil, errors.ErrStageCallbackNotFound
}

// ValidateChangelog validates the changelog configuration of the module.
// It checks for the presence and correctness of required fields:
//   - Ensures 'repository' is specified if 'from' or 'to' types are defined.
//   - Validates that 'from' and 'to' types are either 'commit' or 'tag'.
//   - Confirms 'from' and 'to' values are non-empty.
//   - If 'condition' is specified, checks that:
//   - a Condition type is either 'include' or 'exclude'.
//   - Condition values are non-empty and valid regular expressions.
//
// Returns an error detailing the first encountered validation issue, or nil if the configuration is valid.
func (m *Module) ValidateChangelog() error {
	if m.Repository == "" || (m.Changelog.From.Type == "" && m.Changelog.To.Type == "") {
		return nil
	}

	if err := changeLogFromToValidate(m.Changelog); err != nil {
		return err
	}

	if m.Changelog.Condition.Type != "" {
		if m.Changelog.Condition.Type != types.Include &&
			m.Changelog.Condition.Type != types.Exclude {
			return fmt.Errorf(
				"changelog [%s] condition: type must be %s or %s",
				m.Name,
				types.Include,
				types.Exclude,
			)
		}

		if len(m.Changelog.Condition.Value) == 0 {
			return errors.ErrChangelogConditionValue
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

	switch m.Changelog.Sort {
	case "", types.Asc, types.Desc:
	default:
		return fmt.Errorf("changelog sort must be %s or %s", types.Asc, types.Desc)
	}

	return nil
}

// FindStage searches for a stage with the specified name in the module.
// It iterates through the module's stages and returns the matching stage if found.
// If no stage with the given name exists, it returns an empty Stage and an error
// with the message "stage not found".
//
// name - the name of the stage to search for.
//
// Returns:
//   - Stage: the stage with the matching name.
//   - error: nil if the stage is found; otherwise, an error indicating that the stage was not found.
func (m *Module) FindStage(name string) (types.Stage, error) {
	for _, stage := range m.Stages {
		if stage.Name == name {
			return stage, nil
		}
	}

	return types.Stage{}, fmt.Errorf("stage `%s` not found", name)
}

func changeLogFromToValidate(c types.Changelog) error {
	if c.From.Value == "" || c.To.Value == "" {
		return errors.ErrChangelogValue
	}

	if c.From.Type != types.Commit && c.From.Type != types.Tag {
		return fmt.Errorf("changelog from: type must be %s or %s", types.Commit, types.Tag)
	}

	if c.To.Type != types.Commit && c.To.Type != types.Tag {
		return fmt.Errorf("changelog to: type must be %s or %s", types.Commit, types.Tag)
	}

	return nil
}
