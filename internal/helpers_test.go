package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

type FakeBuildLogger struct{}

func (l *FakeBuildLogger) Info(message string, args ...interface{}) {}

func (l *FakeBuildLogger) Error(message string, err error, args ...interface{}) {}

func (l *FakeBuildLogger) Cleanup() {}

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
    convertTo1251: false

builds:
  release:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"
  lastVersion:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"

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
		if err := CheckContext(context.TODO()); err == nil {
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

		output = CaptureOutput(func() {
			ResultMessage("%s\n", "ok string")
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

func TestReadModuleFromFlags(t *testing.T) {
	t.Run("TestReadModuleFromFlags", func(t *testing.T) {
		_, err := ReadModuleFromFlags(nil)
		if err == nil {
			t.Errorf("ReadModuleFromFlags() did not return an error")
		}

		if !errors.Is(err, NilCmdError) {
			t.Errorf("err = %v, want %v", err, NilCmdError)
		}
	})
}

func TestReadModuleFromFlags_Name(t *testing.T) {
	t.Run("TestReadModuleFromFlags_Name", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.SetContext(context.WithValue(context.Background(), RootDir, RootDir))
		_, err := ReadModuleFromFlags(cmd)
		if err == nil {
			t.Errorf("ReadModuleFromFlags() did not return an error")
		}
	})
}

func TestReadModuleFromFlags_File(t *testing.T) {
	t.Run("TestReadModuleFromFlags_File", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.SetContext(context.WithValue(context.Background(), RootDir, RootDir))
		cmd.SetArgs([]string{"--file", "./test_files/foo"})
		_, err := ReadModuleFromFlags(cmd)
		if err == nil {
			t.Errorf("ReadModuleFromFlags() did not return an error")
		}
	})
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Choose(tt.args.items, tt.args.value, tt.args.title); (err != nil) != tt.wantErr {
				t.Errorf("Choose() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAllModules(t *testing.T) {
	name := fmt.Sprintf("%s_%d", "testing", time.Now().Unix())
	filePath, err := filepath.Abs(fmt.Sprintf("./%s/%s.yaml", ".", name))
	if err != nil {
		t.Error()
	}
	filePath = filepath.Clean(filePath)

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

	type args struct {
		directory string
	}
	tests := []struct {
		want *[]string
		name string
		args args
	}{
		{want: &[]string{"test"}, name: ".", args: args{directory: "."}},
		{want: nil, name: "fake dir", args: args{directory: "some/fake/dir"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AllModules(tt.args.directory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AllModules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckStages(t *testing.T) {
	err := CheckStages(nil)
	if !errors.Is(err, NilModuleError) {
		t.Errorf("CheckStages() error = %v, want %v", err, NilModuleError)
	}
}

func TestCheckStages_NoErrors(t *testing.T) {
	originalCheckPaths := checkPaths
	checkPathsFunc = func(stage Stage, errCh chan<- error) {}
	defer func() { checkPathsFunc = originalCheckPaths }()

	m := &Module{
		Stages: []Stage{
			{Name: "stage1"},
			{Name: "stage2"},
		},
	}

	err := CheckStages(m)
	if err != nil {
		t.Errorf("CheckStages() error = %v, want nil", err)
	}
}

func TestCheckStages_WithErrors(t *testing.T) {
	originalCheckPaths := checkPaths
	checkPathsFunc = func(stage Stage, errCh chan<- error) {
		if stage.Name == "fail" {
			errCh <- fmt.Errorf("failed stage: %s", stage.Name)
		}
	}
	defer func() { checkPathsFunc = originalCheckPaths }()

	m := &Module{
		Stages: []Stage{
			{Name: "ok"},
			{Name: "fail"},
		},
	}

	err := CheckStages(m)
	if err == nil {
		t.Errorf("CheckStages() error = %v, want error", err)
	} else {
		expectedMsg := "errors: [failed stage: fail]"
		if err.Error() != expectedMsg {
			t.Errorf("CheckStages() error = %v, want %v", err, expectedMsg)
		}
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
			if got := isValidPath(tt.args.filePath, tt.args.basePath); got != tt.want {
				t.Errorf("isValidPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkPaths(t *testing.T) {
	errCh := make(chan error)
	stage := Stage{
		From: []string{"."},
	}

	checkPaths(stage, errCh)
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		t.Errorf("checkPaths() = %v, want %v", errs, []error{})
	}
}

func TestHandleStages_NilModule(t *testing.T) {
	t.Run("nil module", func(t *testing.T) {
		err := HandleStages([]string{}, nil, nil, nil, &FakeBuildLogger{}, true)
		if !errors.Is(err, NilModuleError) {
			t.Errorf("HandleStages() error = %v, want %v", err, NilModuleError)
		}
	})
}

func TestHandleStages_NilContext(t *testing.T) {
	m := Module{
		Ctx: context.TODO(),
	}
	t.Run("todo context", func(t *testing.T) {
		err := HandleStages([]string{"fake-stage"}, &m, nil, nil, &FakeBuildLogger{}, true)
		if !errors.Is(err, TODOContextError) {
			t.Errorf("HandleStages() error = %v, want %v", err, TODOContextError)
		}
	})
}

func TestHandleStages_StageNotFound(t *testing.T) {
	m := Module{
		Ctx: context.Background(),
	}
	var wg sync.WaitGroup
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages([]string{"fake-stage"}, &m, &wg, nil, &FakeBuildLogger{}, true)
		if err == nil {
			t.Error("err is nil")
		}
	})
}

func TestHandleStages_NoCustomCommandMode(t *testing.T) {
	m := &Module{
		Ctx: context.Background(),
		Stages: []Stage{
			{Name: "some-fake-stage"},
		},
	}
	var wg sync.WaitGroup
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages([]string{"some-fake-stage"}, m, &wg, nil, &FakeBuildLogger{}, false)
		if !errors.Is(err, NilModuleError) {
			t.Errorf("HandleStages() error = %v, want %v", err, NilModuleError)
		}
	})
}

func TestHandleStages_Ok(t *testing.T) {
	m := &Module{
		Ctx:            context.Background(),
		BuildDirectory: "./testdata",
		Stages: []Stage{
			{Name: "some-fake-stage"},
		},
	}
	var wg sync.WaitGroup
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages([]string{"some-fake-stage"}, m, &wg, nil, &FakeBuildLogger{}, false)
		if err != nil {
			t.Errorf("HandleStages() error = %v, want nil", err)
		}
	})
}
