package fs

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

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

		var wg sync.WaitGroup
		errChan := make(chan error)

		module := FakeModuleConfig{}

		wg.Add(1)
		CopyFromPath(
			context.Background(),
			&wg,
			errChan,
			&module,
			from,
			to,
			types.Replace,
			false,
			[]string{},
		)

		close(errChan)

		defer func() {
			if err := os.Remove(fmt.Sprintf("%s/%s", toPath, fileName)); err != nil {
				t.Error(err)
			}
		}()

		wg.Wait()
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
