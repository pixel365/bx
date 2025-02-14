package config

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
	cfg "github.com/pixel365/bx/internal/config"
)

func resetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration",
		RunE: func(_ *cobra.Command, _ []string) error {
			conf, err := cfg.GetConfig()
			if err != nil {
				return err
			}

			confirm := false
			if err = internal.Confirmation(&confirm,
				"Are you sure you want to reset all settings and clear the configuration file?"); err != nil {
				return err
			}

			if confirm {
				if err = conf.Reset(); err != nil {
					return err
				}

				color.Green("Configuration file cleared")
			}

			return nil
		},
	}

	return cmd
}
