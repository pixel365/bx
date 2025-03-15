package cmd

import (
	"testing"
)

func Test_newVersionCommand(t *testing.T) {
	cmd := newVersionCommand()
	t.Run("version", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "version" {
			t.Errorf("cmd use = %v, want %v", cmd.Use, "version")
		}

		if !cmd.HasAlias("v") {
			t.Error("v alias is missing")
		}

		if cmd.Short != "Print the version information" {
			t.Errorf("cmd short = %v, want %v", cmd.Short, "Print the version information")
		}
	})
}
