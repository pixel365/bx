package push

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pixel365/bx/internal/client"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
)

type mockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}
func (m *mockHttpClient) SetCookies(_ *url.URL, _ []*http.Cookie) {}

func Test_push_ReadModuleFromFlags(t *testing.T) {
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

	authFunc = func(client client.HTTPClient, module *module.Module, password string,
		silent bool) ([]*http.Cookie, error) {
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

	cmd := NewPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_invalid_Version(t *testing.T) {
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
			mod.Account = "test account"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	authFunc = func(client client.HTTPClient, module *module.Module, password string,
		silent bool) ([]*http.Cookie, error) {
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

	cmd := NewPushCommand()
	cmd.SetArgs([]string{"--version", "testingVersion"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_auth(t *testing.T) {
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
			mod.Account = "auth"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	authFunc = func(client client.HTTPClient, module *module.Module, password string,
		silent bool) ([]*http.Cookie, error) {
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

	cmd := NewPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_upload(t *testing.T) {
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

	authFunc = func(client client.HTTPClient, module *module.Module,
		password string, silent bool) ([]*http.Cookie, error) {
		return nil, nil
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	uploadFunc = func(ctx context.Context, client client.HTTPClient, module *module.Module,
		cookies []*http.Cookie, silent bool) error {
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
	ctx := context.Background()
	type args struct {
		client  client.HTTPClient
		module  *module.Module
		cookies []*http.Cookie
		silent  bool
	}

	var c []*http.Cookie
	c = append(c, &http.Cookie{})

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil client", args{nil, &module.Module{}, nil, false}, true},
		{"nil module", args{&mockHttpClient{}, nil, nil, false}, true},
		{"nil cookies", args{&mockHttpClient{}, &module.Module{}, nil, false}, true},
		{"not silent", args{&mockHttpClient{}, &module.Module{}, c, false}, false},
	}

	origSpinnerFunc := spinnerFunc
	spinnerFunc = func(_ string, _ func(context.Context) error) error {
		return nil
	}

	defer func() { spinnerFunc = origSpinnerFunc }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := upload(ctx, tt.args.client, tt.args.module, tt.args.cookies, tt.args.silent); (err != nil) != tt.wantErr {
				t.Errorf("upload() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_push_valid_Version(t *testing.T) {
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

	cmd := NewPushCommand()
	cmd.SetArgs([]string{"--version", "1.0.0"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_valid_label(t *testing.T) {
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
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		mod, err := module.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "another account"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "", nil
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	cmd := NewPushCommand()
	cmd.SetArgs([]string{"--label", "stable"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_invalid_label(t *testing.T) {
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

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "", nil
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	cmd := NewPushCommand()
	cmd.SetArgs([]string{"--label", "some label"})
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_invalid_module(t *testing.T) {
	originalReadModule := readModuleFromFlagsFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, errors.New("some error")
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	cmd := NewPushCommand()
	err := cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_push_invalid_input(t *testing.T) {
	originalReadModule := readModuleFromFlagsFunc
	originalInputPasswordFunc := inputPasswordFunc

	readModuleFromFlagsFunc = func(cmd *cobra.Command) (*module.Module, error) {
		return nil, nil
	}
	defer func() {
		readModuleFromFlagsFunc = originalReadModule
	}()

	inputPasswordFunc = func(cmd *cobra.Command, module *module.Module) (string, error) {
		return "", errors.New("some error")
	}
	defer func() {
		inputPasswordFunc = originalInputPasswordFunc
	}()

	cmd := NewPushCommand()
	err := cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
