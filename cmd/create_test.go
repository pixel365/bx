package cmd

import (
	"errors"
	"testing"

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
