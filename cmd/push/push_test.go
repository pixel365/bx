package push

import (
	"testing"
)

func Test_newPushCommand(t *testing.T) {
	cmd := NewPushCommand()

	t.Run("should create new push command", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "push" {
			t.Errorf("new push command should use 'push', got '%s'", cmd.Use)
		}

		if cmd.Short != "Push module to a Marketplace" {
			t.Errorf("cmd.Short should be 'Push module to a Marketplace' but got '%s'", cmd.Short)
		}

		if cmd.RunE == nil {
			t.Error("cmd RunE should not be nil")
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}

		if cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() should be false")
		}

		if len(cmd.Aliases) > 0 {
			t.Errorf("len(cmd.Aliases) should be 0 but got %d", len(cmd.Aliases))
		}
	})
}
