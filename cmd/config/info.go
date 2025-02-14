package config

import (
	cfg "github.com/pixel365/bx/internal/config"

	"github.com/spf13/cobra"
)

func infoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "info",
		Aliases: []string{"i"},
		Short:   "Get configuration information",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := cfg.GetConfig()
			if err != nil {
				return err
			}

			verbose, _ := c.Flags().GetBool("verbose")
			conf.PrintSummary(verbose)

			return nil
		},
	}

	cmd.Flags().BoolP("verbose", "v", false, "Show extended information")

	return cmd
}
