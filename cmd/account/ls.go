package account

import (
	"errors"

	"github.com/spf13/cobra"
)

func lsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List accounts",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented")
		},
	}

	return cmd
}
