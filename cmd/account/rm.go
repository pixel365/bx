package account

import (
	"errors"

	"github.com/spf13/cobra"
)

func rmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove an account",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented")
		},
	}

	return cmd
}
