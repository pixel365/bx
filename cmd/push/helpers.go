package push

import (
	"context"
	"net/http"

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

	password, err := inputPasswordFunc(cmd, mod)
	if err != nil {
		return err
	}

	silent, _ := cmd.Flags().GetBool("silent")
	httpClient, cookies, err := authFunc(mod, password, silent)
	if err != nil {
		return err
	}

	return uploadFunc(cmd.Context(), httpClient, mod, cookies, silent)
}

func upload(
	ctx context.Context,
	client *request.Client,
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
		return client.UploadZIP(ctx, module, cookies)
	}

	return spinnerFunc(
		"Uploading module to partners.1c-bitrix.ru...",
		func(ctx context.Context) error {
			return client.UploadZIP(ctx, module, cookies)
		},
	)
}
