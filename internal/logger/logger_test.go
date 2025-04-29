package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
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
	t.Run("NewFileZeroLogger", func(t *testing.T) {
		filePath := fmt.Sprintf("./_%d.log", time.Now().UTC().Unix())
		logger := NewFileZeroLogger(filePath, "./")
		defer logger.Cleanup()
		defer func() {
			err := os.Remove(filePath)
			if err != nil {
				return
			}
		}()

		if logger == nil {
			t.Errorf("NewFileZeroLogger() = nil, want not nil")
		}
	})
}

func TestZeroLogger_Info(t *testing.T) {
	t.Run("info", func(t *testing.T) {
		filePath := fmt.Sprintf("./_%d.log", time.Now().UTC().Unix())
		filePath = filepath.Clean(filePath)
		logger := NewFileZeroLogger(filePath, "./")
		defer logger.Cleanup()
		defer func() {
			err := os.Remove(filePath)
			if err != nil {
				return
			}
		}()

		logger.Info("test")
		data, err := os.ReadFile(filePath)
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
		filePath := fmt.Sprintf("./_%d.log", time.Now().UTC().Unix())
		filePath = filepath.Clean(filePath)
		logger := NewFileZeroLogger(filePath, "./")
		defer logger.Cleanup()
		defer func() {
			err := os.Remove(filePath)
			if err != nil {
				return
			}
		}()

		logger.Info("test: %s", "info")
		data, err := os.ReadFile(filePath)
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
		filePath := fmt.Sprintf("./_%d.log", time.Now().UTC().Unix())
		filePath = filepath.Clean(filePath)
		logger := NewFileZeroLogger(filePath, "./")
		defer logger.Cleanup()
		defer func() {
			err := os.Remove(filePath)
			if err != nil {
				return
			}
		}()

		logger.Error(test, errors.New(test))
		data, err := os.ReadFile(filePath)
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
		filePath := fmt.Sprintf("./_%d.log", time.Now().UTC().Unix())
		filePath = filepath.Clean(filePath)
		logger := NewFileZeroLogger(filePath, "./")
		defer logger.Cleanup()
		defer func() {
			err := os.Remove(filePath)
			if err != nil {
				return
			}
		}()

		logger.Error("test: %s", errors.New(test), "error")
		data, err := os.ReadFile(filePath)
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

func TestZeroLogger_Cleanup(t *testing.T) {
	t.Run("Cleanup", func(t *testing.T) {
		filePath := fmt.Sprintf("./_%d.log", time.Now().UTC().Unix())
		filePath = filepath.Clean(filePath)
		logger := NewFileZeroLogger(filePath, "./")
		defer func() {
			err := os.Remove(filePath)
			if err != nil {
				return
			}
		}()

		logger.Cleanup()
	})
}
