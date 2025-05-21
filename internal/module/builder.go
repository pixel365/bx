package module

import (
	"context"
	"fmt"
	"os"

	"github.com/pixel365/bx/internal/interfaces"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/fs"
	"github.com/pixel365/bx/internal/helpers"
)

type ModuleBuilder struct {
	log    interfaces.Logger
	module *Module
}

func NewModuleBuilder(m *Module, logger interfaces.Logger) interfaces.Builder {
	return &ModuleBuilder{
		log:    logger,
		module: m,
	}
}

// Build orchestrates the entire build process for the module.
// It logs the progress of each phase, such as preparation, collection, and Cleanup.
// If any of these phases fails, the build will be rolled back to ensure a clean state.
//
// The method returns an error if any of the steps (Prepare, Collect, or Cleanup) fail.
func (m *ModuleBuilder) Build(ctx context.Context) error {
	if m.module == nil {
		return errors.ErrNilModule
	}

	if err := helpers.CheckContext(ctx); err != nil {
		return err
	}

	m.log.Info("Building module")

	if err := m.Prepare(); err != nil {
		m.log.Error("Failed to prepare build", err)
		if rollbackErr := m.Rollback(); rollbackErr != nil {
			m.log.Error("Failed to rollback", rollbackErr)
		}
		return err
	}

	m.log.Info("Prepare complete")

	if err := m.Collect(ctx); err != nil {
		m.log.Error("Failed to collect build", err)
		if rollbackErr := m.Rollback(); rollbackErr != nil {
			m.log.Error("Failed to rollback", rollbackErr)
		}
		return err
	}

	m.log.Info("Build complete")

	return nil
}

// Prepare sets up the environment for the build process.
// It validates the module, checks the stages, and creates the necessary directories for the build output and logs.
// If any validation or directory creation fails, an error will be returned.
//
// The method returns an error if the module is invalid or if directories cannot be created.
func (m *ModuleBuilder) Prepare() error {
	if m.module == nil {
		return errors.ErrNilModule
	}

	if err := CheckStages(m.module); err != nil {
		m.log.Error("Prepare: check stages failed", err)
		return err
	}

	m.log.Info("Check stages complete")

	if m.module.BuildDirectory == "" {
		m.module.BuildDirectory = "./build"
	}

	path, err := fs.MkDir(m.module.BuildDirectory)
	if err != nil {
		m.log.Error("Prepare: failed to make build directory", err)
		return err
	}

	m.log.Info("Build directory complete: %s", path)

	m.module.BuildDirectory = path

	path, err = fs.MkDir(fmt.Sprintf("%s/%s", m.module.BuildDirectory, m.module.GetVersion()))
	if err != nil {
		m.log.Info("Prepare: failed to make build version directory")
		return err
	}

	m.log.Info("Build version directory complete: %s", path)

	return nil
}

// Cleanup removes any temporary files and directories created during the build process.
// It ensures the environment is cleaned up by deleting the version-specific build directory.
func (m *ModuleBuilder) Cleanup() {
	if m.module == nil {
		return
	}

	versionDir, err := makeVersionDirectory(m.module)
	if err != nil {
		m.log.Error("Cleanup: failed to make version dir", err)
		return
	}

	if err := os.RemoveAll(versionDir); err != nil {
		m.log.Error("Cleanup: failed to remove version directory", err)
		return
	}

	m.log.Info("Cleanup complete")
}

// Rollback reverts any changes made during the build process.
// It deletes the generated zip file and version-specific directories created during the build.
// This function ensures that any temporary build files are removed
// and that the environment is restored to its previous state.
//
// The method returns an error if the rollback process fails.
func (m *ModuleBuilder) Rollback() error {
	if m.module == nil {
		return errors.ErrNilModule
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

		m.log.Info("Removed zip file: %s", zipPath)
	}

	versionDir, err := makeVersionDirectory(m.module)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return err
	}

	m.log.Info("Removed version directory: %s", versionDir)
	m.log.Info("Rollback complete")

	return nil
}

func (m *ModuleBuilder) prepareVersionDirectory() (string, error) {
	if m.module == nil {
		return "", errors.ErrNilModule
	}

	versionDirectory, err := makeVersionDirectory(m.module)
	if err != nil {
		return "", err
	}

	return versionDirectory, nil
}

// collectStages executes the appropriate set of build stages for the current module.
//
// It selects either `Builds.Release` or `Builds.LastVersion` based on the `LastVersion` flag,
// and delegates the execution to `HandleStages`.
// Any errors during stage execution are logged
// and returned to the caller.
//
// Parameters:
//   - ctx: context used for cancellation and timeouts.
//
// Returns:
//   - error: the first error encountered during stage processing, or nil on success.
func (m *ModuleBuilder) collectStages(ctx context.Context) error {
	stages := m.module.Builds.Release
	if m.module.LastVersion {
		stages = m.module.Builds.LastVersion
	}

	if err := HandleStages(ctx, stages, m.module, m.log, false); err != nil {
		m.log.Error("Collect: handle stages failed", err)
		return err
	}

	m.log.Info("Collect complete")

	return nil
}

// Collect gathers the necessary files for the build.
// It processes each stage in parallel using goroutines to handle file copying.
// The function creates the necessary directories for each stage and copies files as defined in the stage configuration.
//
// The method returns an error if any stage fails or if there are issues zipping the collected files.
func (m *ModuleBuilder) Collect(ctx context.Context) error {
	versionDirectory, err := m.prepareVersionDirectory()
	if err != nil {
		return err
	}

	if err := m.collectStages(ctx); err != nil {
		return err
	}

	err = makeVersionDescription(m)
	if err != nil {
		m.log.Error("Failed to collect build description", err)
		return err
	}

	err = makeVersionFile(m)
	if err != nil {
		m.log.Error("Failed to create version.php", err)
	}

	_, err = fs.RemoveEmptyDirs(versionDirectory)
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

	if fs.IsEmptyDir(versionDirectory) {
		return errors.ErrNoChanges
	}

	if !m.module.LastVersion {
		ok, size := fs.IsFileExists(versionDirectory + "/description.ru")
		if !ok || size == 0 {
			return errors.ErrDescriptionDoesNotExists
		}
	}

	zipPath, err := makeZipFilePath(m.module)
	if err != nil {
		return err
	}

	if err := fs.ZipIt(versionDirectory, zipPath); err != nil {
		m.log.Error("Failed to zip build", err)
		return err
	}

	m.log.Info("Zip complete")

	return nil
}
