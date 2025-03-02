package internal

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestDefaultYAML(t *testing.T) {
	const def = `name: "test"
version: "1.0.0"
account: ""
buildDirectory: "./dist/test"
logDirectory: "./logs/test"
variables:
  structPath: "./examples/structure"
  install: "install"
  bitrix: "{structPath}/bitrix"
  local: "{structPath}/local"
stages:
  - name: "components"
    to: "{install}/components"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/components"
      - "{local}/components"
  - name: "templates"
    to: "{install}/templates"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/templates"
      - "{local}/templates"
  - name: "rootFiles"
    to: "."
    actionIfFileExists: "replace"
    from:
      - "{structPath}/simple-file.php"
  - name: "testFiles"
    to: "test"
    actionIfFileExists: "replace"
    from:
      - "{structPath}/simple-file.php"
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

func TestReplaceVariables(t *testing.T) {
	vars := map[string]string{
		"foo":        "bar",
		"var1":       "value1",
		"var-2":      "value2",
		"var_3":      "value3",
		"var--4":     "value4",
		"var-_5":     "value5",
		"var--__--6": "value6",
	}

	type args struct {
		variables map[string]string
		input     string
		depth     int
	}
	tests := []struct {
		name    string
		want    string
		args    args
		wantErr bool
	}{
		{"Single replacement", "some bar", args{vars, "some {foo}", 0}, false},
		{"Negative depth", "", args{vars, "some {foo}", -1}, true},
		{
			"Multiple same variables",
			"some bar bar bar",
			args{vars, "some {foo} {foo} {foo}", 0},
			false,
		},
		{
			"Multiple different variables",
			"some value1 value2 value3 value4 value5 value6",
			args{vars, "some {var1} {var-2} {var_3} {var--4} {var-_5} {var--__--6}", 0},
			false,
		},
		{"Recursion depth limit", "", args{vars, "some {var1}", 6}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReplaceVariables(tt.args.input, tt.args.variables, tt.args.depth)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReplaceVariables() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ReplaceVariables() got = %v, want %v", got, tt.want)
			}
		})
	}
}
