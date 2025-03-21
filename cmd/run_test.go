package cmd

import (
	"errors"
	"testing"

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
