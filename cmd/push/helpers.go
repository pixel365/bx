package push

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/request"
	"github.com/pixel365/bx/internal/types"
	"github.com/pixel365/bx/internal/validators"
)

var (
	inputPasswordFunc       = helpers.UserInput
	newPasswordPromptFunc   = types.NewPrompt
	readModuleFromFlagsFunc = module.ReadModuleFromFlags
	uploadFunc              = upload
	authFunc                = auth
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
	if cmd == nil {
		return errors.NilCmdError
	}

	mod, err := readModuleFromFlagsFunc(cmd)
	if err != nil {
		return err
	}

	version, _ := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)

	if version != "" {
		if err := validators.ValidateVersion(version); err != nil {
			return err
		}
		mod.Version = version
	}

	if err = mod.IsValid(); err != nil {
		return err
	}

	password, err := handlePassword(cmd, mod)
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

func auth(module *module.Module, password string) (*request.Client, []*http.Cookie, error) {
	if module == nil {
		return nil, nil, errors.NilModuleError
	}

	if password == "" {
		return nil, nil, errors.EmptyPasswordError
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	httpClient := request.NewClient(client, jar)

	var err error
	var cookies []*http.Cookie

	err = spinner.New().
		Title("Authorization on partners.1c-bitrix.ru...").
		Type(spinner.Dots).
		ActionWithErr(func(ctx context.Context) error {
			cookies, err = httpClient.Authorization(module.Account, password)
			return err
		}).
		Run()
	if err != nil {
		return nil, nil, err
	}

	return httpClient, cookies, nil
}

// handlePassword manages the process of obtaining and validating the password needed for authentication.
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
func handlePassword(cmd *cobra.Command, module *module.Module) (string, error) {
	password, _ := cmd.Flags().GetString("password")
	password = strings.TrimSpace(password)

	if password == "" {
		password = os.Getenv(module.PasswordEnv())
	}

	if password == "" {
		prompter := newPasswordPromptFunc()
		err := inputPasswordFunc(prompter, &password, "Enter Password:", func(input string) error {
			return validators.ValidatePassword(input)
		})
		if err != nil {
			return "", err
		}
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return "", errors.EmptyPasswordError
	}

	if err := validators.ValidatePassword(password); err != nil {
		return "", err
	}

	return password, nil
}
