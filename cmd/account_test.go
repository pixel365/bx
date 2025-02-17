package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"
)

func TestAccountAdd(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountAdd", func(t *testing.T) {
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

func TestAccountLs(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountLs", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		rootCmd.SetArgs([]string{"account", "ls"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err = rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error("invalid result")
			}
		})

		if output != "test\n" {
			t.Error("invalid result")
		}
	})
}

func TestAccountLsVerbose(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountLsVerbose", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		rootCmd.SetArgs([]string{"account", "ls", "-v"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err = rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error("invalid result")
			}
		})

		now := time.Now().UTC()
		cmp := fmt.Sprintf("Login: test\nCreated At: %s\nUpdated At: %s\nLogged in: false\n",
			now.Format(time.RFC822), now.Format(time.RFC822))
		if output != cmp {
			t.Error(output)
		}
	})
}

func TestAccountRmSuccess(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountRmSuccess", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		rootCmd.SetArgs([]string{"account", "rm", "--login", "test", "--confirm"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err = rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error(err)
			}
		})

		if output != "Account test was deleted.\n" {
			t.Error("invalid result")
		}
	})
}

func TestAccountRmAccountNotFound(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountRmAccountNotFound", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		rootCmd.SetArgs([]string{"account", "rm", "--login", "abc", "--confirm"})
		rootCmd.SetOut(&buf)

		internal.CaptureOutput(func() {
			err = rootCmd.ExecuteContext(ctx)
			if err == nil {
				t.Error("invalid result")
			}

			if err.Error() != internal.NoAccountFound.Error() {
				t.Error(err)
			}
		})
	})
}

func TestAccountNoModules(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountNoModules", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		rootCmd.SetArgs([]string{"account", "modules", "--login", "test"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err = rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error(err)
			}
		})

		if output != fmt.Sprintf("%s\n", internal.NoModulesFound.Error()) {
			t.Error("invalid result")
		}
	})
}

func TestAccountNoModulesNoAccount(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestAccountNoModulesNoAccount", func(t *testing.T) {
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"account", "add", "--login", "test", "--skip-auth", "true"})
		err := rootCmd.ExecuteContext(ctx)
		if err != nil {
			t.Error(err)
		}

		var buf bytes.Buffer
		rootCmd.SetArgs([]string{"account", "modules", "--login", "abc"})
		rootCmd.SetOut(&buf)

		internal.CaptureOutput(func() {
			err = rootCmd.ExecuteContext(ctx)
			if err == nil {
				t.Error("invalid result")
			}

			if err.Error() != internal.NoAccountFound.Error() {
				t.Error(err)
			}
		})
	})
}
