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
	if err := CheckContextActivity(m.Ctx); err != nil {
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
	log.Info().Msg("building module")

	if err := m.Prepare(&log); err != nil {
		log.Error().Err(err).Msg("failed to prepare build")
		return m.Rollback(&log)
	}

	if err := m.Collect(&log); err != nil {
		log.Error().Err(err).Msg("failed to collect build")
		return m.Rollback(&log)
	}

	if err := m.Cleanup(&log); err != nil {
		log.Error().Err(err).Msg("failed to cleanup build")
		return err
	}

	if err := m.Push(&log); err != nil {
		log.Error().Err(err).Msg("failed to push build")
		return err
	}

	return nil
}

func (m *Module) Prepare(log *zerolog.Logger) error {
	if err := m.IsValid(); err != nil {
		log.Error().Err(err).Msg("prepare: module is invalid")
		return err
	}

	if err := CheckStages(m); err != nil {
		log.Error().Err(err).Msg("prepare: check stages failed")
		return err
	}

	if m.BuildDirectory == "" {
		m.BuildDirectory = "./build"
	}

	if m.LogDirectory == "" {
		m.LogDirectory = "./log"
	}

	path, err := mkdir(m.BuildDirectory)
	if err != nil {
		log.Error().Err(err).Msg("prepare: failed to make build directory")
		return err
	}

	m.BuildDirectory = path

	path, err = mkdir(m.LogDirectory)
	if err != nil {
		log.Error().Err(err).Msg("prepare: failed to make log directory")
		return err
	}

	m.LogDirectory = path

	_, err = mkdir(fmt.Sprintf("%s/%s", m.BuildDirectory, m.Version))
	if err != nil {
		log.Error().Err(err).Msg("prepare: failed to make build version directory")
		return err
	}

	return nil
}

func (m *Module) Cleanup(log *zerolog.Logger) error {
	//TODO: implementation
	return nil
}

func (m *Module) Rollback(log *zerolog.Logger) error {
	//TODO: implementation
	return nil
}

func (m *Module) Push(log *zerolog.Logger) error {
	//TODO: implementation
	return nil
}

func (m *Module) Collect(log *zerolog.Logger) error {
	buildDirectory, err := filepath.Abs(fmt.Sprintf("%s/%s", m.BuildDirectory, m.Version))
	if err != nil {
		log.Error().Err(err).Msg("failed to make build directory")
		return err
	}

	buildDirectory = filepath.Clean(buildDirectory)

	var wg sync.WaitGroup
	errCh := make(chan error, len(m.Stages))

	for _, item := range m.Stages {
		if err := CheckContextActivity(m.Ctx); err != nil {
			return err
		}

		wg.Add(1)
		go handleItem(m.Ctx, &wg, errCh, &m.Ignore, item, buildDirectory)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		log.Error().Int("errors", len(errs)).Msg("failed to collect build")
		return fmt.Errorf("errors: %v", errs)
	}

	return nil
}

func handleItem(
	ctx context.Context,
	wg *sync.WaitGroup,
	errCh chan<- error,
	ignore *[]string,
	item Item,
	buildDirectory string,
) {
	defer wg.Done()

	if err := CheckContextActivity(ctx); err != nil {
		errCh <- err
		return
	}

	to, err := mkdir(fmt.Sprintf("%s/%s", buildDirectory, item.To))
	if err != nil {
		errCh <- err
		return
	}

	for _, from := range item.From {
		if err := CheckContextActivity(ctx); err != nil {
			errCh <- err
			return
		}

		wg.Add(1)
		go copyFromPath(ctx, wg, errCh, ignore, from, to, item.ActionIfFileExists)
	}
}
