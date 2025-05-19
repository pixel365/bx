package auth

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/pixel365/bx/internal/client"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/interfaces"
	"github.com/pixel365/bx/internal/module"
)

func Test_Authenticate(t *testing.T) {
	type args struct {
		client   client.HTTPClient
		module   *module.Module
		password string
		silent   bool
	}
	tests := []struct {
		name    string
		args    args
		want    []*http.Cookie
		wantErr bool
	}{
		{"nil module", args{
			module:   nil,
			password: "",
			silent:   false,
		}, nil, true},
		{"empty password", args{
			module:   &module.Module{},
			password: "",
			silent:   false,
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Authenticate(
				tt.args.client,
				tt.args.module,
				tt.args.password,
				tt.args.silent,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_InputPassword(t *testing.T) {
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
			cmd := &cobra.Command{}
			cmd.Flags().String("password", tt.data, "")

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

			res, err := InputPassword(cmd, mod)
			if (err != nil) != tt.wantErr {
				t.Errorf("handlePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if res != tt.want {
				t.Errorf("handlePassword() = %v, want %v", res, tt.want)
			}
		})
	}
}
