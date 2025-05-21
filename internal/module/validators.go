package module

import (
	e "errors"
	"fmt"
	"strings"

	"github.com/pixel365/bx/internal/validators"

	"github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"
)

func ValidateStages(stages []types.Stage) error {
	if len(stages) == 0 {
		return errors.ErrInvalidStages
	}

	for index, stage := range stages {
		if stage.Name == "" {
			return fmt.Errorf("stages [%d]: name is required", index)
		}

		if stage.To == "" {
			return fmt.Errorf("stages [%d]: to is required", index)
		}

		if stage.ActionIfFileExists == "" {
			return fmt.Errorf("stages [%d]: actionIfFileExists is required", index)
		}

		for pathIndex, path := range stage.From {
			if path == "" {
				return fmt.Errorf("stages [%s]: path [%d] is required", stage.Name, pathIndex)
			}
		}

		if err := ValidateRules(stage.Filter, fmt.Sprintf("stage [%d] filter", index)); err != nil {
			return err
		}
	}

	return nil
}

func ValidateRules(rules []string, name string) error {
	if len(rules) > 0 {
		for index, rule := range rules {
			if rule == "" {
				return fmt.Errorf("%s [%d]: rule is required", name, index)
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

func ValidateRelease(steps []string, filter func(string) (types.Stage, error)) error {
	return validateStagesList(steps, "release", filter)
}

func ValidateLastVersion(steps []string, filter func(string) (types.Stage, error)) error {
	return validateStagesList(steps, "lastVersion", filter)
}

func ValidateRun(m *Module) error {
	if m.Run == nil {
		return nil
	}

	if len(m.Run) == 0 {
		return errors.ErrInvalidRun
	}

	for key, stages := range m.Run {
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("run [%s]: key is required", key)
		}

		if strings.Contains(key, " ") {
			return fmt.Errorf("run [%s]: key must not contain spaces", key)
		}

		if err := validateStagesList(stages, fmt.Sprintf("run: %s stages", key), m.FindStage); err != nil {
			return err
		}
	}

	return nil
}

func ValidateMainFields(m *Module) error {
	if m.Name == "" {
		return errors.ErrEmptyModuleName
	}

	if strings.Contains(m.Name, " ") {
		return errors.ErrNameContainsSpace
	}

	if err := validators.ValidateVersion(m.Version); err != nil {
		return err
	}

	switch m.Label {
	case "", types.Alpha, types.Beta, types.Stable:
	default:
		return errors.ErrInvalidLabel
	}

	if m.Account == "" {
		return errors.ErrEmptyAccountName
	}

	return nil
}

func ValidateLog(m *Module) error {
	if m.Log == nil {
		return nil
	}

	if m.Log.Dir == "" {
		return e.New("log dir is required")
	}

	if m.Log.MaxSize <= 0 {
		return e.New("log max size is required")
	}

	if m.Log.MaxBackups <= 0 {
		return e.New("log max backups is required")
	}

	if m.Log.MaxAge <= 0 {
		return e.New("log max age is required")
	}

	return nil
}

func validateStagesList(
	stages []string,
	name string,
	find func(string) (types.Stage, error),
) error {
	if len(stages) == 0 {
		return fmt.Errorf("%s is required", name)
	}

	collection := make(map[string]struct{})
	for index, stage := range stages {
		if stage == "" {
			return fmt.Errorf("%s [%d]: stage is required", name, index)
		}

		if _, exists := collection[stage]; exists {
			return fmt.Errorf("%s [%d]: duplicate stage [%s]", name, index, stage)
		}

		collection[stage] = struct{}{}

		_, err := find(stage)
		if err != nil {
			return fmt.Errorf("%s [%d]: %w", name, index, err)
		}
	}

	return nil
}
