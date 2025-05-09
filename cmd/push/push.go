package push

import (
	"github.com/spf13/cobra"
)

func NewPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Push module to a Marketplace",
		Example: `
# Push module to a registry
bx push --name my_module

# Push a module by file path
bx push -f config.yaml

# Override version
bx push --name my_module --version 1.2.3
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return push(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("version", "v", "", "Version of the module")
	cmd.Flags().StringP("password", "p", "", "Account password")
	cmd.Flags().BoolP("silent", "s", false, "Silent mode")

	return cmd
}
