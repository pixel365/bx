package cmd

import (
	"context"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd(context.Background())

	t.Run("", func(t *testing.T) {
		if cmd == nil {
			t.Error("nil cmd")
		} else {
			if cmd.Use != "bx" {
				t.Errorf("cmd.Use should be 'bx' but got '%s'", cmd.Use)
			}

			if cmd.Short != "Command-line tool for developers of 1C-Bitrix platform modules." {
				t.Errorf("invalid cmd.Short = '%s'", cmd.Short)
			}

			if cmd.HasParent() {
				t.Errorf("cmd.HasParent() = true")
			}

			if cmd.HasFlags() {
				t.Errorf("cmd.HasFlags() = true")
			}

			if !cmd.HasPersistentFlags() {
				t.Errorf("cmd.HasPersistentFlags() = false")
			}

			if !cmd.HasSubCommands() {
				t.Errorf("cmd.HasSubCommands() = false")
			}

			if cmd.Hidden {
				t.Errorf("cmd.Hidden = %v", cmd.Hidden)
			}
		}
	})
}
