package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"

	"github.com/spf13/cobra"
)

func Test_newRunCommand(t *testing.T) {
	cmd := NewRunCommand()

	t.Run("subcommands", func(t *testing.T) {
		if cmd == nil {
			t.Errorf("cmd is nil")
		}

		if cmd.Use != "run" {
			t.Errorf("cmd use = %v, want %v", cmd.Use, "run")
		}

		if cmd.RunE == nil {
			t.Errorf("cmd.RunE is nil")
		}

		if len(cmd.Aliases) > 0 {
			t.Errorf("len(cmd.Aliases) should be 0 but got %d", len(cmd.Aliases))
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}

		if cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() should be false but got true")
		}

		if !cmd.HasExample() {
			t.Errorf("example is required")
		}
	})
}

func Test_run_NoCommandSpecifiedError(t *testing.T) {
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
			mod.Account = "NoCommandSpecifiedError"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewRunCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_run_IsValid(t *testing.T) {
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
			mod.Account = "IsValid"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_run_HandleStages(t *testing.T) {
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
	originalHandleStages := handleStagesFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "HandleStages"
		}

		runCfg := map[string][]string{
			"testCommand": {
				"components",
			},
		}

		mod.Run = runCfg

		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	handleStagesFunc = func(ctx context.Context, stages []string, m *module.Module, wg *sync.WaitGroup, errCh chan<- error,
		logger interfaces.BuildLogger, customCommandMode bool) error {
		return nil
	}
	defer func() {
		handleStagesFunc = originalHandleStages
	}()

	cmd := NewRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
}

func Test_run_readModuleFromFlags_failed(t *testing.T) {
	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, errors.New("error")
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	if err := cmd.Execute(); err == nil {
		t.Errorf("err is nil")
	}
}
