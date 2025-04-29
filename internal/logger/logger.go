package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"
)

type ZeroLogger struct {
	logger  zerolog.Logger
	logFile *os.File
	logDir  string
}

func (l *ZeroLogger) Info(message string, args ...interface{}) {
	if l.logFile == nil {
		return
	}

	if len(args) > 0 {
		l.logger.Info().Msgf(message, args...)
		return
	}

	l.logger.Info().Msg(message)
}

func (l *ZeroLogger) Error(message string, err error, args ...interface{}) {
	if l.logFile == nil {
		return
	}

	if len(args) > 0 {
		l.logger.Error().Err(err).Msgf(message, args...)
		return
	}

	l.logger.Error().Err(err).Msg(message)
}

func (l *ZeroLogger) Cleanup() {
	if l.logFile != nil {
		err := l.logFile.Close()
		if err != nil {
			return
		}

		path := fmt.Sprintf("%s/%s", l.logDir, l.logFile.Name())
		path = filepath.Clean(path)
		err = os.Rename(l.logFile.Name(), path)
		if err != nil {
			return
		}
	}
}

func NewFileZeroLogger(filePath, logDir string) *ZeroLogger {
	filePath = filepath.Clean(filePath)
	logFile, err := os.OpenFile(
		filePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0600,
	)

	if err != nil {
		panic(err)
	}

	logger := &ZeroLogger{
		logger:  zerolog.New(logFile).With().Timestamp().Logger(),
		logFile: logFile,
		logDir:  logDir,
	}

	return logger
}
