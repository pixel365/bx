package cmd

import (
	"testing"
)

func Test_newCheckCommand(t *testing.T) {
	cmd := newCheckCommand()

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
