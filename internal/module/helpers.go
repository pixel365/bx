package module

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/pixel365/bx/internal/repo"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/errors"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/pixel365/bx/internal/fs"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/types"
)

var (
	checkPathsFunc = helpers.CheckPaths
	copyFileFunc   = fs.CopyFile
)

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
		return nil, errors.ErrInvalidFilepath
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
		return nil, errors.ErrNilCmd
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
		return nil, errors.ErrInvalidRootDir
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

// HandleStages executes a sequence of stages defined in the provided module.
//
// For each stage name in the `stages` slice, the corresponding stage is resolved from the module `m`
// and processed concurrently via the `handleStage` function.
// Each stage may produce file copy tasks,
// which are sent to a shared channel (`filesCh`) and handled by a pool of worker goroutines.
// Log messages are sent asynchronously to a logging worker via `logCh`.
//
// The function manages synchronization using multiple WaitGroups and coordinates shutdown via
// context cancellation.
// If any error occurs in stage processing or file copying,
// the first error is captured and the context is canceled.
//
// Parameters:
//   - ctx: base context used for cancellation.
//   - stages: list of stage names to execute.
//   - m: the module containing the stage definitions.
//   - logger: implementation of BuildLogger for reporting progress and messages.
//   - customCommandMode: if true, skips versioned output directory creation.
//
// Returns:
//   - err: the first error encountered during execution, or nil if all stages completed successfully.
func HandleStages(
	ctx context.Context,
	stages []string,
	m *Module,
	logger interfaces.BuildLogger,
	customCommandMode bool,
) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if m == nil {
		err = errors.ErrNilModule
		return
	}

	dir := ""
	if !customCommandMode {
		dir, err = makeVersionDirectory(m)
		if err != nil {
			return
		}
	}

	workersCount := runtime.NumCPU() * 10
	filesCh := make(chan types.Path, workersCount)
	errCh := make(chan error, 1)
	logCh := make(chan string, 100)

	var stagesWorkersWg sync.WaitGroup
	var copyFilesWg sync.WaitGroup
	var once sync.Once

	go errorWorker(errCh, cancel, &once, &err)

	go logWorker(logCh, logger)

	go copyWorkers(ctx, &copyFilesWg, filesCh, errCh, workersCount)

	for _, name := range stages {
		stage, stageErr := m.FindStage(name)
		if stageErr != nil {
			return stageErr
		}

		s := stage
		stagesWorkersWg.Add(1)
		go func(stage *types.Stage) {
			defer stagesWorkersWg.Done()
			handleStage(ctx, filesCh, logCh, errCh, m, *stage, dir, m.StageCallback)
		}(&s)
	}

	go cleanupWorker(&stagesWorkersWg, &copyFilesWg, &once, cancel, filesCh, logCh, errCh)

	<-ctx.Done()
	return
}

// cleanupWorker coordinates the shutdown of all channels and workers.
//
// It waits for stage processing (stage worker goroutines) and copy workers to complete,
// then closes all communication channels and invokes cancellation once to signal global termination.
//
// Parameters:
//   - stagesWorkersWg: WaitGroup tracking stage processing goroutines.
//   - copyFilesWg: WaitGroup tracking file copy worker goroutines.
//   - once: sync.Once to ensure cancellation is called only once.
//   - cancel: context.CancelFunc to terminate the parent context.
//   - filesCh: channel of file copy tasks to be closed after stage completion.
//   - logCh: channel of log messages to be closed after copy workers finish.
//   - errCh: channel of errors to be closed after copy workers finish.
func cleanupWorker(
	stagesWorkersWg, copyFilesWg *sync.WaitGroup,
	once *sync.Once,
	cancel context.CancelFunc,
	filesCh chan types.Path,
	logCh chan string,
	errCh chan error,
) {
	stagesWorkersWg.Wait()
	close(filesCh)

	copyFilesWg.Wait()
	close(logCh)
	close(errCh)

	once.Do(cancel)
}

// copyWorkers launches a fixed number of worker goroutines to process file copy tasks.
//
// Each worker reads from the `filesCh` channel and invokes the `copyFileFunc` function
// to handle the file copying.
// Workers exit when the channel is closed.
//
// Parameters:
//   - ctx: context for cancellation.
//   - wg: WaitGroup used to signal completion of all worker goroutines.
//   - filesCh: channel carrying file copy tasks (type Path).
//   - errCh: channel used to report errors from the workers.
//   - workersCount: number of worker goroutines to spawn.
func copyWorkers(ctx context.Context, wg *sync.WaitGroup, filesCh chan types.Path,
	errCh chan<- error, workersCount int) {
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesCh {
				copyFileFunc(ctx, errCh, file)
			}
		}()
	}
}

// errorWorker monitors the error channel and triggers cancellation on the first non-nil error.
//
// It ensures that only the first error is captured and that cancellation is invoked only once
// using sync.Once.
// It is intended to run as a background goroutine.
//
// Parameters:
//   - ch: channel from which errors are read.
//   - cancel: function to cancel the shared context.
//   - once: sync.Once to ensure cancel is called only once.
//   - err: pointer to the shared error variable to capture the first error.
func errorWorker(ch <-chan error, cancel context.CancelFunc, once *sync.Once, err *error) {
	for e := range ch {
		if e != nil {
			once.Do(func() {
				*err = e
				cancel()
			})
		}
	}
}

// logWorker consumes log messages from the provided channel and sends them to the logger.
//
// It runs as a background goroutine, and processes log messages until the channel is closed.
//
// Parameters:
//   - ch: channel carrying log messages.
//   - logger: implementation of BuildLogger used to output logs.
func logWorker(ch <-chan string, logger interfaces.BuildLogger) {
	for msg := range ch {
		logger.Info(msg)
	}
}

func makeVersionDirectory(module *Module) (string, error) {
	if module == nil || module.BuildDirectory == "" {
		return "", errors.ErrNilModule
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
		return errors.ErrNilModule
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

// handleStage processes a single build stage concurrently and sends copy tasks to the shared file worker pool.
//
// The function performs the following steps:
//  1. Logs the start and completion of the stage via `logCh`.
//  2. Resolves a `Runnable` using the provided callback `cb` and executes its `PreRun` hook.
//  3. Validates the context and resolves the target directory for the stage output.
//  4. Create the target directory (recursively if needed).
//  5. For each input path in `stage.From`, spawns a goroutine that validates the context,
//     builds a copy `Path` struct, and sends it to `filesCh` to be processed by `copyWorkers`.
//  6. Wait for all spawned copy goroutines to finish.
//  7. If the runner was initialized, execute the `PostRun` hook.
//
// Errors from directory creation, pre/post run hooks, or file copy failures are sent to `errCh`.
// If the context is canceled at any point, the function exits early.
//
// Parameters:
//   - ctx: the context used for cancellation and timeout propagation.
//   - filesCh: a shared channel for submitting file copy tasks.
//   - logCh: a channel for emitting informational log messages.
//   - errCh: a channel for reporting execution errors.
//   - module: the current build module used in copy context.
//   - stage: the stage definition to process.
//   - rootDir: optional root output directory prefix; if empty, use `stage.To` as-is.
//   - cb: a callback that returns a `Runnable` with `PreRun` and `PostRun` methods for the stage.
func handleStage(
	ctx context.Context,
	filesCh chan<- types.Path,
	logCh chan<- string,
	errCh chan<- error,
	module *Module,
	stage types.Stage,
	rootDir string,
	cb func(string) (interfaces.Runnable, error),
) {
	logCh <- fmt.Sprintf("Handling stage %s", stage.Name)
	defer func() {
		logCh <- fmt.Sprintf("Finished stage %s", stage.Name)
	}()

	runner, cbErr := cb(stage.Name)

	if cbErr == nil {
		if err := runner.PreRun(ctx); err != nil {
			errCh <- fmt.Errorf("pre-run callback failed for stage %s: %w", stage.Name, err)
			return
		}
	}

	if err := helpers.CheckContext(ctx); err != nil {
		return
	}

	dirPath := stage.To
	if rootDir != "" {
		dirPath = filepath.Join(rootDir, stage.To)
	}

	dirPath = filepath.Clean(dirPath)
	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		errCh <- fmt.Errorf("failed to get absolute path for stage %s: %s", stage.Name, err)
		return
	}

	to, err := fs.MkDir(dirPath)
	if err != nil {
		errCh <- fmt.Errorf("failed to make stage `to` directory: %w", err)
		return
	}

	var wg sync.WaitGroup
	for _, from := range stage.From {
		if err = helpers.CheckContext(ctx); err != nil {
			return
		}

		fromCopy := from
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err = helpers.CheckContext(ctx); err != nil {
				return
			}

			path := types.Path{
				From:           fromCopy,
				To:             to,
				ActionIfExists: stage.ActionIfFileExists,
				Convert:        stage.ConvertTo1251,
			}

			if err := fs.PathProcessing(
				ctx,
				filesCh,
				module,
				path,
				stage.Filter,
			); err != nil {
				errCh <- fmt.Errorf("failed to copy from %s to %s: %w", path.From, path.To, err)
			}
		}()
	}

	wg.Wait()

	if runner != nil {
		if err = runner.PostRun(ctx); err != nil {
			errCh <- fmt.Errorf("post-run callback failed for stage %s: %w", stage.Name, err)
		}
	}
}

func makeVersionDescription(builder *ModuleBuilder) error {
	// If the full latest version is being built, then the version description file is not needed.
	// However, it may be present when copying if specified in the configuration, at the discretion of the developer.
	if builder.module.LastVersion {
		return nil
	}

	descriptionRu := strings.Builder{}
	encoder := charmap.Windows1251.NewEncoder()

	if builder.module.Description != "" {
		encodedDescriptionRu, err := encoder.String(builder.module.Description + "\n")
		if err != nil {
			return fmt.Errorf("encoding description [%s]: %w", builder.module.Description, err)
		}

		_, _ = descriptionRu.WriteString(encodedDescriptionRu)
	} else {
		if builder.module.Repository == "" {
			return nil
		}

		commits, err := repo.ChangelogList(builder.module.Repository, builder.module.Changelog)
		if err != nil {
			return err
		}

		if len(commits) == 0 {
			return nil
		}

		for _, commit := range commits {
			encodedLine, err := encoder.String(commit + "\n")
			if err != nil {
				return fmt.Errorf("encoding commit [%s]: %w", commit, err)
			}
			_, _ = descriptionRu.WriteString(encodedLine)
		}
	}

	footer, err := builder.module.Changelog.EncodedFooter()
	if err != nil {
		return fmt.Errorf(
			"encoding footer template [%s]: %w",
			builder.module.Changelog.FooterTemplate,
			err,
		)
	}
	_, _ = descriptionRu.WriteString(footer)

	err = writeFileForVersion(builder, "description.ru", descriptionRu.String())
	if err != nil {
		return fmt.Errorf("failed to make description file: %w", err)
	}

	return nil
}

func makeVersionFile(builder *ModuleBuilder) error {
	if builder.module.LastVersion {
		return nil
	}

	now := time.Now().Format(time.DateTime)
	buf := strings.Builder{}
	buf.WriteString("<?php\n")
	buf.WriteString("$arModuleVersion = array(\n")
	buf.WriteString("\t\t\"VERSION\" => \"" + builder.module.Version + "\"\n")
	buf.WriteString("\t\t\"VERSION_DATE\" => \"" + now + "\"\n")
	buf.WriteString(");\n")

	err := writeFileForVersion(builder, "/install/version.php", buf.String())
	if err != nil {
		return fmt.Errorf("failed to make version.php file: %w", err)
	}

	return nil
}
