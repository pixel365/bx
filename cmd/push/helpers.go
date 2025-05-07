package push

import (
	"context"
	"net/http"

	"github.com/pixel365/bx/internal/auth"

	"github.com/charmbracelet/huh/spinner"
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

	httpClient, cookies, err := authFunc(mod, password)
	if err != nil {
		return err
	}

	return uploadFunc(httpClient, mod, cookies)
}

func upload(client *request.Client, module *module.Module, cookies []*http.Cookie) error {
	if module == nil {
		return errors.NilModuleError
	}

	if client == nil {
		return errors.NilClientError
	}

	if len(cookies) == 0 {
		return errors.NilCookieError
	}

	err := spinner.New().
		Title("Uploading module to partners.1c-bitrix.ru...").
		Type(spinner.Dots).
		ActionWithErr(func(ctx context.Context) error {
			return client.UploadZIP(module, cookies)
		}).
		Run()
	return err
}
