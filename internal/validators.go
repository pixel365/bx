package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var versionRegex = regexp.MustCompile(`(?m)^(\d+\.\d+\.\d+)$`)

// ValidateModuleName checks if a module with the given name already exists in the specified directory.
//
// The function constructs the expected file path for the module definition using the format "<directory>/<name>.yaml".
// It then checks if the file exists:
// - If the file does not exist, the function returns nil (indicating the module name is available).
// - If the file exists, it returns an error indicating that the module name is already taken.
// - If an error occurs while checking the file, it is returned.
//
// Returns nil if the module name is available, otherwise returns an error.
func ValidateModuleName(name, directory string) error {
	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s.yaml", directory, name))
	if err != nil {
		return err
	}

	if _, err := os.Stat(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
	}

	return fmt.Errorf("module name %s exists", name)
}

// ValidateVersion checks if the given module version is valid.
//
// The function performs the following checks:
// - Trims any leading and trailing whitespace.
// - Ensures the version is not empty.
// - Validates the version format using a predefined regex.
//
// Returns nil if the version is valid, otherwise returns an error.
func ValidateVersion(version string) error {
	version = strings.TrimSpace(version)
	if version == "" {
		return errors.New("module version is required")
	}

	for range versionRegex.FindAllString(version, -1) {
		return nil
	}

	return fmt.Errorf("invalid module version %s", version)
}

// ValidatePassword checks if the given password meets basic validation criteria.
//
// The function performs the following checks:
// - Trims any leading and trailing whitespace.
// - Ensures the password is not empty.
// - Ensures the password is at least 6 characters long.
//
// Returns an error if the password is invalid, otherwise returns nil.
func ValidatePassword(password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return errors.New("password is required")
	}

	if len(password) < 6 {
		return errors.New("password is too short")
	}

	return nil
}

// ValidateArgument checks if the provided argument contains only
// alphanumeric characters, underscores, slashes, or hyphens.
//
// Parameters:
// - arg (string): The argument to validate.
//
// Returns:
// - bool: True if the argument is valid, otherwise false.
func ValidateArgument(arg string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9>./_-]+$`)
	return re.MatchString(arg)
}

func ValidateStages(m *Module) error {
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

func ValidateIgnore(m *Module) error {
	if len(m.Ignore) > 0 {
		for index, rule := range m.Ignore {
			if rule == "" {
				return fmt.Errorf("ignore [%d]: rule is required", index)
			}
		}
	}

	return nil
}

func ValidateVariables(m *Module) error {
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

	return nil
}

func ValidateCallbacks(m *Module) error {
	if len(m.Callbacks) > 0 {
		for index, callback := range m.Callbacks {
			if err := callback.IsValid(); err != nil {
				return fmt.Errorf("callback [%d]: %w", index, err)
			}
		}
	}

	return nil
}

func ValidateBuilds(m *Module) error {
	if len(m.Builds.Release) == 0 {
		return errors.New("release is required")
	}

	for index, stage := range m.Builds.Release {
		if stage == "" {
			return fmt.Errorf("release [%d]: stage is required", index)
		}

		_, err := m.FindStage(stage)
		if err != nil {
			return fmt.Errorf("release [%d]: %w", index, err)
		}
	}

	if len(m.Builds.LastVersion) > 0 {
		for index, stage := range m.Builds.LastVersion {
			if stage == "" {
				return fmt.Errorf("lastVersion [%d]: stage is required", index)
			}

			_, err := m.FindStage(stage)
			if err != nil {
				return fmt.Errorf("lastVersion [%d]: %w", index, err)
			}
		}
	}

	return nil
}
