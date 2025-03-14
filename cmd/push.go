package cmd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"github.com/charmbracelet/huh/spinner"

	"github.com/charmbracelet/huh"

	"github.com/pixel365/bx/internal"

	"github.com/spf13/cobra"
)

func newPushCommand() *cobra.Command {
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

	return cmd
}

// push handles the logic for pushing a module to the Marketplace.
// It validates the module name, reads the module configuration, and authenticates the user.
// The module is then uploaded to the specified server after authentication.
//
// Parameters:
// - cmd (*cobra.Command): The Cobra command that invoked the push function.
// - args ([]string): A slice of arguments passed to the command (unused here).
//
// Returns:
// - error: An error if any validation or upload step fails.
func push(cmd *cobra.Command, _ []string) error {
	module, err := internal.ReadModuleFromFlags(cmd)
	if err != nil {
		return err
	}

	version, err := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)
	if err != nil {
		return err
	}

	module.Ctx = cmd.Context()

	if version != "" {
		if err := internal.ValidateVersion(version); err != nil {
			return err
		}
		module.Version = version
	}

	if err = module.IsValid(); err != nil {
		return err
	}

	password, err := handlePassword(cmd, module)
	if err != nil {
		return err
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	httpClient := internal.NewClient(client, jar)

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
		return err
	}

	err = spinner.New().
		Title("Uploading module to partners.1c-bitrix.ru...").
		Type(spinner.Dots).
		ActionWithErr(func(ctx context.Context) error {
			return httpClient.UploadZIP(module, cookies)
		}).
		Run()
	if err != nil {
		return err
	}

	return nil
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
func handlePassword(cmd *cobra.Command, module *internal.Module) (string, error) {
	password, err := cmd.Flags().GetString("password")
	if err != nil {
		return "", err
	}

	if password == "" {
		password = os.Getenv(module.PasswordEnv())
	}

	if password == "" {
		if err := huh.NewInput().
			Title("Enter Password:").
			Prompt("> ").
			Value(&password).
			EchoMode(2).
			Validate(func(input string) error {
				return internal.ValidatePassword(input)
			}).
			Run(); err != nil {
			return "", err
		}
	}

	password = strings.TrimSpace(password)
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}

	if err := internal.ValidatePassword(password); err != nil {
		return "", err
	}

	return password, nil
}
