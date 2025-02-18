package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bmatcuk/doublestar/v4"
)

func copyFromPath(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	ignore *[]string,
	from, to string,
	existsMode FileExistsMode,
) {
	defer wg.Done()

	if err := CheckContextActivity(ctx); err != nil {
		errCh <- err
		return
	}

	if err := walk(ctx, wg, errCh, from, to, ignore, existsMode); err != nil {
		if !errors.Is(err, doublestar.SkipDir) {
			errCh <- err
		}
	}
}

func walk(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	from, to string,
	patterns *[]string,
	existsMode FileExistsMode,
) error {
	wg.Add(1)
	defer wg.Done()

	var wg2 sync.WaitGroup
	jobs := make(chan struct{}, 10)

	err := filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if ctxErr := CheckContextActivity(ctx); ctxErr != nil {
			errCh <- ctxErr
			return ctxErr
		}

		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(from, path)
		if err != nil {
			return err
		}

		if shouldSkip(relPath, patterns) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		absFrom, err := filepath.Abs(fmt.Sprintf("%s/%s", from, relPath))
		if err != nil {
			return err
		}
		absFrom = filepath.Clean(absFrom)

		isDir, err := IsDir(absFrom)
		if err != nil {
			return err
		}

		if !isDir {
			absTo, err := filepath.Abs(fmt.Sprintf("%s/%s", to, relPath))
			if err != nil {
				return err
			}
			absTo = filepath.Clean(absTo)
			toDir := filepath.Dir(absTo)

			if _, err := mkdir(toDir); err != nil {
				return err
			}

			wg2.Add(1)
			go copyFile(ctx, &wg2, errCh, absFrom, absTo, jobs, existsMode)
		}

		return nil
	})

	wg2.Wait()

	return err
}

func copyFile(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	src, dst string,
	jobs chan struct{},
	existsMode FileExistsMode,
) {
	defer wg.Done()

	if err := CheckContextActivity(ctx); err != nil {
		errCh <- err
		return
	}

	fileName := strings.LastIndex(src, "/")
	if !strings.HasSuffix(dst, src[fileName:]) {
		dst = filepath.Join(dst, src[fileName:])
		dst = filepath.Clean(dst)
	}

	var existingFile os.FileInfo = nil
	stat, err := os.Stat(dst)
	if err == nil {
		existingFile = stat
		if existsMode == Skip {
			return
		}
	}

	jobs <- struct{}{}
	defer func() { <-jobs }()

	in, err := os.Open(src)
	if err != nil {
		errCh <- err
		return
	}

	defer func(in *os.File) {
		err := in.Close()
		if err != nil {
			errCh <- err
		}
	}(in)

	if err := CheckContextActivity(ctx); err != nil {
		errCh <- err
		return
	}

	info, err := in.Stat()
	if err != nil {
		errCh <- err
		return
	}

	allowWrite := true
	if existingFile != nil {
		if existsMode == CopyNew {
			allowWrite = info.ModTime().After(existingFile.ModTime())
		}
	}

	if allowWrite {
		out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			errCh <- err
			return
		}

		defer func(out *os.File) {
			err := out.Close()
			if err != nil {
				errCh <- err
			}
		}(out)

		if err := CheckContextActivity(ctx); err != nil {
			errCh <- err
			return
		}

		_, err = io.Copy(out, in)
		if err != nil {
			errCh <- err
			return
		}

		if err := CheckContextActivity(ctx); err != nil {
			errCh <- err
			return
		}

		err = os.Chtimes(dst, info.ModTime(), info.ModTime())
		if err != nil {
			errCh <- err
			return
		}
	}
}

func shouldSkip(path string, patterns *[]string) bool {
	for _, pattern := range *patterns {
		if ok, err := doublestar.PathMatch(pattern, path); ok || err != nil {
			return true
		}
	}
	return false
}

func mkdir(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	path = filepath.Clean(path)

	if !isValidPath(path, path) {
		return "", fmt.Errorf("invalid path: %s", path)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0750)
		if err != nil {
			return "", err
		}
	}

	return path, nil
}
