package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/errors"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/pixel365/bx/internal/fs"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/types"
)

var checkPathsFunc = helpers.CheckPaths

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

	if !helpers.IsValidPath(filePath, path) {
		return nil, errors.InvalidFilepathError
	}

	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	var m Module
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func ReadModuleFromFlags(cmd *cobra.Command) (*Module, error) {
	if cmd == nil {
		return nil, errors.NilCmdError
	}

	name, _ := cmd.Flags().GetString("name")
	file, _ := cmd.Flags().GetString("file")
	repository, _ := cmd.Flags().GetString("repository")
	description, _ := cmd.Flags().GetString("description")
	version, _ := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)

	file = strings.TrimSpace(file)
	isFile := len(file) > 0

	path, ok := cmd.Context().Value(helpers.RootDir).(string)
	if !ok {
		return nil, errors.InvalidRootDirError
	}

	if !isFile && name == "" {
		err := helpers.Choose(AllModules(path), &name, "")
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

	if repository != "" {
		module.Repository = repository
	}

	if version != "" {
		module.Version = version
	}

	if description != "" {
		module.Description = description
	}

	return module, module.IsValid()
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
		return nil
	}

	for _, file := range files {
		if !file.IsDir() {

			filePath := filepath.Join(directory, file.Name())
			module, err := ReadModule(filePath, "", true)
			if err != nil {
				continue
			}

			modules = append(modules, module.Name)
		}
	}

	return &modules
}

func HandleStages(
	stages []string,
	m *Module,
	wg *sync.WaitGroup,
	errCh chan<- error,
	logger interfaces.BuildLogger,
	customCommandMode bool,
) error {
	if m == nil {
		return errors.NilModuleError
	}

	dir := ""

	var err error
	if !customCommandMode {
		dir, err = makeVersionDirectory(m)
		if err != nil {
			return err
		}
	}

	for _, stageName := range stages {
		if err := helpers.CheckContext(m.Ctx); err != nil {
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

func makeVersionDirectory(module *Module) (string, error) {
	if module == nil || module.BuildDirectory == "" {
		return "", errors.NilModuleError
	}

	path := filepath.Join(module.BuildDirectory, module.GetVersion())
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}

func makeZipFilePath(module *Module) (string, error) {
	path := filepath.Join(module.BuildDirectory, fmt.Sprintf("%s.zip", module.GetVersion()))
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}

func writeFileForVersion(builder *ModuleBuilder, path, content string) error {
	if len(content) == 0 {
		return nil
	}

	versionDir, err := makeVersionDirectory(builder.module)
	if err != nil {
		return err
	}

	fp := filepath.Join(versionDir, path)
	fp = filepath.Clean(fp)

	dirs := strings.Split(fp, "/")
	dirPath := strings.Join(dirs[:len(dirs)-1], "/")
	_, err = fs.MkDir(dirPath)
	if err != nil {
		return err
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil && builder.logger != nil {
			builder.logger.Error("Failed to close "+path, err)
		}
	}()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
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
	if module == nil {
		return errors.NilModuleError
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(module.Stages)*5)

	for _, item := range module.Stages {
		wg.Add(1)
		go func(wg *sync.WaitGroup, item types.Stage) {
			defer wg.Done()
			checkPathsFunc(item, errCh)
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
