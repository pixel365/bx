// Package logger provides a structured file-based logger implementation using the zerolog library.
//
// This package defines a ZeroLogger type that implements logging to a file with support for
// structured logs and automatic cleanup. It is used for tracking build processes, errors,
// and general operational output in the application.
package logger

import (
	"os"
	"path/filepath"

	"github.com/pixel365/bx/internal/types"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZeroLogger is a file-based structured logger powered by zerolog.
//
// It writes logs to a specified file and supports informational and error messages
// with optional formatting. The logger includes automatic cleanup and renaming logic
// to move the log file to a target directory after use.
type ZeroLogger struct {
	logger zerolog.Logger
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
	if len(args) > 0 {
		l.logger.Info().Msgf(message, args...)
		return
	}
	l.logger.Info().Msg(message)
}

// Error logs an error message along with an associated error.
//
// If formatting arguments are provided, the message is formatted accordingly.
// Logging is skipped if the log file is not initialized.
//
// Parameters:
//   - message: The log message or format string.
//   - err: The error object to include in the log.
//   - args: Optional format arguments for the message.
func (l *ZeroLogger) Error(message string, err error, args ...interface{}) {
	if len(args) > 0 {
		l.logger.Error().Err(err).Msgf(message, args...)
		return
	}
	l.logger.Error().Err(err).Msg(message)
}

// NewFileLogger creates and returns a new instance of ZeroLogger,
// which writes logs both to stdout and to a rotating file.
//
// The log file is created in the specified directory with the name "<moduleName>.log".
// File rotation is handled by lumberjack.Logger
// based on the provided configuration in the log parameter.
//
// Parameters:
//   - log: pointer to types.Log containing log directory and rotation settings
//   - moduleName: name of the module used to construct the log filename
//
// Returns:
//   - A pointer to a ZeroLogger instance configured with file and stdout output
func NewFileLogger(log *types.Log, moduleName string) *ZeroLogger {
	if log == nil {
		return &ZeroLogger{}
	}

	fullPath := filepath.Join(log.Dir, moduleName+".log")
	fullPath, _ = filepath.Abs(fullPath)
	fullPath = filepath.Clean(fullPath)

	err := os.MkdirAll(filepath.Dir(fullPath), 0750)
	if err != nil {
		panic("Failed to create log directory: " + fullPath)
	}

	logFile := &lumberjack.Logger{
		Filename:   fullPath,
		MaxSize:    log.MaxSize,
		MaxBackups: log.MaxBackups,
		MaxAge:     log.MaxAge,
		Compress:   log.Compress,
	}

	logger := &ZeroLogger{
		logger: zerolog.New(logFile).With().Timestamp().Logger(),
	}

	return logger
}
