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
