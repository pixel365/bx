package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Build orchestrates the entire build process for the module.
// It logs the progress of each phase, such as preparation, collection, and cleanup.
// If any of these phases fails, the build will be rolled back to ensure a clean state.
//
// The method returns an error if any of the steps (Prepare, Collect, or Cleanup) fail.
func (m *Module) Build() error {
	if err := CheckContext(m.Ctx); err != nil {
		return err
	}

	logFile, err := os.OpenFile(
		fmt.Sprintf("./%s-%s.%s.log", m.Name, m.Version, time.Now().UTC().Format(time.RFC3339)),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)
	if err != nil {
		return err
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			slog.Error(err.Error())
		} else {
			path := fmt.Sprintf("%s/%s", m.LogDirectory, logFile.Name())
			path = filepath.Clean(path)
			err := os.Rename(logFile.Name(), path)
			if err != nil {
				slog.Error(err.Error())
			}
		}
	}(logFile)

	log := zerolog.New(logFile).With().Timestamp().Logger()
	log.Info().Msg("Building module")

	if err := m.Prepare(&log); err != nil {
		log.Error().Err(err).Msg("Failed to prepare build")
		return m.Rollback(&log)
	}

	log.Info().Msg("Prepare complete")

	if err := m.Collect(&log); err != nil {
		log.Error().Err(err).Msg("Failed to collect build")
		return m.Rollback(&log)
	}

	log.Info().Msg("Build complete")

	if err := m.Cleanup(&log); err != nil {
		log.Error().Err(err).Msg("Failed to cleanup build")
		return err
	}

	log.Info().Msg("Cleanup complete")

	return nil
}

// Prepare sets up the environment for the build process.
// It validates the module, checks the stages, and creates the necessary directories for the build output and logs.
// If any validation or directory creation fails, an error will be returned.
//
// The method returns an error if the module is invalid or if directories cannot be created.
func (m *Module) Prepare(log *zerolog.Logger) error {
	if err := m.IsValid(); err != nil {
		log.Error().Err(err).Msg("Prepare: module is invalid")
		return err
	}

	log.Info().Msg("Validation complete")

	if err := CheckStages(m); err != nil {
		log.Error().Err(err).Msg("Prepare: check stages failed")
		return err
	}

	log.Info().Msg("Check stages complete")

	if m.BuildDirectory == "" {
		m.BuildDirectory = "./build"
	}

	if m.LogDirectory == "" {
		m.LogDirectory = "./log"
	}

	path, err := mkdir(m.BuildDirectory)
	if err != nil {
		log.Error().Err(err).Msg("Prepare: failed to make build directory")
		return err
	}

	log.Info().Msgf("Build directory complete: %s", path)

	m.BuildDirectory = path

	path, err = mkdir(m.LogDirectory)
	if err != nil {
		log.Error().Err(err).Msg("Prepare: failed to make log directory")
		return err
	}

	log.Info().Msgf("Log directory complete: %s", path)

	m.LogDirectory = path

	path, err = mkdir(fmt.Sprintf("%s/%s", m.BuildDirectory, m.Version))
	if err != nil {
		log.Error().Err(err).Msg("Prepare: failed to make build version directory")
		return err
	}

	log.Info().Msgf("Build version directory complete: %s", path)

	return nil
}

// Cleanup removes any temporary files and directories created during the build process.
// It ensures the environment is cleaned up by deleting the version-specific build directory.
//
// The method returns an error if the cleanup process fails.
func (m *Module) Cleanup(log *zerolog.Logger) error {
	versionDir, err := makeVersionDirectory(m)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return err
	}

	log.Info().Msg("Cleanup complete")

	return nil
}

// Rollback reverts any changes made during the build process.
// It deletes the generated zip file and version-specific directories created during the build.
// This function ensures that any temporary build files are removed
// and that the environment is restored to its previous state.
//
// The method returns an error if the rollback process fails.
func (m *Module) Rollback(log *zerolog.Logger) error {
	zipPath, err := makeZipFilePath(m)
	if err != nil {
		return err
	}

	if zipStat, err := os.Stat(zipPath); err == nil && !zipStat.IsDir() {
		err := os.Remove(zipPath)
		if err != nil {
			return err
		}

		log.Info().Msgf("Removed zip file: %s", zipPath)
	}

	versionDir, err := makeVersionDirectory(m)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(versionDir); err != nil {
		return err
	}

	log.Info().Msgf("Removed version directory: %s", versionDir)
	log.Info().Msg("Rollback complete")

	return nil
}

// Collect gathers the necessary files for the build.
// It processes each stage in parallel using goroutines to handle file copying.
// The function creates the necessary directories for each stage and copies files as defined in the stage configuration.
//
// The method returns an error if any stage fails or if there are issues zipping the collected files.
func (m *Module) Collect(log *zerolog.Logger) error {
	versionDirectory, err := makeVersionDirectory(m)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(m.Stages))

	for _, item := range m.Stages {
		if err := CheckContext(m.Ctx); err != nil {
			return err
		}

		wg.Add(1)
		go handleStage(m.Ctx, &wg, errCh, log, &m.Ignore, item, versionDirectory)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		log.Error().Int("errors", len(errs)).Msg("Failed to collect build")
		return fmt.Errorf("errors: %v", errs)
	}

	log.Info().Msg("Collect complete")

	zipPath, err := makeZipFilePath(m)
	if err != nil {
		return err
	}

	if err := zipIt(versionDirectory, zipPath); err != nil {
		log.Error().Err(err).Msg("Failed to zip build")
		return err
	}

	log.Info().Msg("Zip complete")

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
//   - log: The logger used to log messages about the process.
//   - ignore: A list of files or directories to be ignored during file copying.
//   - stage: The specific stage being processed, which contains source and destination paths.
//   - buildDirectory: The directory where the build files will be placed.
//
// Returns:
//   - None.
//     Errors will be passed to the `errCh` channel for further handling.
func handleStage(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	log *zerolog.Logger,
	ignore *[]string,
	stage Stage,
	buildDirectory string,
) {
	defer wg.Done()

	var err error
	log.Info().Msg(fmt.Sprintf("Handling stage %s", stage.Name))
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Failed to handle stage %s: %s", stage.Name, err))
			errCh <- err
		} else {
			log.Info().Msg(fmt.Sprintf("Finished stage %s", stage.Name))
		}
	}()

	if err := CheckContext(ctx); err != nil {
		return
	}

	to, err := mkdir(fmt.Sprintf("%s/%s", buildDirectory, stage.To))
	if err != nil {
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
			ignore,
			from,
			to,
			stage.ActionIfFileExists,
			stage.ConvertTo1251,
		)
	}
}

func makeZipFilePath(module *Module) (string, error) {
	path := filepath.Join(module.BuildDirectory, fmt.Sprintf("%s.zip", module.Version))
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}

func makeVersionDirectory(module *Module) (string, error) {
	path := filepath.Join(module.BuildDirectory, module.Version)
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}
