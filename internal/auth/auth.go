package auth

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/client"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/request"
	"github.com/pixel365/bx/internal/types"
	"github.com/pixel365/bx/internal/validators"
)

var (
	inputPasswordFunc     = helpers.UserInput
	newPasswordPromptFunc = types.NewPrompt
)

// Authenticate performs authentication against the partners.1c-bitrix.ru service
// using the provided module and password.
//
// Parameters:
//   - module (*module.Module): The module object.
//     If nil, the function returns errors.ErrNilModule.
//   - password (string): The password used for authentication.
//     If empty, the function returns errors.ErrEmptyPassword.
//   - silent (bool): Skip spinner.
//
// Returns:
//   - []*http.Cookie: A slice of cookies obtained from the authentication response, typically used for
//     session management.
//   - Error: Any error encountered during parameter validation or the authentication process.
//
// Description:
// First validates that both the module and password are provided and non-empty.
// It creates an HTTP client with an associated cookie jar and wraps it using a request.client.
// Using a spinner for progress feedback, it calls httpClient.Authenticate with the module’s account and password.
// If authentication succeeds, the function returns the configured client and cookies;
// otherwise, it returns the encountered error.
func Authenticate(
	client client.HTTPClient,
	module *module.Module,
	password string,
	silent bool,
) ([]*http.Cookie, error) {
	if module == nil {
		return nil, errors.ErrNilModule
	}

	if password == "" {
		return nil, errors.ErrEmptyPassword
	}

	var err error
	var cookies []*http.Cookie

	if silent {
		cookies, err = request.Authenticate(client, module.Account, password)
	} else {
		err = helpers.Spinner("Authenticate on partners.1c-bitrix.ru...",
			func(ctx context.Context) error {
				cookies, err = request.Authenticate(client, module.Account, password)
				return err
			})
	}

	if err != nil {
		return nil, err
	}

	return cookies, nil
}

// InputPassword manages the process of obtaining and validating the password needed for authentication.
// It first checks if the password was provided as a flag, then checks environment variables, and if neither are found,
// it prompts the user to enter a password interactively.
//
// Parameters:
// - cmd (*cobra.Command): The Cobra command that invoked the function.
// - module (*internal.Module): The module for which the password is being provided (may use environment variable).
//
// Returns:
// - string: The validated password.
// - error: An error if the password is invalid or the prompt fails.
func InputPassword(cmd *cobra.Command, module *module.Module) (string, error) {
	password, _ := cmd.Flags().GetString("password")
	password = strings.TrimSpace(password)

	if password == "" {
		password = os.Getenv(module.PasswordEnv())
	}

	if password == "" {
		prompter := newPasswordPromptFunc()
		err := inputPasswordFunc(
			prompter,
			&password,
			"Enter Password:",
			validators.ValidatePassword,
		)
		if err != nil {
			return "", err
		}
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return "", errors.ErrEmptyPassword
	}

	if err := validators.ValidatePassword(password); err != nil {
		return "", err
	}

	return password, nil
}
