package config

import "github.com/spf13/cobra"

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Manage configuration",
		Aliases: []string{"conf"},
	}

	cmd.AddCommand(infoCmd())
	cmd.AddCommand(resetCmd())

	return cmd
}
