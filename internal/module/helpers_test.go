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
)

type FakeBuildLogger struct {
	Logs []string
	mu   sync.Mutex
}

func (l *FakeBuildLogger) Info(msg string, _ ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Logs = append(l.Logs, msg)
}

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

func Test_makeVersionFile(t *testing.T) {
	defer func() {
		stat, err := os.Stat("testdata")
		if err != nil {
			return
		}

		if stat.IsDir() {
			_ = os.RemoveAll(stat.Name())
		}
	}()

	type args struct {
		builder *ModuleBuilder
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{builder: &ModuleBuilder{module: &Module{}}}, "empty module", true},
		{args{builder: &ModuleBuilder{module: &Module{LastVersion: true}}}, "last version", false},
		{args{builder: &ModuleBuilder{module: &Module{
			BuildDirectory: "testdata",
			Version:        "1.0.0",
		}}}, "valid version", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeVersionFile(tt.args.builder); (err != nil) != tt.wantErr {
				t.Errorf("makeVersionFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_makeVersionDescription(t *testing.T) {
	defer func() {
		stat, err := os.Stat("testdata")
		if err != nil {
			return
		}

		if stat.IsDir() {
			_ = os.RemoveAll(stat.Name())
		}
	}()

	type args struct {
		builder *ModuleBuilder
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
			},
			logger: nil,
		}}, "empty repository", false},
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
				Description:    "some description",
			},
			logger: nil,
		}}, "has description", false},
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
				Repository:     ".",
			},
			logger: nil,
		}}, "has repository", false},
		{args{builder: &ModuleBuilder{module: &Module{LastVersion: true}}}, "last version", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeVersionDescription(tt.args.builder); (err != nil) != tt.wantErr {
				t.Errorf("makeVersionDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
