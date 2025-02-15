package config

import (
	"errors"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func infoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Short:   "Get configuration information",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, ok := c.Context().Value(internal.CfgContextKey).(internal.Printer)
			if !ok {
				return errors.New("no config found in context")
			}

			verbose, _ := c.Flags().GetBool("verbose")
			conf.PrintSummary(verbose)

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")

	return cmd
}
