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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	t.Parallel()

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
	assert.Equal(t, def, DefaultYAML())
}

func TestGetModulesDir(t *testing.T) {
	_, err := GetModulesDir()
	require.NoError(t, err)
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
			err := CheckPath(tt.args.path)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCheckContextOk(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	err := CheckContext(ctx)
	require.NoError(t, err)
}

func TestCheckContextDone(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := CheckContext(ctx)
	require.Error(t, err)
}

func TestCheckContextNil(t *testing.T) {
	t.Parallel()
	err := CheckContext(nil) //nolint:staticcheck // SA1012: passing nil context is intentional
	require.Error(t, err)
	assert.ErrorIs(t, err, errors2.ErrNilContext)
}

func TestCaptureOutput(t *testing.T) {
	t.Parallel()
	output := CaptureOutput(func() {
		fmt.Println("ok")
	})

	assert.Equal(t, "ok\n", output)

	output = CaptureOutput(func() {
		fmt.Println("ok string")
	})

	assert.Equal(t, "ok string\n", output)
}

func TestReplaceVariables(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			got, err := ReplaceVariables(tt.args.input, tt.args.variables, tt.args.depth)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_Cleanup(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			assert.NotPanics(t, func() {
				Cleanup(tt.args.resource, nil)
			})
		})
	}
}

func Test_Cleanup_chan(t *testing.T) {
	ch := make(chan error, 1)

	fileName := fmt.Sprintf("mod-%d.yaml", time.Now().UTC().Unix())
	filePath, _ := filepath.Abs("./" + fileName)
	filePath = filepath.Clean(filePath)

	file, err := os.Create(filePath)
	assert.NoError(t, err)

	defer func() {
		err := os.Remove(file.Name())
		assert.NoError(t, err)
	}()

	go func() {
		Cleanup(file, ch)
		close(ch)
	}()

	for e := range ch {
		t.Errorf("TestCleanup_chan failed: %v", e)
	}
}

func TestChoose(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			err := Choose(tt.args.items, tt.args.value, tt.args.title)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
			got := IsValidPath(tt.args.filePath, tt.args.basePath)
			assert.Equal(t, tt.want, got)
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

	assert.Len(t, errs, 0)
}

func TestUserInput_success(t *testing.T) {
	prompter := FakePromptSuccessor{}
	value := ""
	err := UserInput(&prompter, &value, "title", func(string) error { return nil })
	require.NoError(t, err)
}

func TestUserInput_fail(t *testing.T) {
	prompter := FakePromptFailer{}
	value := ""
	err := UserInput(
		&prompter,
		&value,
		"title",
		func(string) error { return errors.New("fake error") },
	)
	require.Error(t, err)
}

func TestSortSemanticVersions(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			got := SortSemanticVersions(tt.args.versions)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_normalizeVersion(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			got := normalizeVersion(tt.args.version)
			assert.Equal(t, tt.want, got)
		})
	}
}
