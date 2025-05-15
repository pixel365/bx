package fs

import (
	"context"
	errors2 "errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/types"
)

type FakeModuleConfig struct{}

func (f FakeModuleConfig) GetVariables() map[string]string { return nil }
func (f FakeModuleConfig) GetRun() map[string][]string     { return nil }
func (f FakeModuleConfig) GetStages() []types.Stage        { return nil }
func (f FakeModuleConfig) GetIgnore() []string {
	return []string{
		"**/*.log",
		"*.json",
		"**/*some*/*",
	}
}
func (f FakeModuleConfig) GetChanges() *types.Changes { return nil }
func (f FakeModuleConfig) IsLastVersion() bool        { return false }

type FakeFileInfo struct {
	Dir bool
}

func (f FakeFileInfo) Name() string       { return "" }
func (f FakeFileInfo) Size() int64        { return 0 }
func (f FakeFileInfo) Mode() fs.FileMode  { return fs.ModeDir }
func (f FakeFileInfo) ModTime() time.Time { return time.Time{} }
func (f FakeFileInfo) Sys() interface{}   { return nil }
func (f FakeFileInfo) IsDir() bool        { return f.Dir }

func Test_mkdir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		path, err := MkDir(name)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()
	})
}

func Test_zipIt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		path, err := MkDir(name)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()

		archivePath := fmt.Sprintf("./_%d.zip", time.Now().UTC().Unix())
		if err := ZipIt(path, archivePath); err != nil {
			t.Error(err)
		}
		defer func() {
			if err := os.Remove(archivePath); err != nil {
				t.Error(err)
			}
		}()
	})
}

func Test_shouldSkip(t *testing.T) {
	patterns := []string{
		"**/*.log",
		"*.json",
		"**/*some*/*",
	}
	type args struct {
		path     string
		patterns []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{".", nil}, false},
		{"2", args{"./testing/errors.log", patterns}, true},
		{"3", args{"./testing/errors.json", patterns}, false},
		{"4", args{"errors.json", patterns}, true},
		{"5", args{"./testing/data/awesome/cfg.yaml", patterns}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkip(tt.args.path, tt.args.patterns); got != tt.want {
				t.Errorf("shouldSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CopyFromPath_ok(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		from := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		fromPath, err := MkDir(from)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(fromPath); err != nil {
				t.Error(err)
			}
		}()

		to := fmt.Sprintf("./__%d", time.Now().UTC().Unix())
		toPath, err := MkDir(to)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(toPath); err != nil {
				t.Error(err)
			}
		}()

		fileName := fmt.Sprintf("%d.txt", time.Now().UTC().Unix())
		filePath := filepath.Join(from, fileName)
		filePath = filepath.Clean(filePath)
		file, err := os.Create(filePath)
		if err != nil {
			t.Error(err)
		}

		err = file.Close()
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(filePath); err != nil {
				t.Error(err)
			}
		}()

		errChan := make(chan types.Path, 1)

		module := FakeModuleConfig{}

		path := types.Path{
			From:           from,
			To:             to,
			ActionIfExists: types.Replace,
			Convert:        false,
		}

		if err = PathProcessing(
			context.Background(),
			errChan,
			&module,
			path,
			[]string{},
		); err != nil {
			close(errChan)
			t.Error(err)
		}

		defer func() {
			_ = os.Remove(fmt.Sprintf("%s/%s", toPath, fileName))
		}()

		close(errChan)
	})
}

func TestPathProcessingContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	t.Run("cancelled context", func(t *testing.T) {
		err := PathProcessing(ctx, nil, nil, types.Path{}, nil)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func Test_isConvertable(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{args{"/some/lang/file.php"}, "php", true},
		{args{"/some/path/file.php"}, "php", false},
		{args{"/some/path/description.ru"}, "description.ru", true},
		{args{"/some/path/image.jpg"}, "jpg", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isConvertable(tt.args.path); got != tt.want {
				t.Errorf("isConvertable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isEmptyDir(t *testing.T) {
	t.Run("empty dir", func(t *testing.T) {
		name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		path, err := MkDir(name)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()

		if !IsEmptyDir(path) {
			t.Errorf("IsEmptyDir() = %v, want %v", IsEmptyDir(path), true)
		}
	})
}

func Test_removeEmptyDirs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		name2 := fmt.Sprintf("./%s/%d", name, time.Now().UTC().Unix())
		path, err := MkDir(name)
		if err != nil {
			t.Error(err)
		}
		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()

		path2, err := MkDir(name2)
		if err != nil {
			t.Error(err)
		}

		status, err := RemoveEmptyDirs(path)
		if err != nil {
			t.Errorf("RemoveEmptyDirs() error = %v", err)
		}
		if !status {
			t.Errorf("RemoveEmptyDirs() = %v, want %v", status, true)
		}

		if !status || err != nil {
			defer func() {
				if err := os.Remove(path2); err != nil {
					t.Error(err)
				}
			}()
		}
	})
}

func Test_shouldInclude(t *testing.T) {
	type args struct {
		path     string
		patterns []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty patterns", args{"./testing.php", []string{}}, true},
		{"empty path", args{"", []string{"**/*.php"}}, true},
		{"included path", args{"./testing.php", []string{"**/*.php"}}, true},
		{"excluded json", args{"./testing.json", []string{"!**/*.json"}}, false},
		{"included php", args{"./testing.php", []string{"!**/*.json"}}, true},
		{
			"excluded test file",
			args{"./some_test.php", []string{"**/*.php", "!**/*_test.php"}},
			false,
		},
		{
			"mutually exclusive rules",
			args{"./testing.php", []string{"**/*.php", "!**/*.php"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldInclude(tt.args.path, tt.args.patterns); got != tt.want {
				t.Errorf("shouldInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsFileExists_true(t *testing.T) {
	filePath := fmt.Sprintf("./%d.txt", time.Now().UTC().Unix())
	filePath = filepath.Clean(filePath)
	file, err := os.Create(filePath)
	if err != nil {
		t.Error(err)
	}

	_, err = file.WriteString("str")
	if err != nil {
		t.Error(err)
	}

	err = file.Close()
	if err != nil {
		t.Error(err)
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Error(err)
		}
	}()

	t.Run("file exists", func(t *testing.T) {
		ok, size := IsFileExists(filePath)
		if !ok || size == 0 {
			t.Errorf("IsFileExists() = %v, want %v", ok, true)
		}
	})
}

func TestIsFileExists_false(t *testing.T) {
	t.Run("file exists", func(t *testing.T) {
		ok, size := IsFileExists("./some-file.txt")
		if ok || size > 0 {
			t.Errorf("IsFileExists() = %v, want %v", ok, false)
		}
	})
}

func Test_skip(t *testing.T) {
	type args struct {
		info os.FileInfo
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{info: nil}, "nil info", false},
		{args{info: FakeFileInfo{Dir: true}}, "is dir", true},
		{args{info: FakeFileInfo{Dir: false}}, "is not dir", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := skip(tt.args.info); (err != nil) != tt.wantErr {
				t.Errorf("skip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_visitor(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		cfg         interfaces.ModuleConfig
		err         error
		filesCh     chan<- types.Path
		path        types.Path
		filterRules []string
	}
	tests := []struct {
		want filepath.WalkFunc
		name string
		args args
	}{
		{want: func(_ string, _ fs.FileInfo, _ error) error {
			return errors.ErrNilContext
		}, name: "nil context", args: args{cfg: FakeModuleConfig{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			visit := visitor(
				ctx,
				tt.args.filesCh,
				tt.args.cfg,
				tt.args.path,
				tt.args.filterRules,
			)
			if visit == nil {
				t.Errorf("visitor() = %v, want non-nil", visit)
			}

			if err := visit("", FakeFileInfo{}, tt.args.err); !errors2.Is(err, context.Canceled) {
				t.Errorf("visit() error = %v, want %v", err, context.Canceled)
			}
		})
	}
}
