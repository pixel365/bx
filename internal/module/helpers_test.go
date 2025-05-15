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

func TestHandleStages_NoCustomCommandMode(t *testing.T) {
	ctx := context.Background()
	m := &Module{
		Stages: []types.Stage{
			{Name: "some-fake-stage"},
		},
	}
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages(
			ctx,
			[]string{"some-fake-stage"},
			m,
			&FakeBuildLogger{},
			false,
		)
		if !errors.Is(err, errors2.ErrNilModule) {
			t.Errorf("HandleStages() error = %v, want %v", err, errors2.ErrNilModule)
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

func TestCopyWorkers(t *testing.T) {
	var mu sync.Mutex
	var called []types.Path

	copyFileFunc = func(ctx context.Context, errCh chan<- error, path types.Path) {
		mu.Lock()
		called = append(called, path)
		mu.Unlock()
	}

	filesCh := make(chan types.Path, 3)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	copyWorkers(ctx, &wg, filesCh, errCh, 2)

	filesCh <- types.Path{From: "a.txt", To: "x"}
	filesCh <- types.Path{From: "b.txt", To: "y"}
	close(filesCh)

	wg.Wait()

	if len(called) != 2 {
		t.Errorf("expected 2 calls, got %d", len(called))
	}
}

func TestErrorWorker(t *testing.T) {
	errCh := make(chan error, 2)
	var once sync.Once
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var capturedErr error
	go errorWorker(errCh, cancel, &once, &capturedErr)

	expectedErr := errors.New("fail")
	errCh <- expectedErr
	time.Sleep(50 * time.Millisecond)

	if !errors.Is(capturedErr, expectedErr) {
		t.Errorf("expected %v, got %v", expectedErr, capturedErr)
	}
}

func TestLogWorker(t *testing.T) {
	logCh := make(chan string, 2)
	mock := &FakeBuildLogger{}

	go logWorker(logCh, mock)
	logCh <- "hello"
	logCh <- "world"
	close(logCh)

	time.Sleep(50 * time.Millisecond)

	if len(mock.Logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(mock.Logs))
	}
	if mock.Logs[0] != "hello" || mock.Logs[1] != "world" {
		t.Errorf("unexpected logs: %v", mock.Logs)
	}
}

func TestCleanupWorker(t *testing.T) {
	var stageWg, copyWg sync.WaitGroup
	stageWg.Add(1)
	copyWg.Add(1)

	filesCh := make(chan types.Path, 1)
	logCh := make(chan string, 1)
	errCh := make(chan error, 1)

	var canceled bool
	cancel := func() { canceled = true }

	var once sync.Once
	go cleanupWorker(&stageWg, &copyWg, &once, cancel, filesCh, logCh, errCh)

	stageWg.Done()
	copyWg.Done()

	time.Sleep(50 * time.Millisecond)

	select {
	case _, ok := <-filesCh:
		if ok {
			t.Error("filesCh should be closed")
		}
	default:
		t.Error("filesCh not closed")
	}

	if !canceled {
		t.Error("cancel not called")
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
