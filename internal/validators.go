package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)^(\d\.\d\.\d)$`)

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

	for range re.FindAllString(version, -1) {
		return nil
	}

	return fmt.Errorf("invalid module version %s", version)
}

func CheckContextActivity(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return nil
	}
}
