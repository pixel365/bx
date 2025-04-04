package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func Test_newCreateCommand(t *testing.T) {
	cmd := newCreateCommand()

	t.Run("", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "create" {
			t.Errorf("cmd.Use should be 'create' but got '%s'", cmd.Use)
		}

		if len(cmd.Aliases) != 1 {
			t.Errorf("len(cmd.Aliases) should be 1 but got %d", len(cmd.Aliases))
		}

		if cmd.Aliases[0] != "c" {
			t.Errorf("cmd.Aliases[0] should be 'c' but got '%s'", cmd.Aliases[0])
		}

		if cmd.Short != "Create a new module" {
			t.Errorf("cmd.Short should be 'Create a new module' but got '%s'", cmd.Short)
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}
	})
}

func Test_create_nil(t *testing.T) {
	t.Run("nil command", func(t *testing.T) {
		err := create(nil, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}

		if !errors.Is(err, internal.NilCmdError) {
			t.Errorf("err = %v, want %v", err, internal.NilCmdError)
		}
	})
}

func Test_create_empty_name(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"--name", ""})

	t.Run("empty name", func(t *testing.T) {
		err := create(cmd, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}
	})
}

func Test_create_not_empty_name(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), internal.RootDir, "."))
	cmd.SetArgs([]string{"--name", "test-module"})

	t.Run("not empty name", func(t *testing.T) {
		err := create(cmd, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}
	})
}

func Test_create_success(t *testing.T) {
	moduleName := fmt.Sprintf("mod-%d", time.Now().UTC().Unix())
	defer func() {
		err := os.Remove(fmt.Sprintf("./%s.yaml", moduleName))
		if err != nil {
			return
		}
	}()

	cmd := newCreateCommand()
	cmd.SetContext(context.WithValue(context.Background(), internal.RootDir, "."))
	cmd.SetArgs([]string{"--name", moduleName})

	t.Run("success", func(t *testing.T) {
		err := cmd.Execute()
		if err != nil {
			t.Errorf("err is nil but got %v", err)
		}
	})
}

func Test_create_module_exists(t *testing.T) {
	moduleName := fmt.Sprintf("mod-%d", time.Now().UTC().Unix())
	fileName := fmt.Sprintf("%s.yaml", moduleName)
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

	cmd := newCreateCommand()
	cmd.SetContext(context.WithValue(context.Background(), internal.RootDir, "."))
	cmd.SetArgs([]string{"--name", moduleName})

	t.Run("module exists", func(t *testing.T) {
		err := cmd.Execute()
		if err == nil {
			t.Errorf("err is nil but got %v", err)
		}
	})
}
