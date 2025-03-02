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

		wg.Add(1)
		copyFromPath(context.Background(), &wg, errChan, &patterns, from, to, Replace, false)

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
		{args{"/some/path/file.php"}, "php", true},
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
