package label

import (
	"errors"

	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/auth"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/types"
)

var (
	readModuleFromFlagsFunc = module.ReadModuleFromFlags
	authFunc                = auth.Authenticate
	inputPasswordFunc       = auth.InputPassword
)

func NewLabelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Change module label",
		Example: `
# Change module label
bx label stable
`,
		RunE: label,
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")
	cmd.Flags().StringP("version", "v", "", "Version of the module")
	cmd.Flags().StringP("file", "f", "", "Path to a module")
	cmd.Flags().StringP("password", "p", "", "Account password")
	cmd.Flags().BoolP("silent", "s", false, "Silent mode")

	return cmd
}

func label(cmd *cobra.Command, _ []string) error {
	if len(cmd.Flags().Args()) != 1 {
		return errors.New("label is required")
	}

	l := types.VersionLabel(cmd.Flags().Args()[0])
	switch l {
	case types.Alpha, types.Beta, types.Stable:
	default:
		return errors2.ErrInvalidLabel
	}

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

	versions := make(types.Versions, 1)
	versions[mod.Version] = l

	err = httpClient.ChangeLabels(mod, cookies, versions)
	if err != nil {
		return err
	}

	return nil
}
