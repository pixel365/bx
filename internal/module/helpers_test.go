package module

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/spf13/cobra"

	errors2 "github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/types"
)

type FakeBuildLogger struct{}

func (l *FakeBuildLogger) Info(_ string, _ ...interface{}) {}

func (l *FakeBuildLogger) Error(_ string, _ error, _ ...interface{}) {}

func (l *FakeBuildLogger) Cleanup() {}

func TestReadModuleFromFlags(t *testing.T) {
	t.Run("TestReadModuleFromFlags", func(t *testing.T) {
		_, err := ReadModuleFromFlags(nil)
		if err == nil {
			t.Errorf("ReadModuleFromFlags() did not return an error")
		}

		if !errors.Is(err, errors2.ErrNilCmd) {
			t.Errorf("err = %v, want %v", err, errors2.ErrNilCmd)
		}
	})
}

func TestReadModuleFromFlags_Name(t *testing.T) {
	t.Run("TestReadModuleFromFlags_Name", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, helpers.RootDir))
		_, err := ReadModuleFromFlags(cmd)
		if err == nil {
			t.Errorf("ReadModuleFromFlags() did not return an error")
		}
	})
}

func TestReadModuleFromFlags_File(t *testing.T) {
	t.Run("TestReadModuleFromFlags_File", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, helpers.RootDir))
		cmd.SetArgs([]string{"--file", "./test_files/foo"})
		_, err := ReadModuleFromFlags(cmd)
		if err == nil {
			t.Errorf("ReadModuleFromFlags() did not return an error")
		}
	})
}

func TestAllModules(t *testing.T) {
	name := fmt.Sprintf("%s_%d", "testing", time.Now().Unix())
	filePath, err := filepath.Abs(fmt.Sprintf("./%s/%s.yaml", ".", name))
	if err != nil {
		t.Error()
	}
	filePath = filepath.Clean(filePath)

	err = os.WriteFile(filePath, []byte(helpers.DefaultYAML()), 0600)
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

func TestHandleStages_NilModule(t *testing.T) {
	ctx := context.Background()
	t.Run("nil module", func(t *testing.T) {
		err := HandleStages(ctx, []string{}, nil, nil, nil, &FakeBuildLogger{}, true)
		if !errors.Is(err, errors2.ErrNilModule) {
			t.Errorf("HandleStages() error = %v, want %v", err, errors2.ErrNilModule)
		}
	})
}

func TestHandleStages_NilContext(t *testing.T) {
	ctx := context.TODO()
	m := Module{}
	t.Run("todo context", func(t *testing.T) {
		err := HandleStages(ctx, []string{"fake-stage"}, &m, nil, nil, &FakeBuildLogger{}, true)
		if !errors.Is(err, errors2.ErrTODOContext) {
			t.Errorf("HandleStages() error = %v, want %v", err, errors2.ErrTODOContext)
		}
	})
}

func TestHandleStages_StageNotFound(t *testing.T) {
	ctx := context.Background()
	m := Module{}
	var wg sync.WaitGroup
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages(ctx, []string{"fake-stage"}, &m, &wg, nil, &FakeBuildLogger{}, true)
		if err == nil {
			t.Error("err is nil")
		}
	})
}

func TestHandleStages_NoCustomCommandMode(t *testing.T) {
	ctx := context.Background()
	m := &Module{
		Stages: []types.Stage{
			{Name: "some-fake-stage"},
		},
	}
	var wg sync.WaitGroup
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages(
			ctx,
			[]string{"some-fake-stage"},
			m,
			&wg,
			nil,
			&FakeBuildLogger{},
			false,
		)
		if !errors.Is(err, errors2.ErrNilModule) {
			t.Errorf("HandleStages() error = %v, want %v", err, errors2.ErrNilModule)
		}
	})
}

func TestHandleStages_Ok(t *testing.T) {
	ctx := context.Background()
	m := &Module{
		BuildDirectory: "./testdata",
		Stages: []types.Stage{
			{Name: "some-fake-stage"},
		},
	}
	var wg sync.WaitGroup
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages(
			ctx,
			[]string{"some-fake-stage"},
			m,
			&wg,
			nil,
			&FakeBuildLogger{},
			false,
		)
		if err != nil {
			t.Errorf("HandleStages() error = %v, want nil", err)
		}
	})
}

func Test_makeVersionDirectory(t *testing.T) {
	mod1 := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	mod2 := &Module{
		BuildDirectory: "testdata/build",
		Version:        "1.0.1",
	}

	mod3 := &Module{
		BuildDirectory: "",
		Version:        "1.0.1",
	}

	cur, _ := os.Getwd()

	type args struct {
		module *Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{mod1}, fmt.Sprintf("%s/testdata/1.0.0", cur), false},
		{"2", args{mod2}, fmt.Sprintf("%s/testdata/build/1.0.1", cur), false},
		{"nil module", args{nil}, "", true},
		{"empty build directory", args{mod3}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeVersionDirectory(tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeVersionDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeVersionDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeZipFilePath(t *testing.T) {
	mod1 := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	mod2 := &Module{
		BuildDirectory: "testdata/build",
		Version:        "1.0.1",
	}

	cur, _ := os.Getwd()

	type args struct {
		module *Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{mod1}, fmt.Sprintf("%s/testdata/1.0.0.zip", cur), false},
		{"2", args{mod2}, fmt.Sprintf("%s/testdata/build/1.0.1.zip", cur), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeZipFilePath(tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeZipFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeZipFilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckStages(t *testing.T) {
	err := CheckStages(nil)
	if !errors.Is(err, errors2.ErrNilModule) {
		t.Errorf("CheckStages() error = %v, want %v", err, errors2.ErrNilModule)
	}
}

func TestCheckStages_NoErrors(t *testing.T) {
	originalCheckPaths := helpers.CheckPaths
	checkPathsFunc = func(stage types.Stage, errCh chan<- error) {}
	defer func() { checkPathsFunc = originalCheckPaths }()

	m := &Module{
		Stages: []types.Stage{
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
	originalCheckPaths := helpers.CheckPaths
	checkPathsFunc = func(stage types.Stage, errCh chan<- error) {
		if stage.Name == "fail" {
			errCh <- fmt.Errorf("failed stage: %s", stage.Name)
		}
	}
	defer func() { checkPathsFunc = originalCheckPaths }()

	m := &Module{
		Stages: []types.Stage{
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
