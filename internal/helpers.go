package internal

import (
	"bytes"
	"context"
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
	return `name: test
version: 1.0.0
account: test
repository: ""
buildDirectory: "./dist"
logDirectory: "./logs"
stages:
  - name: "components"
    to: "install/components"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/bitrix/components
      - ./examples/structure/local/components

  - name: "templates"
    to: "install/templates"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/bitrix/templates
      - ./examples/structure/local/templates

  - name: "rootFiles"
    to: "."
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/simple-file.php

  - name: "testFiles"
    to: "test"
    actionIfFileExists: "replace"
    from:
      - ./examples/structure/simple-file.php

ignore:
  - "**/*.log"
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

func CheckStages(module *Module) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(module.Stages)*5)

	for _, item := range module.Stages {
		wg.Add(1)
		go func(wg *sync.WaitGroup, item Stage) {
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

func CheckContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled: %w", ctx.Err())
	default:
		return nil
	}
}

func checkPaths(item Stage, ch chan<- error) {
	for _, path := range item.From {
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
