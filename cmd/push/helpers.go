package push

import (
	"context"
	"net/http"
	"time"

	"github.com/pixel365/bx/internal/client"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/helpers"

	"github.com/pixel365/bx/internal/auth"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/request"
)

var (
	readModuleFromFlagsFunc = module.ReadModuleFromFlags
	uploadFunc              = upload
	authFunc                = auth.Authenticate
	inputPasswordFunc       = auth.InputPassword
	spinnerFunc             = helpers.Spinner
)

// push handles the logic for pushing a module to the Marketplace.
// It validates the module name, reads the module configuration, and authenticates the user.
// The module is then uploaded to the specified server after authentication.
//
// Parameters:
//   - cmd (*cobra.Command): The Cobra command that invoked the push function.
//   - args ([]string): A slice of arguments passed to the command (unused here).
//
// Returns:
//   - error: An error if any validation or upload step fails.
func push(cmd *cobra.Command, _ []string) error {
	mod, err := readModuleFromFlagsFunc(cmd)
	if err != nil {
		return err
	}

	label, _ := cmd.Flags().GetString("label")
	if label != "" {
		switch types.VersionLabel(label) {
		case types.Alpha, types.Beta, types.Stable:
			mod.Label = types.VersionLabel(label)
		default:
			return errors.ErrInvalidLabel
		}
	}

	password, err := inputPasswordFunc(cmd, mod)
	if err != nil {
		return err
	}

	silent, _ := cmd.Flags().GetBool("silent")

	httpClient := client.NewClient(10 * time.Second)

	cookies, err := authFunc(httpClient, mod, password, silent)
	if err != nil {
		return err
	}

	err = uploadFunc(cmd.Context(), httpClient, mod, cookies, silent)
	if err != nil {
		return err
	}

	versions := make(types.Versions, 1)
	versions[mod.Version] = mod.GetLabel()

	return request.ChangeLabels(httpClient, mod, cookies, versions)
}

func upload(
	ctx context.Context,
	client client.HTTPClient,
	module *module.Module,
	cookies []*http.Cookie,
	silent bool,
) error {
	if module == nil {
		return errors.ErrNilModule
	}

	if client == nil {
		return errors.ErrNilClient
	}

	if len(cookies) == 0 {
		return errors.ErrNilCookie
	}

	if silent {
		return request.UploadZIP(ctx, client, module, cookies)
	}

	return spinnerFunc(
		"Uploading module to partners.1c-bitrix.ru...",
		func(ctx context.Context) error {
			return request.UploadZIP(ctx, client, module, cookies)
		},
	)
}
