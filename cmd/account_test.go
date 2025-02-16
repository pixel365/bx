package cmd

import (
	"context"
	"errors"
	"testing"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"
)

func TestAccountAdd(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		if len(cfg.GetAccounts()) != 1 {
			t.Error("no accounts added")
		}

		if cfg.GetAccounts()[0].Login != "test" {
			t.Error("invalid login")
		}

		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err = rootCmd.ExecuteContext(ctx)
		if err == nil {
			t.Error("invalid result")
		}

		if !errors.Is(err, internal.AccountAlreadyExists) {
			t.Errorf("invalid result: %v", err)
		}
	})
}
