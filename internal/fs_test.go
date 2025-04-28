package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func Test_mkdir(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		path, err := mkdir(name)
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
		path, err := mkdir(name)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()

		archivePath := fmt.Sprintf("./_%d.zip", time.Now().UTC().Unix())
		if err := zipIt(path, archivePath); err != nil {
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
		patterns *[]string
		path     string
	}
	tests := []struct {
		args args
		name string
		want bool
	}{
		{args{nil, "."}, "1", false},
		{args{&patterns, "./testing/errors.log"}, "2", true},
		{args{&patterns, "./testing/errors.json"}, "3", false},
		{args{&patterns, "errors.json"}, "4", true},
		{args{&patterns, "./testing/data/awesome/cfg.yaml"}, "5", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldSkip(tt.args.path, tt.args.patterns); got != tt.want {
				t.Errorf("shouldSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyFromPath(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		from := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		fromPath, err := mkdir(from)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(fromPath); err != nil {
				t.Error(err)
			}
		}()

		to := fmt.Sprintf("./__%d", time.Now().UTC().Unix())
		toPath, err := mkdir(to)
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
		patterns := []string{
			"**/*.log",
			"*.json",
			"**/*some*/*",
		}

		module := Module{Ignore: patterns}

		wg.Add(1)
		copyFromPath(
			context.Background(),
			&wg,
			errChan,
			&module,
			from,
			to,
			Replace,
			false,
			[]string{},
		)

		defer func() {
			if err := os.Remove(fmt.Sprintf("%s/%s", toPath, fileName)); err != nil {
				t.Error(err)
			}
		}()
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
		path, err := mkdir(name)
		if err != nil {
			t.Error(err)
		}

		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()

		if !isEmptyDir(path) {
			t.Errorf("isEmptyDir() = %v, want %v", isEmptyDir(path), true)
		}
	})
}

func Test_removeEmptyDirs(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		name := fmt.Sprintf("./_%d", time.Now().UTC().Unix())
		name2 := fmt.Sprintf("./%s/%d", name, time.Now().UTC().Unix())
		path, err := mkdir(name)
		if err != nil {
			t.Error(err)
		}
		defer func() {
			if err := os.Remove(path); err != nil {
				t.Error(err)
			}
		}()

		path2, err := mkdir(name2)
		if err != nil {
			t.Error(err)
		}

		status, err := removeEmptyDirs(path)
		if err != nil {
			t.Errorf("removeEmptyDirs() error = %v", err)
		}
		if !status {
			t.Errorf("removeEmptyDirs() = %v, want %v", status, true)
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
