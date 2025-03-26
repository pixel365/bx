package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func Test_newBuildCommand(t *testing.T) {
	cmd := newBuildCommand()

	t.Run("parameters", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "build" {
			t.Errorf("cmd use = %v, want %v", cmd.Use, "build")
		}

		if cmd.Short != "Build a module" {
			t.Errorf("cmd short = %v, want %v", cmd.Short, "Build a module")
		}

		if len(cmd.Aliases) != 1 {
			t.Errorf("len(cmd.Aliases) = %v, want %v", len(cmd.Aliases), 1)
		}

		if cmd.Aliases[0] != "b" {
			t.Errorf("cmd.Aliases[0] = %v, want %v", cmd.Aliases[0], "b")
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}

		if cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() should be false")
		}

		if cmd.RunE == nil {
			t.Errorf("cmd.RunE is nil")
		}
	})
}

func Test_build_nil(t *testing.T) {
	t.Run("nil command", func(t *testing.T) {
		err := build(nil, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}

		if !errors.Is(err, internal.NilCmdError) {
			t.Errorf("err = %v, want %v", err, internal.NilCmdError)
		}
	})
}

func Test_build_IsValid(t *testing.T) {
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
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	cmd := newBuildCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
