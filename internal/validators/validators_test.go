package validators

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/helpers"
)

func TestValidateModuleName_NotExisting(t *testing.T) {
	err := ValidateModuleName("not_exists", "./")
	require.NoError(t, err)
}

func TestValidateModuleName_Existing(t *testing.T) {
	name := fmt.Sprintf("%s_%d", "testing", time.Now().Unix())
	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s.yaml", ".", name))
	require.NoError(t, err)

	err = os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
	if err != nil {
		t.Error(err)
	}
	defer func(name string) {
		err := os.Remove(name)
		require.NoError(t, err)
	}(filePath)

	err = ValidateModuleName(name, ".")
	require.Error(t, err)
}

func TestValidateVersion(t *testing.T) {
	t.Parallel()

	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{version: "1.0.0"}, wantErr: false},
		{name: "2", args: args{version: "v1.0.0"}, wantErr: true},
		{name: "3", args: args{version: "3.0.10"}, wantErr: false},
		{name: "4", args: args{version: ""}, wantErr: true},
		{name: "5", args: args{version: "some version"}, wantErr: true},
		{name: "6", args: args{version: "111.000.123"}, wantErr: false},
		{name: "7", args: args{version: "111.00x0.123"}, wantErr: true},
		{name: "8", args: args{version: "111.00x0.123"}, wantErr: true},
		{name: "9", args: args{version: "1x11.00x0.123"}, wantErr: true},
		{name: "10", args: args{version: "1x11.00x0.123x"}, wantErr: true},
		{name: "11", args: args{version: "1..1.1"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateVersion(tt.args.version)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{password: "123456"}, false},
		{"empty", args{password: ""}, true},
		{"only spaces", args{password: "    "}, true},
		{"short", args{password: "123"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidatePassword(tt.args.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestValidateArgument(t *testing.T) {
	t.Parallel()

	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"success", args{"arg"}, true},
		{"fail", args{"?arg"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := ValidateArgument(tt.args.arg)
			if tt.want {
				require.True(t, res)
			} else {
				require.False(t, res)
			}
		})
	}
}
