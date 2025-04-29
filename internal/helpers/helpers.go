package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"

	"github.com/charmbracelet/huh"
)

type Cfg string

const RootDir Cfg = "root_dir"

func Choose(items *[]string, value *string, title string) error {
	if items == nil || len(*items) == 0 {
		return errors.NoItemsError
	}

	var options []huh.Option[string]
	for i, item := range *items {
		if item == "" {
			return fmt.Errorf("empty item at index %d", i)
		}

		options = append(options, huh.NewOption(item, item))
	}

	return huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value).
		Run()
}

func CaptureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}

	os.Stdout = stdout
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	return buf.String()
}

func ResultMessage(format string, a ...any) {
	if len(a) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format, a...)
	}
}

func GetModulesDir() (string, error) {
	dirPath, _ := os.Getwd()
	return filepath.Abs(fmt.Sprintf("%s/.bx", dirPath))
}

func DefaultYAML() string {
	return `name: "test"
version: "1.0.0"
account: ""
buildDirectory: "./dist/test"
logDirectory: "./logs/test"

variables:
  structPath: "./examples/structure"
  install: "install"
  bitrix: "{structPath}/bitrix"
  local: "{structPath}/local"

stages:
  - name: "components"
    to: "{install}/components"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/components"
      - "{local}/components"
  - name: "templates"
    to: "{install}/templates"
    actionIfFileExists: "replace"
    from:
      - "{bitrix}/templates"
      - "{local}/templates"
  - name: "rootFiles"
    to: "."
    actionIfFileExists: "replace"
    from:
      - "{structPath}/simple-file.php"
  - name: "testFiles"
    to: "test"
    actionIfFileExists: "replace"
    from:
      - "{structPath}/simple-file.php"
    convertTo1251: false

builds:
  release:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"
  lastVersion:
    - "components"
    - "templates"
    - "rootFiles"
    - "testFiles"

ignore:
  - "**/*.log"
`
}

// CheckPath checks if a given path is valid and exists on the filesystem.
//
// This function cleans the provided path and verifies its validity using the `IsValidPath` function.
// It then checks if the file or directory exists using `os.Stat`.
// If the file or directory does not exist or is not valid, an appropriate error is returned.
//
// Parameters:
//   - path (string): The path to be validated and checked for existence.
//
// Returns:
//   - error: An error if the path is invalid or does not exist, otherwise returns nil.
func CheckPath(path string) error {
	path = filepath.Clean(path)
	if !IsValidPath(path, path) {
		return errors.InvalidFilepathError
	}

	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	return nil
}

// IsDir checks if the given path points to a directory.
//
// This function first checks if the path is valid using the `CheckPath` function.
// Then, it retrieves the file information using `os.Stat` to determine whether the
// given path is a directory or not.
//
// Parameters:
//   - path (string): The path to be checked.
//
// Returns:
//   - bool: `true` if the path is a directory, `false` otherwise.
//   - error: An error if the path is invalid or if there is an issue retrieving the file information.
func IsDir(path string) (bool, error) {
	err := CheckPath(path)
	if err != nil {
		return false, err
	}
	fi, _ := os.Stat(path)

	return fi.Mode().IsDir(), nil
}

// CheckContext checks whether the provided context has been canceled or expired.
//
// This function checks if the context has been canceled or the deadline exceeded.
// If the context is done, it returns an error indicating the cancellation or expiration.
// If the context is still active, it returns nil.
//
// Parameters:
//   - ctx (context.Context): The context to check.
//
// Returns:
//   - error: Returns an error if the context is done (canceled or expired), otherwise nil.
func CheckContext(ctx context.Context) error {
	if ctx == context.TODO() {
		return errors.TODOContextError
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return nil
	}
}

// CheckPaths checks whether the paths in the "From" field of a stage are valid.
//
// This function iterates over the paths in the "From" field of the given stage and checks
// if each path is valid by calling the CheckPath function. If any path is invalid,
// it sends the error to the provided channel and exits early.
//
// Parameters:
//   - stage (Stage): The stage containing the "From" paths to check.
//   - ch (chan<- error): A channel to send any errors encountered during the path checks.
//
// The function does not return any value. If an error occurs during path validation,
// the error is sent to the provided channel.
func CheckPaths(stage types.Stage, ch chan<- error) {
	for _, path := range stage.From {
		err := CheckPath(path)
		if err != nil {
			ch <- err
		}
	}
}

// IsValidPath checks if the given filePath is a valid path relative to the basePath.
//
// This function determines whether the absolute path of filePath is within the
// directory specified by basePath. It ensures that the file path does not contain
// any ".." segments, which would indicate an attempt to traverse up the directory
// structure, and that the filePath is within the basePath directory.
//
// Parameters:
//   - filePath (string): The path to check for validity.
//   - basePath (string): The base directory to check against.
//
// Returns:
//   - bool: Returns true if the filePath is valid (i.e., is within the basePath directory),
//     otherwise returns false.
func IsValidPath(filePath, basePath string) bool {
	absBasePath, _ := filepath.Abs(basePath)
	absFilePath, _ := filepath.Abs(filePath)

	if strings.HasPrefix(absBasePath, "..") {
		return false
	}

	return strings.HasPrefix(absFilePath, absBasePath)
}

// ReplaceVariables recursively replaces variables in the input string with
// their corresponding values from the provided map of variables.
//
// The function supports up to 5 levels of recursion (controlled by the `depth`
// parameter) to allow nested variable replacements. If the depth exceeds 5 or
// if the depth is less than 0, an error will be returned.
//
// Parameters:
//   - input (string): The input string containing variables in the format `{variableName}` to replace.
//   - variables (map[string]string): A map containing variable names as keys and their replacement values as strings.
//   - depth (int): The current recursion depth, which should start from 0. The function supports up to 5 levels of recursion.
//
// Returns:
//   - string: The updated string with variables replaced by their corresponding values, or an error if no replacement could be made.
//   - error: An error is returned if the depth exceeds 5, if the depth is negative, or if the replacement results in an empty string.
func ReplaceVariables(input string, variables map[string]string, depth int) (string, error) {
	if depth < 0 {
		return "", errors.SmallDepthError
	}

	if depth > 5 {
		return "", errors.LargeDepthError
	}

	variableRegex := regexp.MustCompile(`\{([a-zA-Z0-9-_]+)}`)
	updated := variableRegex.ReplaceAllStringFunc(input, func(match string) string {
		key := strings.Trim(match, "{}")
		if value, ok := variables[key]; ok {
			return value
		}

		return ""
	})

	if updated == "" {
		return "", errors.ReplacementError
	}

	if updated == input {
		return updated, nil
	}

	return ReplaceVariables(updated, variables, depth+1)
}

// Cleanup closes the provided resource and handles any errors that occur during closure.
//
// If the resource is nil, the function returns immediately without taking any action.
// If resource is not nil, its Close() method is called.
// If Close() returns an error and the provided
//
// channel ch is not nil, the error is sent to ch.
//
// Parameters:
//   - resource: an object that implements the io.Closer interface, representing the resource to be closed.
//   - ch: a channel for reporting errors.
//     If ch is nil, any errors from resource.Close() are ignored.
func Cleanup(resource io.Closer, ch chan<- error) {
	if resource == nil {
		return
	}

	if err := resource.Close(); err != nil && ch != nil {
		ch <- err
	}
}
