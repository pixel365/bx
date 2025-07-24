package build

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
	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"

	"github.com/spf13/cobra"
)

type FakeSuccessBuilder struct{}
type FakeFailBuilder struct{}

func (m *FakeSuccessBuilder) Build(ctx context.Context) error   { return nil }
func (m *FakeSuccessBuilder) Prepare() error                    { return nil }
func (m *FakeSuccessBuilder) Rollback() error                   { return nil }
func (m *FakeSuccessBuilder) Collect(ctx context.Context) error { return nil }
func (m *FakeSuccessBuilder) Cleanup()                          {}

func (m *FakeFailBuilder) Build(ctx context.Context) error   { return errors.New("build error") }
func (m *FakeFailBuilder) Prepare() error                    { return errors.New("prepare error") }
func (m *FakeFailBuilder) Rollback() error                   { return errors.New("rollback error") }
func (m *FakeFailBuilder) Collect(ctx context.Context) error { return errors.New("collect error") }
func (m *FakeFailBuilder) Cleanup()                          {}

func Test_newBuildCommand(t *testing.T) {
	cmd := NewBuildCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "build", cmd.Use)
	assert.Equal(t, "Build a module", cmd.Short)
	assert.Len(t, cmd.Aliases, 1)
	assert.Equal(t, "b", cmd.Aliases[0])
	assert.True(t, cmd.HasFlags())
	assert.False(t, cmd.HasSubCommands())
	assert.NotNil(t, cmd.RunE)
}

func Test_build_IsValid(t *testing.T) {
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
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewBuildCommand()
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_build_success(t *testing.T) {
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
	originalBuilder := builderFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		mod.Account = "build_success"
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	builderFunc = func(m *module.Module, logger interfaces.Logger) interfaces.Builder {
		return &FakeSuccessBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--last", "--description", "some description"})
	err = cmd.Execute()
	require.NoError(t, err)
}

func Test_build_fail(t *testing.T) {
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
	originalBuilder := builderFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		mod.Account = "build_fail"
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	builderFunc = func(m *module.Module, logger interfaces.Logger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_build_invalid_version(t *testing.T) {
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
	originalBuilder := builderFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		mod.Account = "build_invalid_version"
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	builderFunc = func(m *module.Module, logger interfaces.Logger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--version", " invalid module version "})
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_build_valid_version(t *testing.T) {
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
	originalBuilder := builderFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		mod.Account = "build_valid_version"
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	builderFunc = func(m *module.Module, logger interfaces.Logger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--version", "1.0.0"})
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_build_repository(t *testing.T) {
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
	originalBuilder := builderFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		mod.Account = "build_repository"
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	builderFunc = func(m *module.Module, logger interfaces.Logger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--repository", "."})
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_build_invalid_last(t *testing.T) {
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
	originalBuilder := builderFunc
	originalLastFunc := validateLastVersionFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		mod.Account = "build_invalid_last"
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	builderFunc = func(m *module.Module, logger interfaces.Logger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	validateLastVersionFunc = func(steps []string, filter func(string) (types.Stage, error)) error {
		return errors.New("invalid last version")
	}
	defer func() {
		validateLastVersionFunc = originalLastFunc
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--last", "."})
	err = cmd.Execute()
	require.Error(t, err)
}

func Test_build_read_module_error(t *testing.T) {
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
		return nil, errors.New("read module error")
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewBuildCommand()
	err = cmd.Execute()
	require.Error(t, err)
}
