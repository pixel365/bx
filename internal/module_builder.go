package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/text/encoding/charmap"
)

type ModuleBuilder struct {
	logger BuildLogger
	module *Module
}

func NewModuleBuilder(m *Module, logger BuildLogger) Builder {
	return &ModuleBuilder{
		logger: logger,
		module: m,
	}
}

// Build orchestrates the entire build process for the module.
// It logs the progress of each phase, such as preparation, collection, and Cleanup.
// If any of these phases fails, the build will be rolled back to ensure a clean state.
//
// The method returns an error if any of the steps (Prepare, Collect, or Cleanup) fail.
func (m *ModuleBuilder) Build() error {
	if m.module == nil {
		return NilModuleError
	}

	if err := CheckContext(m.module.Ctx); err != nil {
		return err
	}

	m.logger.Info("Building module")

	if err := m.Prepare(); err != nil {
		m.logger.Error("Failed to prepare build", err)
		if rollbackErr := m.Rollback(); rollbackErr != nil {
			m.logger.Error("Failed to rollback", rollbackErr)
		}
		return err
	}

	m.logger.Info("Prepare complete")

	if err := m.Collect(); err != nil {
		m.logger.Error("Failed to collect build", err)
		if rollbackErr := m.Rollback(); rollbackErr != nil {
			m.logger.Error("Failed to rollback", rollbackErr)
		}
		return err
	}

	m.logger.Info("Build complete")

	return nil
}

// Prepare sets up the environment for the build process.
// It validates the module, checks the stages, and creates the necessary directories for the build output and logs.
// If any validation or directory creation fails, an error will be returned.
//
// The method returns an error if the module is invalid or if directories cannot be created.
func (m *ModuleBuilder) Prepare() error {
	if m.module == nil {
		return NilModuleError
	}

	if err := m.module.IsValid(); err != nil {
		m.logger.Error("Prepare: module is invalid", err)
		return err
	}

	m.logger.Info("Validation complete")

	if err := CheckStages(m.module); err != nil {
		m.logger.Error("Prepare: check stages failed", err)
		return err
	}

	m.logger.Info("Check stages complete")

	if m.module.BuildDirectory == "" {
		m.module.BuildDirectory = "./build"
	}

	if m.module.LogDirectory == "" {
		m.module.LogDirectory = "./log"
	}

	path, err := mkdir(m.module.BuildDirectory)
	if err != nil {
		m.logger.Error("Prepare: failed to make build directory", err)
		return err
	}

	m.logger.Info("Build directory complete: %s", path)

	m.module.BuildDirectory = path

	path, err = mkdir(m.module.LogDirectory)
	if err != nil {
		m.logger.Error("Prepare: failed to make log directory", err)
		return err
	}

	m.logger.Info("Log directory complete: %s", path)

	m.module.LogDirectory = path

	path, err = mkdir(fmt.Sprintf("%s/%s", m.module.BuildDirectory, m.module.GetVersion()))
	if err != nil {
		m.logger.Info("Prepare: failed to make build version directory")
		return err
	}

	m.logger.Info("Build version directory complete: %s", path)

	return nil
}

// Cleanup removes any temporary files and directories created during the build process.
// It ensures the environment is cleaned up by deleting the version-specific build directory.
func (m *ModuleBuilder) Cleanup() {
	if m.module == nil {
		return
	}

	defer m.logger.Cleanup()

	versionDir, err := makeVersionDirectory(m.module)
	if err != nil {
		m.logger.Error("Cleanup: failed to make version dir", err)
		return
	}

	if err := os.RemoveAll(versionDir); err != nil {
		m.logger.Error("Cleanup: failed to remove version directory", err)
		return
	}

	m.logger.Info("Cleanup complete")
}

// Rollback reverts any changes made during the build process.
// It deletes the generated zip file and version-specific directories created during the build.
// This function ensures that any temporary build files are removed
// and that the environment is restored to its previous state.
//
// The method returns an error if the rollback process fails.
func (m *ModuleBuilder) Rollback() error {
	if m.module == nil {
		return NilModuleError
	}

	zipPath, err := makeZipFilePath(m.module)
	if err != nil {
		return err
	}

	if zipStat, err := os.Stat(zipPath); err == nil && !zipStat.IsDir() {
		err := os.Remove(zipPath)
		if err != nil {
			return err
		}

		m.logger.Info("Removed zip file: %s", zipPath)
	}

	versionDir, err := makeVersionDirectory(m.module)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return err
	}

	m.logger.Info("Removed version directory: %s", versionDir)
	m.logger.Info("Rollback complete")

	return nil
}

// Collect gathers the necessary files for the build.
// It processes each stage in parallel using goroutines to handle file copying.
// The function creates the necessary directories for each stage and copies files as defined in the stage configuration.
//
// The method returns an error if any stage fails or if there are issues zipping the collected files.
func (m *ModuleBuilder) Collect() error {
	if m.module == nil {
		return NilModuleError
	}

	versionDirectory, err := makeVersionDirectory(m.module)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(m.module.Stages))

	stages := m.module.Builds.Release
	if m.module.LastVersion {
		stages = m.module.Builds.LastVersion
	}

	if err := HandleStages(stages, m.module, &wg, errCh, m.logger, false); err != nil {
		m.logger.Error("Collect: handle stages failed", err)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		m.logger.Error("Failed to collect build", nil)
		return fmt.Errorf("errors: %v", errs)
	}

	m.logger.Info("Collect complete")

	err = makeVersionDescription(m)
	if err != nil {
		m.logger.Error("Failed to collect build description", err)
		return err
	}

	_, err = removeEmptyDirs(versionDirectory)
	if err != nil {
		return err
	}

	_, err = os.Stat(versionDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("version directory does not exist: %s", versionDirectory)
		}
		return err
	}

	if isEmptyDir(versionDirectory) {
		return NoChangesError
	}

	zipPath, err := makeZipFilePath(m.module)
	if err != nil {
		return err
	}

	if err := zipIt(versionDirectory, zipPath); err != nil {
		m.logger.Error("Failed to zip build", err)
		return err
	}

	m.logger.Info("Zip complete")

	return nil
}

// handleStage processes an individual stage during the build.
// It manages file copying based on the configuration for each stage, including handling concurrency using goroutines.
// For each stage, it creates the necessary directories and copies files from the source to the destination directory.
//
// Parameters:
//   - ctx: The context used to manage cancellation or timeouts.
//   - wg: The wait group to synchronize the completion of all goroutines.
//   - errCh: A channel for capturing any errors that occur during the process.
//   - logger: The logger used to log messages about the process.
//   - ignore: A list of files or directories to be ignored during file copying.
//   - stage: The specific stage being processed, which contains source and destination paths.
//   - rootDir: The directory where the files will be placed.
//   - callback: Stage Callback
//
// Returns:
//   - None.
//     Errors will be passed to the `errCh` channel for further handling.
func handleStage(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	logger BuildLogger,
	module *Module,
	stage Stage,
	rootDir string,
	cb func(string) (Runnable, error),
) {
	callback, cbErr := cb(stage.Name)

	defer func() {
		if cbErr == nil {
			wg.Add(1)
			go callback.PostRun(ctx, wg, logger)
		}
		wg.Done()
	}()

	if cbErr == nil {
		wg.Add(1)
		go callback.PreRun(ctx, wg, logger)
	}

	var err error

	if logger != nil {
		logger.Info("Handling stage %s", stage.Name)
	}

	defer func() {
		if err != nil {
			if logger != nil {
				logger.Error(fmt.Sprintf("Failed to handle stage %s: %s", stage.Name, err), err)
			}
			errCh <- err
		} else {
			if logger != nil {
				logger.Info("Finished stage %s", stage.Name)
			}
		}
	}()

	if err := CheckContext(ctx); err != nil {
		return
	}

	dirPath := stage.To
	if rootDir != "" {
		dirPath = filepath.Join(rootDir, stage.To)
	}

	dirPath = filepath.Clean(dirPath)
	dirPath, err = filepath.Abs(dirPath)
	if err != nil {
		if logger != nil {
			logger.Error(
				fmt.Sprintf("Failed to get absolute path for stage %s: %s", stage.Name, err),
				err,
			)
		}
		return
	}

	to, err := mkdir(dirPath)
	if err != nil {
		if logger != nil {
			logger.Error("Failed to make stage `to` directory", err)
		}
		return
	}

	for _, from := range stage.From {
		if err := CheckContext(ctx); err != nil {
			return
		}

		wg.Add(1)
		go copyFromPath(
			ctx,
			wg,
			errCh,
			module,
			from,
			to,
			stage.ActionIfFileExists,
			stage.ConvertTo1251,
			stage.Filter,
		)
	}
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

func makeVersionDirectory(module *Module) (string, error) {
	if module == nil || module.BuildDirectory == "" {
		return "", NilModuleError
	}

	path := filepath.Join(module.BuildDirectory, module.GetVersion())
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}

func makeVersionDescription(builder *ModuleBuilder) error {
	// If the full latest version is being built, then the version description file is not needed.
	// However, it may be present when copying if specified in the configuration, at the discretion of the developer.
	if builder.module.LastVersion {
		return nil
	}

	versionDir, err := makeVersionDirectory(builder.module)
	if err != nil {
		return err
	}

	descriptionRu := ""

	encoder := charmap.Windows1251.NewEncoder()

	if builder.module.Description != "" {

		encodedDescriptionRu, err := encoder.String(builder.module.Description + "\n")
		if err != nil {
			return fmt.Errorf("encoding description [%s]: %w", builder.module.Description, err)
		}
		descriptionRu = encodedDescriptionRu

	} else {

		if builder.module.Repository == "" {
			return nil
		}

		commits, err := ChangelogList(builder.module.Repository, builder.module.Changelog)
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
			descriptionRu += encodedLine
		}

	}

	descriptionPath := filepath.Join(versionDir, "description.ru")
	descriptionPath = filepath.Clean(descriptionPath)

	file, err := os.Create(descriptionPath)
	if err != nil {
		return fmt.Errorf(
			"failed to create description file [%s]: %w",
			descriptionPath,
			err,
		)
	}

	defer func() {
		if err := file.Close(); err != nil && builder.logger != nil {
			builder.logger.Error("Failed to close description file", err)
		}
	}()

	_, err = file.WriteString(descriptionRu)
	if err != nil {
		return fmt.Errorf("failed to write description to %s: %w", descriptionPath, err)
	}

	return nil
}
