package run

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"

	"github.com/spf13/cobra"
)

func Test_newRunCommand(t *testing.T) {
	cmd := NewRunCommand()

	assert.NotNil(t, cmd)
	assert.NotNil(t, cmd.RunE)
	assert.Equal(t, "run", cmd.Use)
	assert.Empty(t, cmd.Aliases)
	assert.True(t, cmd.HasFlags())
	assert.False(t, cmd.HasSubCommands())
	assert.True(t, cmd.HasExample())
}

func Test_run_NoCommandSpecifiedError(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	require.NoError(t, err)

	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "NoCommandSpecifiedError"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewRunCommand()
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_run_IsValid(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	require.NoError(t, err)

	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "IsValid"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_run_HandleStages(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	require.NoError(t, err)

	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	originalHandleStages := handleStagesFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "HandleStages"
		}

		runCfg := map[string][]string{
			"testCommand": {
				"components",
			},
		}

		mod.Run = runCfg

		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	handleStagesFunc = func(ctx context.Context, stages []string, m *module.Module, logger interfaces.Logger,
		customCommandMode bool) error {
		return nil
	}
	defer func() {
		handleStagesFunc = originalHandleStages
	}()

	cmd := NewRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	err = cmd.Execute()
	require.NoError(t, err)
}

func Test_run_readModuleFromFlags_failed(t *testing.T) {
	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, errors.New("error")
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewRunCommand()
	cmd.SetArgs([]string{"--cmd", "testCommand"})
	err := cmd.Execute()
	require.Error(t, err)
}
