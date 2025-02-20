package internal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

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
			fmt.Println(err)
		} else {
			path := fmt.Sprintf("%s/%s", m.LogDirectory, logFile.Name())
			path = filepath.Clean(path)
			err := os.Rename(logFile.Name(), path)
			if err != nil {
				fmt.Println(err)
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

	if err := m.Push(&log); err != nil {
		log.Error().Err(err).Msg("Failed to push build")
		return err
	}

	log.Info().Msg("Push complete")

	return nil
}

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

func (m *Module) Push(log *zerolog.Logger) error {
	//TODO: implementation
	return nil
}

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

func handleStage(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	log *zerolog.Logger,
	ignore *[]string,
	item Stage,
	buildDirectory string,
) {
	defer wg.Done()

	var err error
	log.Info().Msg(fmt.Sprintf("Handling stage %s", item.Name))
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Failed to handle stage %s: %s", item.Name, err))
			errCh <- err
		} else {
			log.Info().Msg(fmt.Sprintf("Finished stage %s", item.Name))
		}
	}()

	if err := CheckContext(ctx); err != nil {
		return
	}

	to, err := mkdir(fmt.Sprintf("%s/%s", buildDirectory, item.To))
	if err != nil {
		return
	}

	for _, from := range item.From {
		if err := CheckContext(ctx); err != nil {
			return
		}

		wg.Add(1)
		go copyFromPath(ctx, wg, errCh, ignore, from, to, item.ActionIfFileExists)
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
