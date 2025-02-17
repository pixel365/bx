package cmd

import (
	"context"
	"testing"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"
)

func TestNewRootCmd(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		ctx := context.Background()
		cfg, _ := config.NewMockConfig()
		var manager internal.ConfigManager = cfg
		root := NewRootCmd(ctx, manager)

		if root == nil {
			t.Error("root is nil")
		} else {
			if !root.HasSubCommands() {
				t.Error("subcommands is not set")
			}

			if root.Use != "bx" {
				t.Error("invalid use")
			}

			if root.Short != "Command-line tool for developers of 1C-Bitrix platform modules." {
				t.Error("invalid short")
			}

			if root.HasParent() {
				t.Error("parent is set")
			}

			if root.HasFlags() {
				t.Error("flags is set")
			}

			if root.Hidden {
				t.Error("hidden is set")
			}

			if !root.HasPersistentFlags() {
				t.Error("persistent flags is not set")
			}
		}
	})
}
