package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateModuleName_NotExisting(t *testing.T) {
	t.Run("TestValidateModuleName_NotExisting", func(t *testing.T) {
		if err := ValidateModuleName("not_exists", "./"); err != nil {
			t.Error(err)
		}
	})
}

func TestValidateModuleName_Existing(t *testing.T) {
	t.Run("TestValidateModuleName_Existing", func(t *testing.T) {
		name := fmt.Sprintf("%s_%d", "testing", time.Now().Unix())
		filePath, err := filepath.Abs(fmt.Sprintf("%s/%s.yaml", ".", name))
		if err != nil {
			t.Error()
		}

		err = os.WriteFile(filePath, []byte(DefaultYAML()), 0600)
		if err != nil {
			t.Error(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Error(err)
			}
		}(filePath)

		err = ValidateModuleName(name, ".")
		if err == nil {
			t.Errorf("error expected")
		}
	})
}

func TestValidateVersion(t *testing.T) {
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
			if err := ValidateVersion(tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
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
			if err := ValidatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateArgument(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty argument", args{""}, false},
		{"invalid argument", args{"*"}, false},
		{"valid argument", args{"--name"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateArgument(tt.args.arg); got != tt.want {
				t.Errorf("ValidateArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}
