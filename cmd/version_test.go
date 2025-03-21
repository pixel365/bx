package cmd

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/pixel365/bx/internal"
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

func Test_printVersion(t *testing.T) {
	t.Run("version", func(t *testing.T) {
		output := internal.CaptureOutput(func() {
			printVersion()
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

func TestNewVersionCommand_Run(t *testing.T) {
	cmd := newVersionCommand()
	output := internal.CaptureOutput(func() {
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
}
