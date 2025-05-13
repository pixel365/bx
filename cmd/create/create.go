package create

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/module"
	"github.com/pixel365/bx/internal/validators"

	"gopkg.in/yaml.v3"

	"github.com/spf13/cobra"
)

var (
	moduleNameInputFunc = helpers.UserInput
	newModulePromptFunc = types.NewPrompt
)

func NewCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Aliases: []string{"c"},
		Short:   "Create a new module",
		Example: `
# Create a new module
bx create --name my_module
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return create(cmd, args)
		},
	}

	cmd.Flags().StringP("name", "n", "", "Name of the module")

	return cmd
}

// create handles the logic of creating a new module by generating a YAML file for the module configuration.
// It validates the module name and uses default values to create the module's YAML configuration file.
//
// Parameters:
//   - cmd (*cobra.Command): The Cobra command that invoked the create function.
//   - args ([]string): A slice of arguments passed to the command (unused here).
//
// Returns:
//   - error: An error if the module name is invalid or any other error occurs during the creation process.
func create(cmd *cobra.Command, _ []string) error {
	directory, ok := cmd.Context().Value(helpers.RootDir).(string)
	if !ok {
		return errors.ErrInvalidRootDir
	}

	name, _ := cmd.Flags().GetString("name")
	name = strings.TrimSpace(name)
	if name == "" {
		prompter := newModulePromptFunc()
		err := moduleNameInputFunc(prompter, &name, "Enter Module Name:", func(input string) error {
			return validators.ValidateModuleName(input, directory)
		})
		if err != nil {
			return err
		}
	} else {
		if err := validators.ValidateModuleName(name, directory); err != nil {
			return err
		}
	}

	var mod module.Module
	def := []byte(helpers.DefaultYAML())
	_ = yaml.Unmarshal(def, &mod)

	mod.Name = name
	mod.Account = ""

	out, _ := mod.ToYAML()

	filePath, _ := filepath.Abs(fmt.Sprintf("%s/%s.yaml", directory, name))

	return os.WriteFile(filePath, out, 0600)
}
