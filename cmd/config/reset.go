package config

import (
	"errors"
	"fmt"

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

			confirm, _ := c.Flags().GetBool("yes")
			if !confirm {
				if err := internal.Confirmation(&confirm,
					"Are you sure you want to reset all settings and clear the configuration file?"); err != nil {
					return err
				}
			}

			if confirm {
				if err := conf.Reset(); err != nil {
					return err
				}

				fmt.Println("Configuration file cleared")
			}

			return nil
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "Confirm reset configuration")

	return cmd
}
