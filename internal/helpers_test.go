package internal

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDefaultYAML(t *testing.T) {
	const def = `name: test
version: 1.0.0
account: test
repository: ""
buildDirectory: "./dist/test"
logDirectory: "./logs/test"
stages:
  - name: "components"
    to: "install/components"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/bitrix/components
      - ./examples/structure/local/components

  - name: "templates"
    to: "install/templates"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/bitrix/templates
      - ./examples/structure/local/templates

  - name: "rootFiles"
    to: "."
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/simple-file.php

  - name: "testFiles"
    to: "test"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/simple-file.php

ignore:
  - "**/*.log"
`
	t.Run("TestDefaultYAML", func(t *testing.T) {
		if DefaultYAML() != def {
			t.Error("Default YAML does not match")
		}
	})
}

func TestGetModulesDir(t *testing.T) {
	t.Run("TestGetModulesDir", func(t *testing.T) {
		_, err := GetModulesDir("")
		if err != nil {
			t.Errorf("GetModulesDir() returned an error: %s", err)
		}
	})
}

func TestCheckPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"invalid symbols", args{"asdfasdfasdf2#@$@"}, true},
		{"normal path", args{"./"}, false},
		{"not found", args{fmt.Sprintf("./test_404_%d", time.Now().UTC().Unix())}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckPath(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CheckPath() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsDir(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"current", args{"."}, true, false},
		{"404", args{fmt.Sprintf("./test_%d", time.Now().UTC().Unix())}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDir(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckContextOk(t *testing.T) {
	ctx := context.Background()
	t.Run("TestCheckContextOk", func(t *testing.T) {
		if err := CheckContext(ctx); err != nil {
			t.Errorf("CheckContext() returned an error: %s", err)
		}
	})
}

func TestCheckContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Run("TestCheckContextDone", func(t *testing.T) {
		cancel()
		if err := CheckContext(ctx); err == nil {
			t.Error("CheckContext() did not return an error")
		}
	})
}

func TestCaptureOutput(t *testing.T) {
	t.Run("TestCaptureOutput", func(t *testing.T) {
		output := CaptureOutput(func() {
			ResultMessage("ok")
		})

		if output != "ok\n" {
			t.Errorf("CaptureOutput() = %v, want %v", output, "ok\n")
		}
	})
}
