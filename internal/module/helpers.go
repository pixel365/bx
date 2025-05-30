package module

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/encoding/charmap"

	"github.com/pixel365/bx/internal/repo"

	"github.com/pixel365/bx/internal/errors"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/pixel365/bx/internal/fs"
	"github.com/pixel365/bx/internal/helpers"
)

var (
	checkPathsFunc = helpers.CheckPaths
	copyFileFunc   = fs.CopyFile
)

// ReadModule reads a module from a YAML file or directory path and returns a Module object.
//
// This function attempts to read a module from the specified path. If the `file` flag is true,
// the function treats `path` as the file path directly. Otherwise, it expects a YAML file with
// the name of the module, combining the `path` and `name` parameters to form the file path.
//
// Parameters:
//   - path (string): The directory or file path where the module file is located.
//   - name (string): The name of the module. Used to construct the file path when `file` is false.
//   - file (bool): Flag indicating whether the `path` is a direct file path or a directory where
//     a module file should be looked for.
//
// Returns:
//   - *Module: A pointer to a `Module` object if the file can be successfully read and unmarshalled.
//   - error: An error if reading or unmarshalling the file fails.
func ReadModule(path, name string, file bool) (*Module, error) {
	var filePath string
	var err error

	if !file {
		filePath, err = filepath.Abs(path + "/" + name + ".yaml")
	} else {
		filePath, err = filepath.Abs(path)
	}

	if err != nil {
		return nil, err
	}

	if !helpers.IsValidPath(filePath, path) {
		return nil, errors.ErrInvalidFilepath
	}

	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	var m Module
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func ReadModuleFromFlags(cmd *cobra.Command) (*Module, error) {
	if cmd == nil {
		return nil, errors.ErrNilCmd
	}

	name, _ := cmd.Flags().GetString("name")
	file, _ := cmd.Flags().GetString("file")
	repository, _ := cmd.Flags().GetString("repository")
	description, _ := cmd.Flags().GetString("description")
	version, _ := cmd.Flags().GetString("version")
	version = strings.TrimSpace(version)

	file = strings.TrimSpace(file)
	isFile := len(file) > 0

	path, ok := cmd.Context().Value(helpers.RootDir).(string)
	if !ok {
		return nil, errors.ErrInvalidRootDir
	}

	if !isFile && name == "" {
		err := helpers.Choose(AllModules(path), &name, "")
		if err != nil {
			return nil, err
		}
	}

	if isFile {
		path = file
	}

	module, err := ReadModule(path, name, isFile)
	if err != nil {
		return nil, err
	}

	if repository != "" {
		module.Repository = repository
	}

	if version != "" {
		module.Version = version
	}

	if description != "" {
		module.Description = description
	}

	return module, module.IsValid()
}

// AllModules returns a list of module names found in the specified directory.
//
// The function reads the directory, checks for files (skipping directories), and attempts to read
// each file as a module using the ReadModule function. If a file can be successfully read as a
// module, its name is added to the list.
//
// Parameters:
//   - directory (string): The path to the directory to scan for modules.
//
// Returns:
//   - *[]string: A pointer to a slice of strings containing the names of all successfully read modules.
func AllModules(directory string) *[]string {
	var modules []string

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(directory, file.Name())
			module, err := ReadModule(filePath, "", true)
			if err != nil {
				continue
			}

			modules = append(modules, module.Name)
		}
	}

	return &modules
}

func makeVersionDirectory(module *Module) (string, error) {
	if module == nil || module.BuildDirectory == "" {
		return "", errors.ErrNilModule
	}

	path := filepath.Join(module.BuildDirectory, module.GetVersion())
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}

func makeZipFilePath(module *Module) (string, error) {
	path := filepath.Join(module.BuildDirectory, fmt.Sprintf("%s.zip", module.GetVersion()))
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	path = filepath.Clean(path)
	return path, nil
}

func writeFileForVersion(builder *ModuleBuilder, path, content string) error {
	if len(content) == 0 {
		return nil
	}

	versionDir, err := makeVersionDirectory(builder.module)
	if err != nil {
		return err
	}

	fp := filepath.Join(versionDir, path)
	fp = filepath.Clean(fp)

	dirs := strings.Split(fp, "/")
	dirPath := strings.Join(dirs[:len(dirs)-1], "/")
	_, err = fs.MkDir(dirPath)
	if err != nil {
		return err
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil && builder.log != nil {
			builder.log.Error("Failed to close "+path, err)
		}
	}()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func makeVersionDescription(builder *ModuleBuilder) error {
	// If the full latest version is being built, then the version description file is not needed.
	// However, it may be present when copying if specified in the configuration, at the discretion of the developer.
	if builder.module.LastVersion {
		return nil
	}

	descriptionRu := strings.Builder{}
	encoder := charmap.Windows1251.NewEncoder()

	if builder.module.Description != "" {
		encodedDescriptionRu, err := encoder.String(builder.module.Description + "\n")
		if err != nil {
			return fmt.Errorf("encoding description [%s]: %w", builder.module.Description, err)
		}

		_, _ = descriptionRu.WriteString(encodedDescriptionRu)
	} else {
		if builder.module.Repository == "" {
			return nil
		}

		commits, err := repo.ChangelogList(builder.module.Repository, builder.module.Changelog)
		if err != nil {
			return err
		}

		if len(commits) == 0 {
			return nil
		}

		for _, commit := range commits {
			encodedLine, err := encoder.String(commit + "<br>")
			if err != nil {
				return fmt.Errorf("encoding commit [%s]: %w", commit, err)
			}
			_, _ = descriptionRu.WriteString(encodedLine)
		}
	}

	footer, err := builder.module.Changelog.EncodedFooter()
	if err != nil {
		return fmt.Errorf(
			"encoding footer template [%s]: %w",
			builder.module.Changelog.FooterTemplate,
			err,
		)
	}
	_, _ = descriptionRu.WriteString(footer)

	err = writeFileForVersion(builder, "description.ru", descriptionRu.String())
	if err != nil {
		return fmt.Errorf("failed to make description file: %w", err)
	}

	return nil
}

func makeVersionFile(builder *ModuleBuilder) error {
	if builder.module.LastVersion {
		return nil
	}

	buf := versionPhpContent(builder.module.Version, time.Now())

	err := writeFileForVersion(builder, "/install/version.php", buf.String())
	if err != nil {
		return fmt.Errorf("failed to make version.php file: %w", err)
	}

	return nil
}

func versionPhpContent(version string, date time.Time) strings.Builder {
	buf := strings.Builder{}
	buf.WriteString("<?php\n")
	buf.WriteString("$arModuleVersion = array(\n")
	buf.WriteString("\t\t\"VERSION\" => \"" + version + "\",\n")
	buf.WriteString("\t\t\"VERSION_DATE\" => \"" + date.Format(time.DateTime) + "\",\n")
	buf.WriteString(");\n")

	return buf
}
