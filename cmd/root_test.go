package cmd

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd(context.Background())

	t.Run("", func(t *testing.T) {
		if cmd == nil {
			t.Error("cmd is nil")
		}

		if cmd.Use != "bx" {
			t.Errorf("cmd.Use should be 'bx' but got '%s'", cmd.Use)
		}

		if cmd.Short != "Command-line tool for developers of 1C-Bitrix platform modules." {
			t.Errorf("invalid cmd.Short = '%s'", cmd.Short)
		}

		if cmd.HasParent() {
			t.Errorf("cmd.HasParent() = true")
		}

		if cmd.HasFlags() {
			t.Errorf("cmd.HasFlags() = true")
		}

		if !cmd.HasSubCommands() {
			t.Errorf("cmd.HasSubCommands() = false")
		}

		if cmd.Hidden {
			t.Errorf("cmd.Hidden = %v", cmd.Hidden)
		}
	})
}

func Test_initRootDir_GetModDir(t *testing.T) {
	t.Run("mod dir", func(t *testing.T) {
		origGetModulesDirFunc := getModulesDirFunc
		getModulesDirFunc = func() (string, error) {
			return "", errors.New("get modules dir error")
		}
		defer func() {
			getModulesDirFunc = origGetModulesDirFunc
		}()

		_ = NewRootCmd(context.Background())
		_, err := initRootDir()
		if err == nil {
			t.Errorf("err is nil")
		}
	})
}

func Test_initRootDir_ValidDir(t *testing.T) {
	t.Run("valid dir", func(t *testing.T) {
		origGetModulesDirFunc := getModulesDirFunc
		getModulesDirFunc = func() (string, error) {
			return ".", nil
		}
		defer func() {
			getModulesDirFunc = origGetModulesDirFunc
		}()

		_ = NewRootCmd(context.Background())
		res, _ := initRootDir()
		if res != "." {
			t.Errorf("res = %v, want %v", res, "/some/dir")
		}
	})
}

func Test_initRootDir_InvalidDir(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		origGetModulesDirFunc := getModulesDirFunc
		origOsStat := osStat
		osStat = func(string) (os.FileInfo, error) {
			return nil, errors.New("get modules dir error")
		}
		defer func() {
			osStat = origOsStat
		}()

		getModulesDirFunc = func() (string, error) {
			return ".", nil
		}
		defer func() {
			getModulesDirFunc = origGetModulesDirFunc
		}()

		_ = NewRootCmd(context.Background())
		_, err := initRootDir()
		if err == nil {
			t.Errorf("err is nil")
		}
	})
}

func Test_initRootDir_MkDirError(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		origGetModulesDirFunc := getModulesDirFunc
		origOsStat := osStat
		origMkDir := mkDir

		osStat = func(string) (os.FileInfo, error) {
			return nil, os.ErrNotExist
		}
		defer func() {
			osStat = origOsStat
		}()

		getModulesDirFunc = func() (string, error) {
			return ".", nil
		}
		defer func() {
			getModulesDirFunc = origGetModulesDirFunc
		}()

		mkDir = func(name string, perm os.FileMode) error {
			return errors.New("mkdir error")
		}
		defer func() {
			mkDir = origMkDir
		}()

		_ = NewRootCmd(context.Background())
		_, err := initRootDir()
		if err == nil {
			t.Errorf("err is nil")
		}
	})
}

func Test_Cmd_PersistentPreRunE(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		origInitRootDir := initRootDirFunc
		initRootDirFunc = func() (string, error) {
			return "", nil
		}
		defer func() {
			initRootDirFunc = origInitRootDir
		}()

		cmd := NewRootCmd(context.Background())
		if err := cmd.PersistentPreRunE(cmd, []string{}); err != nil {
			t.Errorf("cmd.PersistentPreRunE() = %v, want nil", err)
		}
	})
}

func Test_Cmd_PersistentPreRunE_Err(t *testing.T) {
	t.Run("invalid dir", func(t *testing.T) {
		origInitRootDir := initRootDirFunc
		initRootDirFunc = func() (string, error) {
			return "", errors.New("init root dir error")
		}
		defer func() {
			initRootDirFunc = origInitRootDir
		}()

		cmd := NewRootCmd(context.Background())
		if err := cmd.PersistentPreRunE(cmd, []string{}); err == nil {
			t.Errorf("cmd.PersistentPreRunE() = %v, want nil", err)
		}
	})
}
