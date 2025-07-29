package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/types"
)

type InfoRecord struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
}

type ErrorRecord struct {
	Level   string `json:"level"`
	Time    string `json:"time"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

var test = "test"

const (
	lInfo  = "info"
	lError = "error"
)

func TestNewFileZeroLogger(t *testing.T) {
	log := types.Log{
		Dir:        "./",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		LocalTime:  false,
		Compress:   false,
	}
	name := fmt.Sprintf("_%d.mod", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("%s/%s", log.Dir, name)
	logger := NewFileLogger(&log, name)
	defer func() {
		err := os.Remove(filePath)
		if err != nil {
			return
		}
	}()

	assert.NotNil(t, logger)
}

func TestNilZeroLogger(t *testing.T) {
	logger := NewFileLogger(nil, "")
	assert.Equal(t, &ZeroLogger{}, logger)
}

func TestZeroLogger_Info(t *testing.T) {
	log := types.Log{
		Dir:        "./",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		LocalTime:  false,
		Compress:   false,
	}
	name := fmt.Sprintf("_%d.mod", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("%s/%s", log.Dir, name)
	filePath = filepath.Clean(filePath)
	logger := NewFileLogger(&log, name)
	defer func() {
		err := os.Remove(filepath.Clean(filePath + ".log"))
		if err != nil {
			return
		}
	}()

	logger.Info("test")
	data, err := os.ReadFile(filepath.Clean(filePath + ".log"))
	require.NoError(t, err)

	var record InfoRecord
	err = json.Unmarshal(data, &record)
	require.NoError(t, err)

	assert.Equal(t, lInfo, record.Level)
	assert.Equal(t, test, record.Message)
}

func TestZeroLogger_Info_With_Args(t *testing.T) {
	log := types.Log{
		Dir:        "./",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		LocalTime:  false,
		Compress:   false,
	}
	name := fmt.Sprintf("_%d.mod", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("%s/%s", log.Dir, name)
	filePath = filepath.Clean(filePath)
	logger := NewFileLogger(&log, name)
	defer func() {
		err := os.Remove(filepath.Clean(filePath + ".log"))
		if err != nil {
			return
		}
	}()

	logger.Info("test: %s", "info")
	data, err := os.ReadFile(filepath.Clean(filePath + ".log"))
	require.NoError(t, err)

	var record InfoRecord
	err = json.Unmarshal(data, &record)
	require.NoError(t, err)

	assert.Equal(t, lInfo, record.Level)
	assert.Equal(t, "test: info", record.Message)
}

func TestZeroLogger_Error(t *testing.T) {
	log := types.Log{
		Dir:        "./",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		LocalTime:  false,
		Compress:   false,
	}
	name := fmt.Sprintf("_%d.mod", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("%s/%s", log.Dir, name)
	filePath = filepath.Clean(filePath)
	logger := NewFileLogger(&log, name)
	defer func() {
		err := os.Remove(filepath.Clean(filePath + ".log"))
		if err != nil {
			return
		}
	}()

	logger.Error(test, errors.New(test))
	data, err := os.ReadFile(filepath.Clean(filePath + ".log"))
	require.NoError(t, err)

	var record ErrorRecord
	err = json.Unmarshal(data, &record)
	require.NoError(t, err)

	assert.Equal(t, test, record.Error)
	assert.Equal(t, test, record.Message)
	assert.Equal(t, test, record.Error)
}

func TestZeroLogger_Error_With_Args(t *testing.T) {
	log := types.Log{
		Dir:        "./",
		MaxSize:    1,
		MaxBackups: 1,
		MaxAge:     1,
		LocalTime:  false,
		Compress:   false,
	}
	name := fmt.Sprintf("_%d.mod", time.Now().UTC().Unix())
	filePath := fmt.Sprintf("%s/%s", log.Dir, name)
	filePath = filepath.Clean(filePath)
	logger := NewFileLogger(&log, name)
	defer func() {
		err := os.Remove(filepath.Clean(filePath + ".log"))
		if err != nil {
			return
		}
	}()

	logger.Error("test: %s", errors.New(test), "error")
	data, err := os.ReadFile(filepath.Clean(filePath + ".log"))
	require.NoError(t, err)

	var record ErrorRecord
	err = json.Unmarshal(data, &record)
	require.NoError(t, err)

	assert.Equal(t, lError, record.Level)
	assert.Equal(t, "test: error", record.Message)
	assert.Equal(t, test, record.Error)
}

func TestNewFileLogger_PanicOnInvalidDir(t *testing.T) {
	log := &types.Log{
		Dir: string([]byte{0x00}),
	}

	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()

	NewFileLogger(log, "test-module")
}
