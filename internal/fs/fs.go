package fs

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/types"

	"golang.org/x/text/encoding/charmap"

	"github.com/pixel365/bx/internal/helpers"

	"github.com/bmatcuk/doublestar/v4"
)

// CopyFromPath performs the process of copying files from a source directory to a target directory.
//
// The function checks if the context is valid, and if so, it invokes the `walk` function to traverse the source
// directory and perform the copy operation. It uses a `sync.WaitGroup` to wait for completion, and reports
// any errors encountered through the provided error channel.
//
// Parameters:
//   - ctx (context.Context): The context to control the execution and cancellation of the operation.
//   - wg (*sync.WaitGroup): The wait group to synchronize the completion of the operation.
//   - errCh (chan<- error): A channel for reporting errors encountered during the operation.
//   - cfg (internal.ModuleConfig): Module instance
//   - from (string): The source directory to copy from.
//   - to (string): The destination directory to copy to.
//   - existsMode (FileExistsAction): The action to take if the file already exists in the destination.
//   - convert (bool): A flag to indicate whether to convert the files during the copy process.
//   - filterRules ([]string): Filter rules for selective copying of files based on patterns.
//
// If any error is encountered during the execution, it is reported through the `errCh` channel.
func CopyFromPath(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	cfg interfaces.ModuleConfig,
	from, to string,
	existsMode types.FileExistsAction,
	convert bool,
	filterRules []string,
) {
	defer wg.Done()

	if err := helpers.CheckContext(ctx); err != nil {
		errCh <- err
		return
	}

	if err := walk(ctx, wg, errCh, from, to, cfg, existsMode, convert, filterRules); err != nil {
		if !errors.Is(err, doublestar.SkipDir) {
			errCh <- err
		}
	}
}

// walk traverses the source directory recursively and copies files to the destination directory
// while checking the context for cancellation. It processes the paths according to the specified patterns
// and handles existing files based on the provided `existsMode`.
//
// The function uses a wait group to synchronize the completion of all file copy operations and reports
// any errors encountered during the process through the provided error channel.
//
// Parameters:
//   - ctx (context.Context): The context to control the execution and cancellation of the operation.
//   - wg (*sync.WaitGroup): The wait group to synchronize the completion of the operation.
//   - errCh (chan<- error): A channel for reporting errors encountered during the operation.
//   - from (string): The source directory to walk through.
//   - to (string): The destination directory where files should be copied.
//   - cfg (internal.ModuleConfig): Module instance
//   - existsMode (FileExistsAction): The action to take if the file already exists in the destination directory.
//   - convert (bool): A flag to indicate whether to convert files during the copy process.
//   - filterRules ([]string): Filter rules for selective copying of files based on patterns.
//
// Returns:
//   - error: If an error occurs during the traversal or file copying process, it will be returned. If the
//     process completes successfully, it returns nil.
func walk(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	from, to string,
	cfg interfaces.ModuleConfig,
	existsMode types.FileExistsAction,
	convert bool,
	filterRules []string,
) error {
	wg.Add(1)
	defer wg.Done()

	var wg2 sync.WaitGroup
	jobs := make(chan struct{}, 10)

	var changes *types.Changes = nil
	if !cfg.IsLastVersion() {
		changes = cfg.GetChanges()
	}

	skipFn := func(info os.FileInfo) error {
		if info.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}

	err := filepath.Walk(from, func(path string, info os.FileInfo, err error) error {
		if ctxErr := helpers.CheckContext(ctx); ctxErr != nil {
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

		absFrom, err := filepath.Abs(fmt.Sprintf("%s/%s", from, relPath))
		if err != nil {
			return err
		}
		absFrom = filepath.Clean(absFrom)

		if shouldSkip(relPath, cfg.GetIgnore()) {
			return skipFn(info)
		}

		if !shouldInclude(absFrom, filterRules) {
			return skipFn(info)
		}

		isDir, err := helpers.IsDir(absFrom)
		if err != nil {
			return err
		}

		if !isDir {
			if changes != nil {
				if !changes.IsChangedFile(absFrom) {
					return nil
				}
			}

			absTo, err := filepath.Abs(fmt.Sprintf("%s/%s", to, relPath))
			if err != nil {
				return err
			}
			absTo = filepath.Clean(absTo)
			toDir := filepath.Dir(absTo)

			if _, err := MkDir(toDir); err != nil {
				return err
			}

			wg2.Add(1)
			go copyFile(ctx, &wg2, errCh, absFrom, absTo, jobs, existsMode, convert)
		}

		return nil
	})

	wg2.Wait()

	return err
}

// copyFile copies a file from the source path to the destination path, taking into account the context cancellation,
// existing file handling mode (`existsMode`), and optional conversion of file content.
//
// The function checks if the destination file exists and takes action based on the specified `existsMode`:
//   - If `Skip`, it does nothing if the file already exists.
//   - If `ReplaceIfNewer`, it replaces the file only if the source file is newer than the destination file.
//   - If neither of the above, it always replaces the file.
//
// If the `convert` flag is set to `true` and the file is convertible (determined by the `isConvertable` function),
// it converts the file content using Windows-1251 encoding before writing it to the destination.
//
// The function handles errors by sending them to the provided error channel and checks the context at various stages
// to support cancellation during file copying.
//
// Parameters:
//   - ctx (context.Context): The context to control the execution and cancellation of the operation.
//   - wg (*sync.WaitGroup): The wait group to synchronize the completion of the operation.
//   - errCh (chan<- error): A channel for reporting errors encountered during the operation.
//   - src (string): The source file path to copy from.
//   - dst (string): The destination file path to copy to.
//   - jobs (chan struct{}): A channel for managing concurrent file copy operations with a limited number of concurrent jobs.
//   - existsMode (FileExistsAction): The action to take if the file already exists in the destination directory.
//   - convert (bool): A flag to indicate whether to convert the file content during the copy process.
//
// Returns:
//   - None: Errors are sent to the error channel `errCh`, and no return value is provided.
func copyFile(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	src, dst string,
	jobs chan struct{},
	existsMode types.FileExistsAction,
	convert bool,
) {
	defer wg.Done()

	if err := helpers.CheckContext(ctx); err != nil {
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
		if existsMode == types.Skip {
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

	defer helpers.Cleanup(in, errCh)

	if err := helpers.CheckContext(ctx); err != nil {
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
		if existsMode == types.ReplaceIfNewer {
			allowWrite = info.ModTime().After(existingFile.ModTime())
		}
	}

	if allowWrite {
		out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
		if err != nil {
			errCh <- err
			return
		}

		defer helpers.Cleanup(out, errCh)

		if err := helpers.CheckContext(ctx); err != nil {
			errCh <- err
			return
		}

		var writer io.Writer
		if convert && isConvertable(src) {
			encoder := charmap.Windows1251.NewEncoder()
			writer = encoder.Writer(out)
		} else {
			writer = out
		}

		_, err = io.Copy(writer, in)
		if err != nil {
			errCh <- err
			return
		}

		if err := helpers.CheckContext(ctx); err != nil {
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

// shouldSkip checks if a given file path should be skipped based on a list of glob patterns.
// It compares the provided file path against each pattern in the `patterns` slice using the `doublestar.PathMatch`
// function to determine if the path matches any pattern. If any pattern matches or an error occurs during the match,
// the function returns true, indicating the path should be skipped. Otherwise, it returns false.
//
// Parameters:
//   - path (string): The file path to check against the patterns.
//   - patterns (*[]string): A slice of glob patterns that define the paths to skip. If `patterns` is nil or empty,
//     the function returns false, meaning the path should not be skipped.
//
// Returns:
//   - bool: Returns `true` if the path matches any of the patterns or an error occurs, indicating the path should be skipped.
//     Returns `false` otherwise, meaning the path should not be skipped.
func shouldSkip(path string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	for _, pattern := range patterns {
		if ok, err := doublestar.PathMatch(pattern, path); ok || err != nil {
			return true
		}
	}
	return false
}

// shouldInclude determines whether a given file path should be included based on a set of patterns.
//
// It supports both inclusion and exclusion patterns:
//   - Patterns **without** a `!` prefix define files that should be **included**.
//   - Patterns **with** a `!` prefix define files that should be **excluded**.
//
// If no patterns are provided, the function returns `true`, allowing all files by default.
//
// The function follows these rules:
//  1. If no **inclusion** patterns are specified, all files are initially allowed.
//  2. A file matches **if it matches at least one inclusion pattern** (or if no inclusions are defined).
//  3. A file is **excluded** if it matches an exclusion pattern, even if it was previously included.
//
// Parameters:
//   - `path string`: The file path to evaluate.
//   - `patterns []string`: A list of patterns to determine inclusion/exclusion.
//
// Returns:
//   - `bool`: `true` if the file should be included, `false` if it should be excluded.
func shouldInclude(path string, patterns []string) bool {
	if len(patterns) == 0 || path == "" {
		return true
	}

	var include, exclude []string
	for _, pattern := range patterns {
		if strings.HasPrefix(pattern, "!") {
			exclude = append(exclude, strings.TrimPrefix(pattern, "!"))
			continue
		}
		include = append(include, pattern)
	}

	allow := len(include) == 0

	if !allow {
		for _, pattern := range include {
			if ok, err := doublestar.PathMatch(pattern, path); ok || err != nil {
				allow = true
				break
			}
		}
	}

	if allow {
		for _, pattern := range exclude {
			if ok, err := doublestar.PathMatch(pattern, path); ok || err != nil {
				allow = false
				break
			}
		}
	}

	return allow
}

// MkDir creates a directory at the specified path, including any necessary parent directories.
// It first ensures the provided path is an absolute and clean path, and checks whether the path is valid.
// If the path does not exist, it will be created with permissions 0750. If the directory already exists, it does nothing.
//
// Parameters:
//   - path (string): The path where the directory should be created. The function will resolve this to an absolute path
//     and ensure it is valid before attempting to create the directory.
//
// Returns:
//   - string: The absolute path of the created directory.
//   - error: If an error occurs during path resolution or directory creation, an error is returned.
func MkDir(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	path = filepath.Clean(path)

	if !helpers.IsValidPath(path, path) {
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

// ZipIt creates a ZIP archive from the specified directory, including all its files and subdirectories.
// It recursively walks through the directory, adding files and directories to the ZIP archive, preserving
// the relative paths of the files.
//
// Parameters:
//   - dirPath (string): The path to the directory to be archived. It is walked recursively, and all files and
//     subdirectories are included in the ZIP archive.
//   - archivePath (string): The path where the ZIP archive should be created. The function will resolve this to
//     an absolute path and ensure the archive is saved at that location.
//
// Returns:
//   - error: If an error occurs during any part of the zipping process (such as file opening, writing, or walking),
//     an error is returned.
func ZipIt(dirPath, archivePath string) error {
	archivePath, err := filepath.Abs(archivePath)
	if err != nil {
		return err
	}

	// 'x.y.z' or '.last_version' folder inside the archive
	subdir := filepath.Base(dirPath)

	archivePath = filepath.Clean(archivePath)
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}

	defer helpers.Cleanup(zipFile, nil)

	zipWriter := zip.NewWriter(zipFile)
	defer helpers.Cleanup(zipWriter, nil)

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		filePath = filepath.Clean(filePath)
		if err != nil {
			return err
		}

		if filePath == dirPath {
			return nil
		}

		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(subdir + "/" + relPath)

		if info.IsDir() {
			_, err := zipWriter.Create(relPath + "/")
			return err
		}

		fileInArchive, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		srcFile, err := os.Open(filePath)
		if err != nil {
			return err
		}

		_, err = io.Copy(fileInArchive, srcFile)
		helpers.Cleanup(srcFile, nil)

		return err
	})
	return err
}

// isConvertable checks whether a given file path is eligible for conversion based on its extension.
//
// It returns true if the file has a ".php" extension or if it ends with "description.ru".
// Otherwise, it returns false.
//
// Parameters:
//   - path (string): The file path to check for conversion eligibility.
//
// Returns:
//   - bool: Returns true if the file path ends with ".php" or "description.ru", otherwise returns false.
func isConvertable(path string) bool {
	if path == "" {
		return false
	}

	return (strings.Contains(path, "/lang/") && strings.HasSuffix(path, ".php")) ||
		strings.HasSuffix(path, "description.ru")
}

// IsEmptyDir checks whether the specified directory exists and is empty.
//
// Parameters:
//   - path: The path to the directory to check.
//
// Returns:
//   - true if the directory exists and is empty.
//   - false if the directory does not exist, is not accessible, or contains at least one entry.
//
// Notes:
//   - If the directory cannot be opened (e.g., due to permission issues), the function returns false.
//   - If an error occurs while reading directory entries, it is assumed to be empty
//     (e.g., when the directory does not exist).
//   - Logs an error if closing the directory fails.
func IsEmptyDir(path string) bool {
	path = filepath.Clean(path)
	dir, err := os.Open(path)
	if err != nil {
		return false
	}
	defer helpers.Cleanup(dir, nil)

	entries, err := dir.Readdirnames(1)
	if err != nil {
		return true
	}

	return len(entries) == 0
}

// RemoveEmptyDirs recursively removes empty directories within the specified root directory.
//
// Parameters:
//   - root: The path of the directory to start the Cleanup.
//
// Returns:
//   - A boolean indicating whether the directory itself is empty after processing.
//   - An error if any issue occurs while reading or removing directories.
//
// Behavior:
//   - Traverses all subdirectories recursively.
//   - If a subdirectory becomes empty after processing its contents, it gets removed.
//   - If any error occurs (e.g., permission issues), the error is returned.
//   - The function ensures that only empty directories are removed, leaving files untouched.
//
// Example:
//
//	_, err := RemoveEmptyDirs("/path/to/root")
//	if err != nil {
//	    log.Fatalf("Failed to remove empty directories: %v", err)
//	}
func RemoveEmptyDirs(root string) (bool, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return false, fmt.Errorf("error reading directory %s: %w", root, err)
	}

	empty := true
	for _, entry := range entries {
		if entry.IsDir() {
			inner := filepath.Join(root, entry.Name())
			isEmpty, err := RemoveEmptyDirs(inner)
			if err != nil {
				return false, fmt.Errorf("error removing empty directory %s: %w", inner, err)
			}

			if isEmpty {
				if err := os.Remove(inner); err != nil {
					return false, fmt.Errorf("failed to remove empty directory %s: %w", inner, err)
				}
			} else {
				empty = false
			}
		} else {
			empty = false
		}
	}

	return empty, nil
}

func IsFileExists(path string) (bool, int64) {
	path = filepath.Clean(path)
	file, err := os.Stat(path)
	if err != nil {
		return false, 0
	}

	return file.Mode().IsRegular(), file.Size()
}
