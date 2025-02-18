package internal

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/charmbracelet/huh"
)

type Cfg string

const (
	RootDir Cfg = "root_dir"

	Yes = "Yes"
	No  = "No"
)

type Printer interface {
	PrintSummary(verbose bool)
}

type OptionProvider interface {
	Option() string
}

func Confirmation(flag *bool, title string) error {
	if err := huh.NewConfirm().
		Title(title).
		Affirmative(Yes).
		Negative(No).
		Value(flag).
		Run(); err != nil {
		return err
	}

	return nil
}

func Choose(items *[]string, value *string, title string) error {
	if len(*items) == 0 {
		switch any(items).(type) {
		default:
			return errors.New("no items")
		}
	}

	var options []huh.Option[string]
	for _, item := range *items {
		options = append(options, huh.NewOption(item, item))
	}

	if err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(value).
		Run(); err != nil {
		return err
	}

	return nil
}

func CaptureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}

	os.Stdout = stdout
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)

	return buf.String()
}

func ResultMessage(format string, a ...any) {
	if len(a) == 0 {
		fmt.Println(format)
	} else {
		fmt.Printf(format, a...)
	}
}

func GetModulesDir(path string) (string, error) {
	var err error
	dirPath := path
	if dirPath == "" {
		dirPath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	return filepath.Abs(fmt.Sprintf("%s/.bx", dirPath))
}

func DefaultYAML() string {
	return `name: test  # The name of the project or build.
version: 1.0.0  # The version of the project or build.
account: test  # The account associated with the project.
repository: ""  # The repository URL where the project is stored (can be empty if not specified).
buildDirectory: "./dist"  # Directory where the build artifacts will be output.
logDirectory: "./logs"  # Directory where log files will be stored.

mapping:
  - name: "components"  # Name of the mapping, describing what the mapping represents (e.g., components).
    # This can be any name that makes sense for your project, used for your own convenience.
    relativePath: "install/components"  # Relative path in the project to map files to.
    ifFileExists: "replace"  # Action to take if the file already exists (options: replace, skip, copy-new).
    paths:
      - ./examples/structure/bitrix/components  # List of paths to files that will be mapped.
      - ./examples/structure/local/components

  - name: "templates"
    relativePath: "install/templates"
    ifFileExists: "replace"
    paths:
      - ./examples/structure/bitrix/templates
      - ./examples/structure/local/templates

  - name: "rootFiles"
    relativePath: "."
    ifFileExists: "replace"
    paths:
      - ./examples/structure/simple-file.php

  - name: "testFiles"
    relativePath: "test"
    ifFileExists: "replace"
    paths:
      - ./examples/structure/simple-file.php

ignore:
  - "**/*.log"  # List of files or patterns to ignore during the build or processing (e.g., log files).
`
}

func AllModules(directory string) *[]string {
	var modules []string

	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(directory, file.Name())
		module, err := ReadModule(filePath, "", true)
		if err != nil {
			continue
		}

		modules = append(modules, module.Name)
	}

	return &modules
}

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

	if !isValidPath(filePath, path) {
		return nil, errors.New("invalid file path")
	}

	data, err := os.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, err
	}

	var module Module
	if err := yaml.Unmarshal(data, &module); err != nil {
		return nil, err
	}

	return &module, nil
}

func CheckPath(path string) error {
	path = filepath.Clean(path)
	if !isValidPath(path, path) {
		return fmt.Errorf("invalid path: %s", path)
	}

	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	return nil
}

func IsDir(path string) (bool, error) {
	err := CheckPath(path)
	if err != nil {
		return false, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.Mode().IsDir(), nil
}

func CheckMapping(module *Module) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(module.Mapping)*5)

	for _, item := range module.Mapping {
		wg.Add(1)
		go func(wg *sync.WaitGroup, item Item) {
			defer wg.Done()
			checkPaths(item, errCh)
		}(&wg, item)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors: %v", errs)
	}

	return nil
}

func checkPaths(item Item, ch chan<- error) {
	for _, path := range item.Paths {
		err := CheckPath(path)
		if err != nil {
			ch <- err
			return
		}
	}
}

func isValidPath(filePath, basePath string) bool {
	absBasePath, _ := filepath.Abs(basePath)
	absFilePath, _ := filepath.Abs(filePath)

	if strings.Contains(absFilePath, "..") {
		return false
	}

	return strings.HasPrefix(absFilePath, absBasePath)
}
