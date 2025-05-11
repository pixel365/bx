// Package logger provides a structured file-based logger implementation using the zerolog library.
//
// This package defines a ZeroLogger type that implements logging to a file with support for
// structured logs and automatic cleanup. It is used for tracking build processes, errors,
// and general operational output in the application.
package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pixel365/bx/internal/helpers"

	"github.com/rs/zerolog"
)

// ZeroLogger is a file-based structured logger powered by zerolog.
//
// It writes logs to a specified file and supports informational and error messages
// with optional formatting. The logger includes automatic cleanup and renaming logic
// to move the log file to a target directory after use.
type ZeroLogger struct {
	logger  zerolog.Logger
	logFile *os.File
	logDir  string
}

// Info logs an informational message using zerolog.
//
// If formatting arguments are provided, the message is formatted using fmt-style verbs.
// Logging is skipped if the log file is not initialized.
//
// Parameters:
//   - message: The log message or format string.
//   - args: Optional format arguments for the message.
func (l *ZeroLogger) Info(message string, args ...interface{}) {
	if l.logFile != nil {
		if len(args) > 0 {
			l.logger.Info().Msgf(message, args...)
			return
		}
		l.logger.Info().Msg(message)
	}
}

// Error logs an error message along with an associated error.
//
// If formatting arguments are provided, the message is formatted accordingly.
// Logging is skipped if the log file is not initialized.
//
// Parameters:
//   - message: The log message or format string.
//   - err:     The error object to include in the log.
//   - args:    Optional format arguments for the message.
func (l *ZeroLogger) Error(message string, err error, args ...interface{}) {
	if l.logFile != nil {
		if len(args) > 0 {
			l.logger.Error().Err(err).Msgf(message, args...)
			return
		}
		l.logger.Error().Err(err).Msg(message)
	}
}

// Cleanup closes the log file and moves it to the configured log directory.
//
// If a rename operation fails, it is silently ignored. If no file was opened,
// Cleanup does nothing.
func (l *ZeroLogger) Cleanup() {
	if l.logFile != nil {
		helpers.Cleanup(l.logFile, nil)

		path := fmt.Sprintf("%s/%s", l.logDir, l.logFile.Name())
		path = filepath.Clean(path)

		err := os.Rename(l.logFile.Name(), path)
		if err != nil {
			return
		}
	}
}

// NewFileZeroLogger creates a new ZeroLogger instance that logs to a file.
//
// The file is opened in append mode and created if it does not exist. If the file
// cannot be opened, the function panics. The resulting logger is timestamped
// and configured for structured output.
//
// Parameters:
//   - filePath: Path to the log file to write to.
//   - logDir:   Directory where the file will be moved after cleanup.
//
// Returns:
//   - *ZeroLogger: A configured logger instance ready for use.
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
