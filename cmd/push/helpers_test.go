package push

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/request"
)

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

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
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

	cmd := NewPushCommand()
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

	originalReadModule := readModuleFromFlagsFunc
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "test account"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
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

	cmd := NewPushCommand()
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

	originalReadModule := readModuleFromFlagsFunc
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc

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

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
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

	cmd := NewPushCommand()
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

	originalReadModule := readModuleFromFlagsFunc
	originalAuthFunc := authFunc
	originalUploadFunc := uploadFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "upload"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
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

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "", nil
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	cmd := NewPushCommand()
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

	originalReadModule := readModuleFromFlagsFunc
	originalAuthFunc := authFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "Version"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	authFunc = func(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
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

	cmd := NewPushCommand()
	cmd.SetArgs([]string{"--version", "1.0.0"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
