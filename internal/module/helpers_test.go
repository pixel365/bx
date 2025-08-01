package module

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	t.Parallel()
	_, err := ReadModuleFromFlags(nil)
	require.Error(t, err)
	assert.ErrorIs(t, err, errors2.ErrNilCmd)
}

func TestReadModuleFromFlags_Name(t *testing.T) {
	t.Parallel()
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, helpers.RootDir))
	_, err := ReadModuleFromFlags(cmd)
	require.Error(t, err)
}

func TestReadModuleFromFlags_File(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), helpers.RootDir, helpers.RootDir))
	cmd.SetArgs([]string{"--file", "./test_files/foo"})
	_, err := ReadModuleFromFlags(cmd)
	require.Error(t, err)
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
			got := AllModules(tt.args.directory)
			assert.Equal(t, tt.want, got)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
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
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
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
			err := makeVersionFile(tt.args.builder)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
			log: nil,
		}}, "empty repository", false},
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
				Description:    "some description",
			},
			log: nil,
		}}, "has description", false},
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
				Repository:     ".",
			},
			log: nil,
		}}, "has repository", false},
		{args{builder: &ModuleBuilder{module: &Module{LastVersion: true}}}, "last version", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := makeVersionDescription(tt.args.builder)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_versionPhpContent(t *testing.T) {
	t.Parallel()
	date, err := time.Parse(time.RFC3339, "2025-05-20T23:00:00Z")
	require.NoError(t, err)

	buf := versionPhpContent("1.0.0", date)
	buf2 := strings.Builder{}
	buf2.WriteString("<?php\n")
	buf2.WriteString("$arModuleVersion = array(\n")
	buf2.WriteString("\t\t\"VERSION\" => \"1.0.0\",\n")
	buf2.WriteString("\t\t\"VERSION_DATE\" => \"" + date.Format(time.DateTime) + "\",\n")
	buf2.WriteString(");\n")

	assert.Equal(t, buf.String(), buf2.String())
}
