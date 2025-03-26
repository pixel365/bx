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

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
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
	cmd := newPushCommand()
	cmd.SetArgs([]string{"--password", "123456", "--name", "test"})
	_ = cmd.Flags().Set("password", "123456")
	_ = cmd.Flags().Set("name", "test")

	module := &internal.Module{}

	type args struct {
		cmd    *cobra.Command
		module *internal.Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"success", args{cmd, module}, "123456", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handlePassword(tt.args.cmd, tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("handlePassword() got = %v, want %v", got, tt.want)
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

		if !errors.Is(err, internal.NilCmdError) {
			t.Errorf("err = %v, want %v", err, internal.NilCmdError)
		}
	})
}

func Test_push_ReadModuleFromFlags(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
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

	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "ReadModuleFromFlags"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *internal.Module, password string) (*internal.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	cmd := newPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}
func Test_push_Version(t *testing.T) {
	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath := filepath.Join(fmt.Sprintf("./%s", fileName))
	filePath = filepath.Clean(filePath)

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
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

	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "Version"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *internal.Module, password string) (*internal.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
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

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
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

	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "auth"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *internal.Module, password string) (*internal.Client, []*http.Cookie, error) {
		return nil, nil, errors.New("auth error")
	}
	defer func() {
		authFunc = originalAuthFunc
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

	err := os.WriteFile(filePath, []byte(internal.DefaultYAML()), 0600)
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

	readModuleFromFlags = func(cmd *cobra.Command) (*internal.Module, error) {
		mod, err := internal.ReadModule(filePath, "", true)
		if err == nil {
			mod.Account = "upload"
		}
		return mod, err
	}
	defer func() {
		readModuleFromFlags = originalReadModule
	}()

	authFunc = func(module *internal.Module, password string) (*internal.Client, []*http.Cookie, error) {
		return nil, nil, nil
	}
	defer func() {
		authFunc = originalAuthFunc
	}()

	uploadFunc = func(client *internal.Client, module *internal.Module, cookies []*http.Cookie) error {
		return errors.New("upload error")
	}
	defer func() {
		uploadFunc = originalUploadFunc
	}()

	cmd := newPushCommand()
	err = cmd.Execute()
	if err == nil {
		t.Errorf("err is nil")
	}
}

func Test_upload(t *testing.T) {
	type args struct {
		client  *internal.Client
		module  *internal.Module
		cookies []*http.Cookie
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil client", args{nil, &internal.Module{}, nil}, true},
		{"nil module", args{&internal.Client{}, nil, nil}, true},
		{"nil cookies", args{&internal.Client{}, &internal.Module{}, nil}, true},
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
		module   *internal.Module
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *internal.Client
		want1   []*http.Cookie
		wantErr bool
	}{
		{"nil module", args{
			module:   nil,
			password: "",
		}, nil, nil, true},
		{"empty password", args{
			module:   &internal.Module{},
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
