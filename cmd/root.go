package cmd

import (
	"context"
	"errors"
	"os"

	"github.com/pixel365/bx/cmd/list"

	"github.com/pixel365/bx/cmd/build"
	"github.com/pixel365/bx/cmd/check"
	"github.com/pixel365/bx/cmd/create"
	"github.com/pixel365/bx/cmd/run"
	"github.com/pixel365/bx/cmd/version"

	"github.com/pixel365/bx/cmd/push"

	errors2 "github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"

	"github.com/spf13/cobra"

	"github.com/joho/godotenv"
)

var (
	osStat            = os.Stat
	mkDir             = os.Mkdir
	initRootDirFunc   = initRootDir
	getModulesDirFunc = helpers.GetModulesDir
)

func NewRootCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bx",
		Short: "Command-line tool for developers of 1C-Bitrix platform modules.",
		PersistentPreRunE: func(command *cobra.Command, _ []string) error {
			_ = godotenv.Load()
			dirPath, err := initRootDirFunc(command)
			if err != nil {
				return err
			}

			ctx = context.WithValue(ctx, helpers.RootDir, dirPath)
			command.SetContext(ctx)

			return nil
		},
	}

	cmd.AddCommand(create.NewCreateCommand())
	cmd.AddCommand(build.NewBuildCommand())
	cmd.AddCommand(check.NewCheckCommand())
	cmd.AddCommand(push.NewPushCommand())
	cmd.AddCommand(run.NewRunCommand())
	cmd.AddCommand(version.NewVersionCommand())
	cmd.AddCommand(list.NewListCommand())

	return cmd
}

// initRootDir is responsible for initializing the root directory for the project.
// It checks if the specified root directory exists, creates it if it doesn't,
// and returns the absolute path to the directory.
//
// Parameters:
//   - cmd (*cobra.Command): The command that called this function, used to retrieve the directory flag.
//
// Returns:
//   - string: The absolute path to the root directory of the project.
//   - error: An error if the directory cannot be created or accessed.
func initRootDir(cmd *cobra.Command) (string, error) {
	if cmd == nil {
		return "", errors2.NilCmdError
	}

	dirPath, err := getModulesDirFunc()
	if err != nil {
		return "", err
	}

	if _, err := osStat(dirPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = mkDir(dirPath, 0750)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	return dirPath, nil
}
