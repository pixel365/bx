package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	errors2 "github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"

	"github.com/spf13/cobra"
)

var fakeError = errors.New("fake error")

func Test_newCheckCommand(t *testing.T) {
	cmd := NewCheckCommand()

	t.Run("parameters", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "check" {
			t.Errorf("cmd use = %v, want %v", cmd.Use, "check")
		}

		if cmd.RunE == nil {
			t.Errorf("cmd run is nil")
		}

		if len(cmd.Aliases) > 0 {
			t.Errorf("len(cmd.Aliases) should be 0 but got %d", len(cmd.Aliases))
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}

		if cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() should be false")
		}
	})
}

func Test_check_nil(t *testing.T) {
	t.Run("nil command", func(t *testing.T) {
		err := check(nil, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}

		if !errors.Is(err, errors2.NilCmdError) {
			t.Errorf("err = %v, want %v", err, errors2.NilCmdError)
		}
	})
}

func Test_check_ReadModuleFromFlags(t *testing.T) {
	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, fakeError
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewCheckCommand()
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}

	if !errors.Is(err, fakeError) {
		t.Errorf("err = %v, want %v", err, "fake error")
	}
}

func Test_check_IsValid(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "check_test"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewCheckCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_check_repository(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "check_repository"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewCheckCommand()
	cmd.SetArgs([]string{"--repository", "."})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_check_success(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	originalCheckStages := checkStagesFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "test"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	checkStagesFunc = func(module *module.Module) error {
		return nil
	}
	defer func() {
		checkStagesFunc = originalCheckStages
	}()

	cmd := NewCheckCommand()
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
}
