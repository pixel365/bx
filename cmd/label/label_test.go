package label

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/client"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/module"
)

func TestNewLabelCommand(t *testing.T) {
	cmd := NewLabelCommand()

	assert.NotNil(t, cmd)
	assert.Equal(t, "label", cmd.Use)
	assert.Equal(t, "Change module label", cmd.Short)
	assert.NotNil(t, cmd.RunE)
	assert.True(t, cmd.HasFlags())
	assert.False(t, cmd.HasSubCommands())
	assert.Empty(t, cmd.Aliases)
}

func TestLabelCommand_no_args(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	require.Error(t, err)
	assert.Equal(t, "label is required", err.Error())
}

func TestLabelCommand_invalid_label(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{"some label"})
	err := cmd.Execute()

	require.Error(t, err)
	assert.ErrorIs(t, err, errors2.ErrInvalidLabel)
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

	authFunc = func(client client.HTTPClient, module *module.Module,
		password string, silent bool) ([]*http.Cookie, error) {
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

	err := cmd.Execute()
	require.Error(t, err)
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

	err := cmd.Execute()
	require.Error(t, err)
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

	err := cmd.Execute()
	require.Error(t, err)
}

func TestLabelCommand_change_labels(t *testing.T) {
	cmd := NewLabelCommand()
	cmd.SetArgs([]string{"stable"})

	originalReadModule := readModuleFromFlagsFunc
	originalInputPasswordFunc := inputPasswordFunc
	originalChangeLabelsFunc := changeLabelsFunc
	origNewClient := newClientFunc
	originalAuthFunc := authFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return &module.Module{Version: "1.0.0", Account: "login"}, nil
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "password123456", nil
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	changeLabelsFunc = func(client client.HTTPClient, module *module.Module,
		cookies []*http.Cookie, versions types.Versions) error {
		return errors.New("change labels error")
	}

	defer func() {
		changeLabelsFunc = originalChangeLabelsFunc
	}()

	newClientFunc = func(ttl time.Duration) client.HTTPClient {
		return &client.MockHttpClient{}
	}
	defer func() {
		newClientFunc = origNewClient
	}()

	authFunc = func(client client.HTTPClient, module *module.Module, password string,
		silent bool) ([]*http.Cookie, error) {
		return nil, nil
	}

	defer func() {
		authFunc = originalAuthFunc
	}()

	err := cmd.Execute()
	require.Error(t, err)
}
