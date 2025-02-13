package account

import (
	"errors"

	"github.com/spf13/cobra"
)

func authCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate with an account",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented")
		},
	}

	return cmd
}
