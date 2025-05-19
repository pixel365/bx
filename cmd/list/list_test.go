package list

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/client"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
)

func TestNewListCommand(t *testing.T) {
	cmd := NewListCommand()

	t.Run("new list command", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "list" {
			t.Errorf("new list command should use 'list', got '%s'", cmd.Use)
		}

		if cmd.Short != "List all module versions" {
			t.Errorf("cmd.Short should be 'List all module versions' but got '%s'", cmd.Short)
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

func Test_list_ReadModuleFromFlags(t *testing.T) {
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
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "ReadModuleFromFlags"
		}
		return mod, err
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

	cmd := NewListCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_list_auth(t *testing.T) {
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
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "some account"
		}
		return mod, err
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

	cmd := NewListCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_list_versions(t *testing.T) {
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
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc
	originalVersionsFunc := versionsFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "auth"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	authFunc = func(client client.HTTPClient, module *module.Module,
		password string, silent bool) ([]*http.Cookie, error) {
		return nil, nil
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

	defer func() {
		versionsFunc = originalVersionsFunc
	}()

	versionsFunc = func(ctx context.Context, client client.HTTPClient, module *module.Module,
		cookies []*http.Cookie) (types.Versions, error) {
		return nil, errors.New("versions error")
	}

	cmd := NewListCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
