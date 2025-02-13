package module

import (
	"errors"

	"github.com/spf13/cobra"
)

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List available modules",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented")
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")

	return cmd
}
