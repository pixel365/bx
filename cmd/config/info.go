package config

import (
	"errors"

	"github.com/spf13/cobra"
)

func infoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Short:   "Get configuration information",
		RunE: func(_ *cobra.Command, _ []string) error {
			return errors.New("not implemented")
		},
	}

	return cmd
}
