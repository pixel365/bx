package config

import (
	"errors"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal"
)

func resetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Reset configuration",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, ok := c.Context().Value(internal.CfgContextKey).(internal.ConfigManager)
			if !ok {
				return errors.New("no config found in context")
			}

			confirm := false
			if err := internal.Confirmation(&confirm,
				"Are you sure you want to reset all settings and clear the configuration file?"); err != nil {
				return err
			}

			if confirm {
				if err := conf.Reset(); err != nil {
					return err
				}

				color.Green("Configuration file cleared")
			}

			return nil
		},
	}

	return cmd
}
