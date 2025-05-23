package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"
)

type FakePromptSuccessor struct{}
type FakePromptFailer struct{}

func (p *FakePromptSuccessor) Input(_ string, _ func(string) error) error { return nil }
func (p *FakePromptSuccessor) GetValue() string                           { return "" }
func (p *FakePromptFailer) Input(_ string, _ func(string) error) error {
	return errors.New("fail")
}
func (p *FakePromptFailer) GetValue() string { return "" }

func TestDefaultYAML(t *testing.T) {
	const def = `name: "test"
version: "1.0.0"
account: ""
buildDirectory: "./dist/test"

log:
  dir: "./logs"
  maxSize: 10
  maxBackups: 5
  maxAge: 30
  localTime: true
  compress: true

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

builds:
  release:
    - "components"
  lastVersion:
    - "components"

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
		_, err := GetModulesDir()
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
		{"invalid path", args{"../some-file.txt"}, true},
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

func TestCheckContextNil(t *testing.T) {
	t.Run("TestCheckContextNil", func(t *testing.T) {
		err := CheckContext(nil) //nolint:staticcheck // SA1012: passing nil context is intentional
		if err == nil {
			t.Error("CheckContext() did not return an error")
		}

		if !errors.Is(err, errors2.ErrNilContext) {
			t.Error("CheckContext() did not return an error")
		}
	})
}

func TestCaptureOutput(t *testing.T) {
	t.Run("TestCaptureOutput", func(t *testing.T) {
		output := CaptureOutput(func() {
			fmt.Println("ok")
		})

		if output != "ok\n" {
			t.Errorf("CaptureOutput() = %v, want %v", output, "ok\n")
		}

		output = CaptureOutput(func() {
			fmt.Println("ok string")
		})

		if output != "ok string\n" {
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

func Test_Cleanup(t *testing.T) {
	type args struct {
		resource io.Closer
	}
	tests := []struct {
		args args
		name string
	}{
		{args: args{nil}, name: "nil"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Cleanup(tt.args.resource, nil)
		})
	}
}

func Test_Cleanup_chan(t *testing.T) {
	ch := make(chan error)

	t.Run("TestCleanup_chan", func(t *testing.T) {
		fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
		filePath, _ := filepath.Abs("./" + fileName)
		filePath = filepath.Clean(filePath)

		file, err := os.Create(filePath)
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			err := os.Remove(file.Name())
			if err != nil {
				return
			}
		}()

		go func() {
			Cleanup(file, ch)
			close(ch)
		}()
	})

	for e := range ch {
		t.Errorf("TestCleanup_chan failed: %v", e)
	}
}

func TestChoose(t *testing.T) {
	empty := ""
	type args struct {
		items *[]string
		value *string
		title string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty items", args{&[]string{}, &empty, ""}, true},
		{"empty item", args{&[]string{""}, &empty, ""}, true},
		{"single item", args{&[]string{"option"}, &empty, ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Choose(tt.args.items, tt.args.value, tt.args.title); (err != nil) != tt.wantErr {
				t.Errorf("Choose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isValidPath(t *testing.T) {
	type args struct {
		filePath string
		basePath string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"invalid path", args{"../some-file.txt", "."}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidPath(tt.args.filePath, tt.args.basePath); got != tt.want {
				t.Errorf("IsValidPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkPaths(t *testing.T) {
	errCh := make(chan error)
	stage := types.Stage{
		From: []string{"."},
	}

	CheckPaths(stage, errCh)
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		t.Errorf("CheckPaths() = %v, want %v", errs, []error{})
	}
}

func TestUserInput_success(t *testing.T) {
	t.Run("TestUserInput_success", func(t *testing.T) {
		prompter := FakePromptSuccessor{}
		value := ""
		err := UserInput(&prompter, &value, "title", func(string) error { return nil })
		if err != nil {
			t.Errorf("UserInput() err = %v", err)
		}
	})
}

func TestUserInput_fail(t *testing.T) {
	t.Run("TestUserInput_fail", func(t *testing.T) {
		prompter := FakePromptFailer{}
		value := ""
		err := UserInput(
			&prompter,
			&value,
			"title",
			func(string) error { return errors.New("fake error") },
		)
		if err == nil {
			t.Errorf("UserInput() err = %v", err)
		}
	})
}

func TestSortSemanticVersions(t *testing.T) {
	m := make(map[string]struct{})

	m["8.9.2"] = struct{}{}
	m["8.8.0"] = struct{}{}
	m["8.20.5"] = struct{}{}
	m["8.19.1"] = struct{}{}
	m["8.9.0"] = struct{}{}
	m["0.0.1"] = struct{}{}
	m["8.20.4"] = struct{}{}
	m["8.2.8"] = struct{}{}
	m["10.0.0"] = struct{}{}
	m["8.2.9"] = struct{}{}
	m["8.2.9-beta"] = struct{}{}
	m["8.2.9-alpha"] = struct{}{}
	m["8.2.9-4d42446a"] = struct{}{}

	type args struct {
		versions iter.Seq[string]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"ok",
			args{versions: maps.Keys(m)},
			[]string{
				"0.0.1",
				"8.2.8",
				"8.2.9-4d42446a",
				"8.2.9-alpha",
				"8.2.9-beta",
				"8.2.9",
				"8.8.0",
				"8.9.0",
				"8.9.2",
				"8.19.1",
				"8.20.4",
				"8.20.5",
				"10.0.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortSemanticVersions(tt.args.versions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortSemanticVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizeVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"with prefix", args{"v0.0.1"}, "v0.0.1"},
		{"without prefix", args{"0.0.1"}, "v0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeVersion(tt.args.version); got != tt.want {
				t.Errorf("normalizeVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
