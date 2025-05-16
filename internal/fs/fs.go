package fs

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/helpers"

	"github.com/bmatcuk/doublestar/v4"
)

// PathProcessing walks through the source directory specified in `path.From`
// and submits file copy tasks to the provided `filesCh` channel.
//
// It applies ignore rules, change tracking (if enabled), and pattern-based filtering
// before sending any file to be processed.
// All context cancellations are respected and short-circuit the operation early.
//
// Parameters:
//   - ctx: Context used to control cancellation and timeouts.
//   - filesCh: Channel to which valid file copy tasks are submitted (type: types.Path).
//   - cfg: Module configuration providing access to ignore rules and optional change tracking.
//   - path: Describes the copy task, including source, destination, overwrite mode, and encoding settings.
//   - filterRules: A list of file path patterns to include (e.g., ["*.php", "*.tpl"]).
//
// Returns:
//   - error: If an error occurs during directory walking or context cancellation.
func PathProcessing(
	ctx context.Context,
	filesCh chan<- types.Path,
	cfg interfaces.ModuleConfig,
	path types.Path,
	filterRules []string,
) error {
	if err := helpers.CheckContext(ctx); err != nil {
		return err
	}

	var changes *types.Changes
	if !cfg.IsLastVersion() {
		changes = cfg.GetChanges()
	}

	err := filepath.Walk(path.From, visitor(
		ctx, filesCh, cfg, path, filterRules, changes,
	))

	return err
}

// visitor returns a filepath.WalkFunc that handles file traversal logic,
// and submits eligible files as copy tasks to `filesCh`.
//
// For each file, it applies the following logic:
//   - Skips directories and ignored files (via cfg.GetIgnore()).
//   - Applies file inclusion rules from `filterRules`.
//   - If `cfg.IsLastVersion()` is false, applies change tracking using `cfg.GetChanges()`.
//   - Resolves destination path, ensures the destination directory exists, and emits the copy task.
//
// Parameters:
//   - ctx: Context for early cancellation.
//   - filesCh: Channel to submit resulting copy tasks.
//   - cfg: Module config containing ignore and change tracking logic.
//   - path: Copy task parameters (source, target, mode).
//   - filterRules: File pattern filters.
//   - changes: List of changes.
//
// Returns:
//   - A `filepath.WalkFunc` suitable for use in `filepath.Walk`.
func visitor(
	ctx context.Context,
	filesCh chan<- types.Path,
	cfg interfaces.ModuleConfig,
	path types.Path,
	filterRules []string,
	changes *types.Changes,
) filepath.WalkFunc {
	return func(fromPath string, info os.FileInfo, err error) error {
		if ctxErr := helpers.CheckContext(ctx); ctxErr != nil {
			return ctxErr
		}

		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(path.From, fromPath)
		if err != nil {
			return err
		}

		absFrom, err := filepath.Abs(fmt.Sprintf("%s/%s", path.From, relPath))
		if err != nil {
			return err
		}
		absFrom = filepath.Clean(absFrom)

		if shouldSkip(relPath, cfg.GetIgnore()) || !shouldInclude(absFrom, filterRules) {
			return skip(info)
		}

		isDir, err := helpers.IsDir(absFrom)
		if err != nil {
			return err
		}

		if isDir {
			return nil
		}

		if changes != nil && !changes.IsChangedFile(absFrom) {
			return nil
		}

		absTo, err := filepath.Abs(fmt.Sprintf("%s/%s", path.To, relPath))
		if err != nil {
			return err
		}
		absTo = filepath.Clean(absTo)
		toDir := filepath.Dir(absTo)

		if _, err = MkDir(toDir); err != nil {
			return err
		}

		newPath := types.Path{
			From:           absFrom,
			To:             absTo,
			ActionIfExists: path.ActionIfExists,
			Convert:        path.Convert,
		}

		filesCh <- newPath

		return nil
	}
}

func skip(info os.FileInfo) error {
	if info == nil {
		return nil
	}

	if info.IsDir() {
		return filepath.SkipDir
	}

	return nil
}

// CopyFile copies a file from the source path to the destination path, taking into account the context cancellation,
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
//   - errCh (chan<- error): A channel for reporting errors encountered during the operation.
//   - file (types.Path): Path params.
//
// Returns:
//   - None: Errors are sent to the error channel `errCh`, and no return value is provided.
func CopyFile(
	ctx context.Context,
	errCh chan<- error,
	file types.Path,
) {
	if err := helpers.CheckContext(ctx); err != nil {
		errCh <- err
		return
	}

	fileName := strings.LastIndex(file.From, "/")
	if !strings.HasSuffix(file.To, file.From[fileName:]) {
		file.To = filepath.Join(file.To, file.From[fileName:])
		file.To = filepath.Clean(file.To)
	}

	var existingFile os.FileInfo = nil
	stat, err := os.Stat(file.To)
	if err == nil {
		existingFile = stat
		if file.ActionIfExists == types.Skip {
			return
		}
	}

	in, err := os.Open(file.From)
	if err != nil {
		errCh <- err
		return
	}

	defer helpers.Cleanup(in, errCh)

	info, err := in.Stat()
	if err != nil {
		errCh <- err
		return
	}

	allowWrite := true
	if existingFile != nil && file.ActionIfExists == types.ReplaceIfNewer {
		allowWrite = info.ModTime().After(existingFile.ModTime())
	}

	if !allowWrite {
		return
	}

	out, err := os.OpenFile(file.To, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		errCh <- err
		return
	}

	defer helpers.Cleanup(out, errCh)

	var writer io.Writer
	if file.Convert && isConvertable(file.From) {
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

	err = os.Chtimes(file.To, info.ModTime(), info.ModTime())
	if err != nil {
		errCh <- err
		return
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
//   - bool: Returns `true` if the path matches any of the patterns or an error occurs, indicating the path should be
//     skipped.
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
// If the path does not exist, it will be created with permissions 0750. If the directory already exists,
// it does nothing.
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

// IsFileExists checks whether the given path exists and points to a regular file.
//
// The function first cleans the input path using filepath.Clean. It then checks
// if the file exists and is a regular file (not a directory, symlink, etc.).
//
// Parameters:
//   - path: The filesystem path to check.
//
// Returns:
//   - bool: true if the file exists and is a regular file; false otherwise.
//   - int64: The size of the file in bytes if it exists, or 0 if it does not.
func IsFileExists(path string) (bool, int64) {
	path = filepath.Clean(path)
	file, err := os.Stat(path)
	if err != nil {
		return false, 0
	}

	return file.Mode().IsRegular(), file.Size()
}
