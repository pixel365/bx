package create

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

	"github.com/spf13/cobra"
)

func Test_newCreateCommand(t *testing.T) {
	cmd := NewCreateCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "create", cmd.Use)
	assert.Len(t, cmd.Aliases, 2)
	assert.Equal(t, []string{"c", "init"}, cmd.Aliases)
	assert.Equal(t, "Create a new module", cmd.Short)
	assert.True(t, cmd.HasFlags())
}

func Test_create_empty_name(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	cmd.SetArgs([]string{"--name", ""})

	origModuleInputNameFunc := moduleNameInputFunc

	moduleNameInputFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return errors.New("empty name")
	}

	defer func() {
		moduleNameInputFunc = origModuleInputNameFunc
	}()

	err := create(cmd, []string{})
	require.Error(t, err)
}

func Test_create_not_empty_name(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, "."))
	cmd.SetArgs([]string{"--name", "test-module"})

	origModuleInputNameFunc := moduleNameInputFunc

	moduleNameInputFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return nil
	}

	defer func() {
		moduleNameInputFunc = origModuleInputNameFunc
	}()

	err := create(cmd, []string{})
	require.NoError(t, err)
}

func Test_create_not_empty_invalid_name(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, "."))
	cmd.SetArgs([]string{"--name", "test-module"})

	origModuleInputNameFunc := moduleNameInputFunc

	moduleNameInputFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return errors.New("invalid name")
	}

	defer func() {
		moduleNameInputFunc = origModuleInputNameFunc
	}()

	err := create(cmd, []string{})
	require.Error(t, err)
}

func Test_create_success(t *testing.T) {
	moduleName := fmt.Sprintf("mod-%d", time.Now().UTC().Unix())
	defer func() {
		err := os.Remove(fmt.Sprintf("./%s.yaml", moduleName))
		if err != nil {
			return
		}
	}()

	cmd := NewCreateCommand()
	cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, "."))
	cmd.SetArgs([]string{"--name", moduleName})

	err := cmd.Execute()
	require.NoError(t, err)
}

func Test_create_module_exists(t *testing.T) {
	moduleName := fmt.Sprintf("mod-%d", time.Now().UTC().Unix())
	fileName := fmt.Sprintf("%s.yaml", moduleName)
	filePath := fmt.Sprintf("./%s", fileName)
	filePath = filepath.Clean(filePath)

	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			t.Error(err)
		}
	}()

	err := os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}

	cmd := NewCreateCommand()
	cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, "."))
	cmd.SetArgs([]string{"--name", moduleName})

	err = cmd.Execute()
	require.Error(t, err)
}
