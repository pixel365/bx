package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pixel365/bx/internal/helpers"
)

func Test_newVersionCommand(t *testing.T) {
	cmd := NewVersionCommand()
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

		if !cmd.HasFlags() {
			t.Error("flags is missing")
		}
	})
}

func TestNewVersionCommand(t *testing.T) {
	t.Run("short", func(t *testing.T) {
		cmd := NewVersionCommand()
		output := helpers.CaptureOutput(func() {
			_ = cmd.Execute()
		})

		want := fmt.Sprintf("%s\n", buildVersion)

		if output != want {
			t.Errorf("output = %v, want %v", output, want)
		}
	})
}

func TestNewVersionCommand_Verbose(t *testing.T) {
	t.Run("verbose", func(t *testing.T) {
		cmd := NewVersionCommand()
		cmd.SetArgs([]string{"--verbose"})
		output := helpers.CaptureOutput(func() {
			_ = cmd.Execute()
		})

		want := fmt.Sprintf("Version: %s\nCommit: %s\nDate: %s\nGo: %s %s/%s\n",
			buildVersion,
			buildCommit,
			buildDate,
			runtime.Version(),
			runtime.GOOS,
			runtime.GOARCH)

		if output != want {
			t.Errorf("output = %v, want %v", output, want)
		}
	})
}
