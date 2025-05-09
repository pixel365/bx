package list

import (
	"fmt"
	"maps"
	"slices"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/auth"
	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/types"
)

var (
	readModuleFromFlagsFunc = module.ReadModuleFromFlags
	authFunc                = auth.Authenticate
	inputPasswordFunc       = auth.InputPassword
)

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all module versions",
		Example: `
# List all module versions by name
bx list --name my_module

# List all module versions by file path
bx list -f config.yaml
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("password", "p", "", "Account password")
	cmd.Flags().BoolP("head", "", false, "Show last module version")
	cmd.Flags().StringP("sort", "", "", "Sort module versions by name")
	cmd.Flags().BoolP("silent", "s", false, "Silent mode")

	return cmd
}

func list(cmd *cobra.Command, _ []string) error {
	mod, err := readModuleFromFlagsFunc(cmd)
	if err != nil {
		return err
	}

	password, err := inputPasswordFunc(cmd, mod)
	if err != nil {
		return err
	}

	silent, _ := cmd.Flags().GetBool("silent")
	httpClient, cookies, err := authFunc(mod, password, silent)
	if err != nil {
		return err
	}

	versions, err := httpClient.Versions(mod, cookies)
	if err != nil {
		return err
	}

	head, _ := cmd.Flags().GetBool("head")
	s, _ := cmd.Flags().GetString("sort")

	sorting := types.Desc

	if s != "" {
		switch s {
		case string(types.Asc), string(types.Desc):
			sorting = types.SortingType(s)
		default:
			return errors.InvalidArgumentError
		}
	}

	items := slices.Sorted(maps.Keys(versions))

	if sorting == types.Asc {
		for _, version := range items {
			fmt.Printf("%s (%s)\n", version, versions[version])
			if head {
				break
			}
		}
		return nil
	}

	for i := len(items) - 1; i >= 0; i-- {
		fmt.Printf("%s (%s)\n", items[i], versions[items[i]])
		if head {
			break
		}
	}

	return nil
}
