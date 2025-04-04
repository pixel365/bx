package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func Test_newRunCommand(t *testing.T) {
	cmd := newRunCommand()

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

func Test_run_nil(t *testing.T) {
	t.Run("nil command", func(t *testing.T) {
		err := run(nil, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}

		if !errors.Is(err, internal.NilCmdError) {
			t.Errorf("err = %v, want %v", err, internal.NilCmdError)
		}
	})
}

func Test_run_NoCommandSpecifiedError(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlags
	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "NoCommandSpecifiedError"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	cmd := newRunCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_run_IsValid(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlags
	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "IsValid"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	cmd := newRunCommand()
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

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlags
	originalHandleStages := handleStages
	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
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
		readModuleFromFlags = originalReadModule
	}()

	handleStages = func(stages []string, m *internal.Module, wg *sync.WaitGroup, errCh chan<- error,
		logger internal.BuildLogger, customCommandMode bool) error {
		return nil
	}
	defer func() {
		handleStages = originalHandleStages
	}()

	cmd := newRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
}
