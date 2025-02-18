package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

func ValidateVersion(version string) error {
	version = strings.TrimSpace(version)
	if version == "" {
		return errors.New("module version is required")
	}

	return nil
}
