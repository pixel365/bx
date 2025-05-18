package label

import (
	"errors"
	"net/http"
	"testing"

	"github.com/pixel365/bx/internal/client"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/module"
)

func TestNewLabelCommand(t *testing.T) {
	cmd := NewLabelCommand()

	t.Run("new label", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "label" {
			t.Errorf("cmd use = %v, want %v", cmd.Use, "label")
		}

		if cmd.Short != "Change module label" {
			t.Errorf("cmd short = %v, want %v", cmd.Short, "Change module label")
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

func TestLabelCommand_no_args(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{})

	t.Run("no args", func(t *testing.T) {
		err := cmd.Execute()
		if err == nil {
			t.Error("Execute() should return an error")
		}

		if err.Error() != "label is required" {
			t.Errorf("error = %v, want %v", err, "label is required")
		}
	})
}

func TestLabelCommand_invalid_label(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{"some label"})

	t.Run("invalid label", func(t *testing.T) {
		err := cmd.Execute()
		if err == nil {
			t.Error("Execute() should return an error")
		}

		if !errors.Is(err, errors2.ErrInvalidLabel) {
			t.Errorf("error = %v, want %v", err, errors2.ErrInvalidLabel)
		}
	})
}

func TestLabelCommand_valid_label(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{"stable"})

	originalReadModule := readModuleFromFlagsFunc
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return &module.Module{Version: "1.0.0"}, nil
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	authFunc = func(client client.HTTPClient, module *module.Module, password string, silent bool) ([]*http.Cookie, error) {
		return nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "", nil
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	t.Run("valid label", func(t *testing.T) {
		err := cmd.Execute()
		if err == nil {
			t.Error("Execute() should return an error")
		}
	})
}

func TestLabelCommand_valid_label2(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{"stable"})

	originalReadModule := readModuleFromFlagsFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, errors.New("module error")
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	t.Run("valid label", func(t *testing.T) {
		err := cmd.Execute()
		if err == nil {
			t.Error("Execute() should return an error")
		}
	})
}

func TestLabelCommand_valid_label3(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{"stable"})

	originalReadModule := readModuleFromFlagsFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return &module.Module{Version: "1.0.0"}, nil
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "", errors.New("password error")
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	t.Run("valid label", func(t *testing.T) {
		err := cmd.Execute()
		if err == nil {
			t.Error("Execute() should return an error")
		}
	})
}
