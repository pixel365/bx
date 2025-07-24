package check

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"

	"github.com/spf13/cobra"
)

var errFake = errors.New("fake error")

func Test_newCheckCommand(t *testing.T) {
	cmd := NewCheckCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "check", cmd.Use)
	assert.NotNil(t, cmd.RunE)
	assert.Empty(t, cmd.Aliases)
	assert.True(t, cmd.HasFlags())
}

func Test_check_ReadModuleFromFlags(t *testing.T) {
	originalReadModule := readModuleFromFlagsFunc
	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, errFake
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewCheckCommand()
	cmd.SetArgs([]string{})
	err := cmd.Execute()
	require.Error(t, err)
	assert.ErrorIs(t, err, errFake)
}

func Test_check_IsValid(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
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
			mod.Account = "check_test"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewCheckCommand()
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_check_repository(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
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
			mod.Account = "check_repository"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewCheckCommand()
	cmd.SetArgs([]string{"--repository", "."})
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_check_success(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Error(err)
		}
	}(filePath)

	originalReadModule := readModuleFromFlagsFunc
	originalCheckStages := checkStagesFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "test"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	checkStagesFunc = func(module *module.Module) error {
		return nil
	}
	defer func() {
		checkStagesFunc = originalCheckStages
	}()

	cmd := NewCheckCommand()
	err = cmd.Execute()
	require.NoError(t, err)
}
