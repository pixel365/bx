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

func push(cmd *cobra.Command, _ []string) error {
	path := cmd.Context().Value(internal.RootDir).(string)
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	file, err := cmd.Flags().GetString("file")
	file = strings.TrimSpace(file)
	if err != nil {
		return err
	}

	isFile := len(file) > 0
	if !isFile && name == "" {
		if err := internal.Choose(internal.AllModules(path), &name, ""); err != nil {
			return err
		}
	}

	if isFile {
		path = file
	}

	version, err := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)
	if err != nil {
		return err
	}

	module, err := internal.ReadModule(path, name, isFile)
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
