package cmd

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd(context.Background())

	assert.NotNil(t, cmd)
	assert.Equal(t, "bx", cmd.Use)
	assert.Equal(t, "Command-line tool for developers of 1C-Bitrix platform modules.", cmd.Short)
	assert.False(t, cmd.HasParent())
	assert.False(t, cmd.HasFlags())
	assert.True(t, cmd.HasSubCommands())
	assert.False(t, cmd.Hidden)
}

func Test_initRootDir_GetModDir(t *testing.T) {
	origGetModulesDirFunc := getModulesDirFunc
	getModulesDirFunc = func() (string, error) {
		return "", errors.New("get modules dir error")
	}
	defer func() {
		getModulesDirFunc = origGetModulesDirFunc
	}()

	_ = NewRootCmd(context.Background())
	_, err := initRootDir()
	require.Error(t, err)
}

func Test_initRootDir_ValidDir(t *testing.T) {
	origGetModulesDirFunc := getModulesDirFunc
	getModulesDirFunc = func() (string, error) {
		return ".", nil
	}
	defer func() {
		getModulesDirFunc = origGetModulesDirFunc
	}()

	_ = NewRootCmd(context.Background())
	res, _ := initRootDir()
	assert.Equal(t, res, ".")
}

func Test_initRootDir_InvalidDir(t *testing.T) {
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
	require.Error(t, err)
}

func Test_initRootDir_MkDirError(t *testing.T) {
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
	require.Error(t, err)
}

func Test_Cmd_PersistentPreRunE(t *testing.T) {
	origInitRootDir := initRootDirFunc
	initRootDirFunc = func() (string, error) {
		return "", nil
	}
	defer func() {
		initRootDirFunc = origInitRootDir
	}()

	cmd := NewRootCmd(context.Background())
	err := cmd.PersistentPreRunE(cmd, []string{})
	require.NoError(t, err)
}

func Test_Cmd_PersistentPreRunE_Err(t *testing.T) {
	origInitRootDir := initRootDirFunc
	initRootDirFunc = func() (string, error) {
		return "", errors.New("init root dir error")
	}
	defer func() {
		initRootDirFunc = origInitRootDir
	}()

	cmd := NewRootCmd(context.Background())
	err := cmd.PersistentPreRunE(cmd, []string{})
	require.Error(t, err)
}
