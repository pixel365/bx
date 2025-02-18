package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Check a module",
		RunE: func(cmd *cobra.Command, args []string) error {
			return check(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")

	return cmd
}

func check(cmd *cobra.Command, _ []string) error {
	path := cmd.Context().Value(internal.RootDir).(string)
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	file, err := cmd.Flags().GetString("file")
	file = strings.TrimSpace(file)
	if err != nil {
		return err
	}

	isFile := len(file) > 0

	if !isFile && name == "" {
		err := internal.Choose(internal.AllModules(path), &name, "")
		if err != nil {
			return err
		}
	}

	if isFile {
		path = file
	}

	module, err := internal.ReadModule(path, name, isFile)
	if err != nil {
		return err
	}

	if err := module.IsValid(); err != nil {
		return err
	}

	if err := internal.CheckMapping(module); err != nil {
		return err
	}

	fmt.Println("ok")

	return nil
}
