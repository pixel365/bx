package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/pixel365/bx/internal/interfaces"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/request"

	"github.com/spf13/cobra"
)

func Test_newPushCommand(t *testing.T) {
	cmd := newPushCommand()

	t.Run("should create new push command", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "push" {
			t.Errorf("new push command should use 'push', got '%s'", cmd.Use)
		}

		if cmd.Short != "Push module to a Marketplace" {
			t.Errorf("cmd.Short should be 'Push module to a Marketplace' but got '%s'", cmd.Short)
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

func Test_handlePassword(t *testing.T) {
	mod := &module.Module{}
	tests := []struct {
		name    string
		data    string
		want    string
		wantErr bool
	}{
		{"success", "123456", "123456", false},
		{"short password", "12345", "", true},
		{"empty password", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := newPushCommand()
			cmd.SetArgs([]string{"--password", tt.data, "--name", "test"})
			_ = cmd.Flags().Set("password", tt.data)
			_ = cmd.Flags().Set("name", "test")

			if tt.data == "" {
				origInput := inputPasswordFunc
				defer func() {
					inputPasswordFunc = origInput
				}()

				inputPasswordFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
					return nil
				}
			}

			res, err := handlePassword(cmd, mod)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if res != tt.want {
				t.Errorf("handlePassword() = %v, want %v", res, tt.want)
			}
		})
	}
}

func Test_push_nil(t *testing.T) {
	t.Run("nil command", func(t *testing.T) {
		err := push(nil, []string{})
		if err == nil {
			t.Errorf("err is nil")
		}

		if !errors.Is(err, errors2.NilCmdError) {
			t.Errorf("err = %v, want %v", err, errors2.NilCmdError)
		}
	})
}

func Test_push_ReadModuleFromFlags(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
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

	originalReadModule := readModuleFromFlags
	originalAuthFunc := authFunc

	readModuleFromFlags = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "ReadModuleFromFlags"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	origInputPasswordFunc := inputPasswordFunc
	inputPasswordFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return errors.New("input error")
	}

	defer func() {
		inputPasswordFunc = origInputPasswordFunc
	}()

	cmd := newPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
func Test_push_invalid_Version(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
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

	originalReadModule := readModuleFromFlags
	originalAuthFunc := authFunc

	readModuleFromFlags = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "test account"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	origInputPasswordFunc := inputPasswordFunc
	inputPasswordFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return nil
	}

	defer func() {
		inputPasswordFunc = origInputPasswordFunc
	}()

	cmd := newPushCommand()
	cmd.SetArgs([]string{"--version", "testingVersion"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_auth(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
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

	originalReadModule := readModuleFromFlags
	originalAuthFunc := authFunc

	readModuleFromFlags = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "auth"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	origInputPasswordFunc := inputPasswordFunc
	inputPasswordFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return nil
	}

	defer func() {
		inputPasswordFunc = origInputPasswordFunc
	}()

	cmd := newPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_upload(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
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

	originalReadModule := readModuleFromFlags
	originalAuthFunc := authFunc
	originalUploadFunc := uploadFunc

	readModuleFromFlags = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "upload"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, nil
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	uploadFunc = func(client *request.Client, module *module.Module, cookies []*http.Cookie) error {
		return errors.New("upload error")
	}
	defer func() {
		uploadFunc = originalUploadFunc
	}()

	origInputPasswordFunc := inputPasswordFunc
	inputPasswordFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return nil
	}

	defer func() {
		inputPasswordFunc = origInputPasswordFunc
	}()

	cmd := newPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_upload(t *testing.T) {
	type args struct {
		client  *request.Client
		module  *module.Module
		cookies []*http.Cookie
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil client", args{nil, &module.Module{}, nil}, true},
		{"nil module", args{&request.Client{}, nil, nil}, true},
		{"nil cookies", args{&request.Client{}, &module.Module{}, nil}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := upload(tt.args.client, tt.args.module, tt.args.cookies); (err != nil) != tt.wantErr {
				t.Errorf("upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_auth(t *testing.T) {
	type args struct {
		module   *module.Module
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *request.Client
		want1   []*http.Cookie
		wantErr bool
	}{
		{"nil module", args{
			module:   nil,
			password: "",
		}, nil, nil, true},
		{"empty password", args{
			module:   &module.Module{},
			password: "",
		}, nil, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := auth(tt.args.module, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("auth() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_push_valid_Version(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
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

	originalReadModule := readModuleFromFlags
	originalAuthFunc := authFunc

	readModuleFromFlags = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "Version"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	origInputPasswordFunc := inputPasswordFunc
	inputPasswordFunc = func(_ interfaces.Prompter, _ *string, _ string, _ func(string) error) error {
		return nil
	}

	defer func() {
		inputPasswordFunc = origInputPasswordFunc
	}()

	cmd := newPushCommand()
	cmd.SetArgs([]string{"--version", "1.0.0"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
