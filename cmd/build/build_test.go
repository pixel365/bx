package build

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	t.Run("parameters", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "build" {
			t.Errorf("cmd use = %v, want %v", cmd.Use, "build")
		}

		if cmd.Short != "Build a module" {
			t.Errorf("cmd short = %v, want %v", cmd.Short, "Build a module")
		}

		if len(cmd.Aliases) != 1 {
			t.Errorf("len(cmd.Aliases) = %v, want %v", len(cmd.Aliases), 1)
		}

		if cmd.Aliases[0] != "b" {
			t.Errorf("cmd.Aliases[0] = %v, want %v", cmd.Aliases[0], "b")
		}

		if !cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() should be true")
		}

		if cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() should be false")
		}

		if cmd.RunE == nil {
			t.Errorf("cmd.RunE is nil")
		}
	})
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
	if err == nil {
		t.Errorf("err is nil")
	}
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

	builderFunc = func(m *module.Module, logger interfaces.BuildLogger) interfaces.Builder {
		return &FakeSuccessBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--last", "--description", "some description"})
	err = cmd.Execute()
	if err != nil {
		t.Error(err)
	}
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

	builderFunc = func(m *module.Module, logger interfaces.BuildLogger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
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

	builderFunc = func(m *module.Module, logger interfaces.BuildLogger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--version", " invalid module version "})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
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

	builderFunc = func(m *module.Module, logger interfaces.BuildLogger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--version", "1.0.0"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
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

	builderFunc = func(m *module.Module, logger interfaces.BuildLogger) interfaces.Builder {
		return &FakeFailBuilder{}
	}
	defer func() {
		builderFunc = originalBuilder
	}()

	cmd := NewBuildCommand()
	cmd.SetArgs([]string{"--repository", "."})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
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

	builderFunc = func(m *module.Module, logger interfaces.BuildLogger) interfaces.Builder {
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
	if err == nil {
		t.Errorf("err is nil")
	}
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
	if err == nil {
		t.Errorf("err is nil")
	}
}
