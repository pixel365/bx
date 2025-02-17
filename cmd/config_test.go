package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	config2 "github.com/pixel365/bx/cmd/config"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"
)

func TestConfig(t *testing.T) {
	t.Run("TestConfig", func(t *testing.T) {
		cmd := config2.NewConfigCmd()
		if cmd.Use != "config" {
			t.Errorf("cmd.Use = %s; want config", cmd.Use)
		}

		if cmd.Short != "Manage configuration" {
			t.Errorf("cmd.Short = %s; want Manage configuration", cmd.Short)
		}

		if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "conf" {
			t.Errorf("cmd.Aliases = %v; want conf", cmd.Aliases)
		}

		if !cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands = false; want true")
		}

		if cmd.HasFlags() {
			t.Errorf("cmd.HasFlags = true; want false")
		}
	})
}

func TestConfigResetSuccessfully(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestConfigResetSuccessfully", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"config", "reset", "-y"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err := rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error(err)
			}
		})

		if output != "Configuration file cleared\n" {
			t.Error(output)
		}
	})
}

func TestConfigInfo(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestConfigInfo", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"config", "info"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err := rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error(err)
			}
		})

		now := time.Now().UTC()
		cmp := fmt.Sprintf("Created At: %s\nUpdated At: %s\n",
			now.Format(time.RFC822), now.Format(time.RFC822))
		if output != cmp {
			t.Error(output)
		}
	})
}

func TestConfigInfoVerbose(t *testing.T) {
	ctx := context.Background()
	cfg, _ := config.NewMockConfig()
	var manager internal.ConfigManager = cfg

	t.Run("TestConfigInfoVerbose", func(t *testing.T) {
		var buf bytes.Buffer
		rootCmd := NewRootCmd(ctx, manager)
		rootCmd.SetArgs([]string{"config", "info", "-v"})
		rootCmd.SetOut(&buf)

		output := internal.CaptureOutput(func() {
			err := rootCmd.ExecuteContext(ctx)
			if err != nil {
				t.Error(err)
			}
		})

		now := time.Now().UTC()
		cmp := fmt.Sprintf("Created At: %s\nUpdated At: %s\nAccounts: 0\nModules: 0\n",
			now.Format(time.RFC822), now.Format(time.RFC822))
		if output != cmp {
			t.Error(output)
		}
	})
}
