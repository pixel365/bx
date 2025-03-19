package internal

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/huh"
)

type Cfg string

const RootDir Cfg = "root_dir"

func Choose(items *[]string, value *string, title string) error {
	if len(*items) == 0 {
		switch any(items).(type) {
		default:
			return NoItemsError
		}
	}

	var options []huh.Option[string]
	for _, item := range *items {
		options = append(options, huh.NewOption(item, item))
	}

	if err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value).
		Run(); err != nil {
		return err
	}

	return nil
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

func GetModulesDir(path string) (string, error) {
	var err error
	dirPath := path
	if dirPath == "" {
		dirPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

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

// AllModules returns a list of module names found in the specified directory.
//
// The function reads the directory, checks for files (skipping directories), and attempts to read
// each file as a module using the ReadModule function. If a file can be successfully read as a
// module, its name is added to the list.
//
// Parameters:
//   - directory (string): The path to the directory to scan for modules.
//
// Returns:
//   - *[]string: A pointer to a slice of strings containing the names of all successfully read modules.
func AllModules(directory string) *[]string {
	var modules []string

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		module, err := ReadModule(filePath, "", true)
		if err != nil {
			continue
		}

		modules = append(modules, module.Name)
	}

	return &modules
}

// ReadModule reads a module from a YAML file or directory path and returns a Module object.
//
// This function attempts to read a module from the specified path. If the `file` flag is true,
// the function treats `path` as the file path directly. Otherwise, it expects a YAML file with
// the name of the module, combining the `path` and `name` parameters to form the file path.
//
// Parameters:
//   - path (string): The directory or file path where the module file is located.
//   - name (string): The name of the module. Used to construct the file path when `file` is false.
//   - file (bool): Flag indicating whether the `path` is a direct file path or a directory where
//     a module file should be looked for.
//
// Returns:
//   - *Module: A pointer to a `Module` object if the file can be successfully read and unmarshalled.
//   - error: An error if reading or unmarshalling the file fails.
func ReadModule(path, name string, file bool) (*Module, error) {
	var filePath string
	var err error

	if !file {
		filePath, err = filepath.Abs(path + "/" + name + ".yaml")
	} else {
		filePath, err = filepath.Abs(path)
	}

	if err != nil {
		return nil, err
	}

	if !isValidPath(filePath, path) {
		return nil, InvalidFilepathError
	}

	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	var module Module
	if err := yaml.Unmarshal(data, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

// CheckPath checks if a given path is valid and exists on the filesystem.
//
// This function cleans the provided path and verifies its validity using the `isValidPath` function.
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
	if !isValidPath(path, path) {
		return fmt.Errorf("invalid path: %s", path)
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
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.Mode().IsDir(), nil
}

// CheckStages validates the paths in the stages of the given module.
//
// This function iterates over the stages in the provided module and concurrently
// checks the paths defined in each stage using goroutines. If any errors are encountered,
// they are collected in a channel and returned as a combined error.
//
// Parameters:
//   - module (*Module): The module containing stages to be validated.
//
// Returns:
//   - error: Returns an error if any validation fails in any stage's paths.
//     If no errors are found, it returns nil.
func CheckStages(module *Module) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(module.Stages)*5)

	for _, item := range module.Stages {
		wg.Add(1)
		go func(wg *sync.WaitGroup, item Stage) {
			defer wg.Done()
			checkPaths(item, errCh)
		}(&wg, item)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors: %v", errs)
	}

	return nil
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
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return nil
	}
}

// checkPaths checks whether the paths in the "From" field of a stage are valid.
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
func checkPaths(stage Stage, ch chan<- error) {
	for _, path := range stage.From {
		err := CheckPath(path)
		if err != nil {
			ch <- err
			return
		}
	}
}

// isValidPath checks if the given filePath is a valid path relative to the basePath.
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
func isValidPath(filePath, basePath string) bool {
	absBasePath, _ := filepath.Abs(basePath)
	absFilePath, _ := filepath.Abs(filePath)

	if strings.Contains(absFilePath, "..") {
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
		return "", SmallDepthError
	}

	if depth > 5 {
		return "", LargeDepthError
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
		return "", ReplacementError
	}

	if updated == input {
		return updated, nil
	}

	return ReplaceVariables(updated, variables, depth+1)
}

func ReadModuleFromFlags(cmd *cobra.Command) (*Module, error) {
	path := cmd.Context().Value(RootDir).(string)
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return nil, err
	}

	file, err := cmd.Flags().GetString("file")
	file = strings.TrimSpace(file)
	if err != nil {
		return nil, err
	}

	isFile := len(file) > 0

	if !isFile && name == "" {
		err := Choose(AllModules(path), &name, "")
		if err != nil {
			return nil, err
		}
	}

	if isFile {
		path = file
	}

	module, err := ReadModule(path, name, isFile)
	if err != nil {
		return nil, err
	}

	module.Ctx = cmd.Context()

	return module, nil
}

func HandleStages(
	stages []string,
	m *Module,
	wg *sync.WaitGroup,
	errCh chan<- error,
	logger BuildLogger,
	customCommandMode bool,
) error {
	var err error
	dir := ""

	if !customCommandMode {
		dir, err = makeVersionDirectory(m)
		if err != nil {
			return err
		}
	}

	for _, stageName := range stages {
		if err := CheckContext(m.Ctx); err != nil {
			return err
		}

		stage, err := m.FindStage(stageName)
		if err != nil {
			return fmt.Errorf("failed to find stage: %w", err)
		}

		wg.Add(1)

		go handleStage(m.Ctx, wg, errCh, logger, m, stage, dir, m.StageCallback)
	}

	return nil
}
