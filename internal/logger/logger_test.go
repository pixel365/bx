package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

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
	t.Run("NewFileLogger", func(t *testing.T) {
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

		if logger == nil {
			t.Errorf("NewFileLogger() = nil, want not nil")
		}
	})
}

func TestNilZeroLogger(t *testing.T) {
	t.Run("NilZeroLogger", func(t *testing.T) {
		logger := NewFileLogger(nil, "")
		if !reflect.DeepEqual(logger, &ZeroLogger{}) {
			t.Errorf("NewFileLogger() = nil, want not nil")
		}
	})
}

func TestZeroLogger_Info(t *testing.T) {
	t.Run("info", func(t *testing.T) {
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
		if err != nil {
			t.Errorf("os.ReadFile(%s) = %v", filePath, err)
			return
		}

		var record InfoRecord
		err = json.Unmarshal(data, &record)
		if err != nil {
			t.Errorf("json.Unmarshal() = %v", err)
		}

		if record.Level != lInfo {
			t.Errorf("record.Level = %s, want info", record.Level)
		}

		if record.Message != test {
			t.Errorf("record.Message = %s, want %s", record.Message, test)
		}
	})
}

func TestZeroLogger_Info_With_Args(t *testing.T) {
	t.Run("info", func(t *testing.T) {
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
		if err != nil {
			t.Errorf("os.ReadFile(%s) = %v", filePath, err)
			return
		}

		var record InfoRecord
		err = json.Unmarshal(data, &record)
		if err != nil {
			t.Errorf("json.Unmarshal() = %v", err)
		}

		if record.Level != lInfo {
			t.Errorf("record.Level = %s, want info", record.Level)
		}

		if record.Message != "test: info" {
			t.Errorf("record.Message = %s, want %s", record.Message, "test: info")
		}
	})
}

func TestZeroLogger_Error(t *testing.T) {
	t.Run("error", func(t *testing.T) {
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
		if err != nil {
			t.Errorf("os.ReadFile(%s) = %v", filePath, err)
			return
		}

		var record ErrorRecord
		err = json.Unmarshal(data, &record)
		if err != nil {
			t.Errorf("json.Unmarshal() = %v", err)
		}

		if record.Level != lError {
			t.Errorf("record.Level = %s, want info", record.Level)
		}

		if record.Message != test {
			t.Errorf("record.Message = %s, want %s", record.Message, test)
		}

		if record.Error != test {
			t.Errorf("record.Error = %s, want %s", record.Error, test)
		}
	})
}

func TestZeroLogger_Error_With_Args(t *testing.T) {
	t.Run("error", func(t *testing.T) {
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
		if err != nil {
			t.Errorf("os.ReadFile(%s) = %v", filePath, err)
			return
		}

		var record ErrorRecord
		err = json.Unmarshal(data, &record)
		if err != nil {
			t.Errorf("json.Unmarshal() = %v", err)
		}

		if record.Level != lError {
			t.Errorf("record.Level = %s, want info", record.Level)
		}

		if record.Message != "test: error" {
			t.Errorf("record.Message = %s, want %s", record.Message, "test: error")
		}

		if record.Error != test {
			t.Errorf("record.Error = %s, want %s", record.Error, test)
		}
	})
}

func TestNewFileLogger_PanicOnInvalidDir(t *testing.T) {
	log := &types.Log{
		Dir: string([]byte{0x00}),
	}

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("NewFileLogger() did not panic on invalid dir")
		} else {
			t.Log(err)
		}
	}()

	NewFileLogger(log, "test-module")
}
