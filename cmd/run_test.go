package cmd

import (
	"testing"
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
