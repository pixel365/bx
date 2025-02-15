package module

import (
	"fmt"
	"strings"

	"github.com/pixel365/bx/internal"
	"github.com/pixel365/bx/internal/config"
	"github.com/pixel365/bx/internal/model"

	"github.com/spf13/cobra"
)

func rmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm",
		Short: "Remove a module",
		RunE: func(c *cobra.Command, _ []string) error {
			conf, err := config.GetConfig()
			if err != nil {
				return err
			}

			name, _ := c.Flags().GetString("name")
			name = strings.TrimSpace(name)
			if name == "" {
				if err = internal.Choose(&conf.Modules, &name,
					"OptionProvider the module you want to delete:"); err != nil {
					return err
				}
			}

			confirm, _ := c.Flags().GetBool("yes")
			if !confirm {
				if err = internal.Confirmation(&confirm,
					fmt.Sprintf("Are you sure you want to delete module %s?", name)); err != nil {
					return err
				}
			}

			if confirm {
				deleted := false
				var modules []model.Module
				for _, m := range conf.Modules {
					if m.Name == name {
						deleted = true
						continue
					}

					modules = append(modules, m)
				}

				if !deleted {
					return fmt.Errorf("module %s not found", name)
				}

				conf.Modules = modules

				if err := conf.Save(); err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "The name of the module")
	cmd.Flags().BoolP("yes", "y", false, "Confirm deletion")

	return cmd
}
