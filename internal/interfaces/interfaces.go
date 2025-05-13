// Package interfaces defines contracts for core components involved in
// the build and execution workflow of the BX system.
//
// These interfaces abstract behaviors such as module building, configuration,
// user prompting, and execution logging, enabling flexible and testable implementations.
package interfaces

import (
	"context"
	"sync"

	"github.com/pixel365/bx/internal/types"
)

// Builder defines the core methods required for building and managing a module lifecycle.
//
// Methods:
//   - Build: Executes the primary build logic.
//   - Prepare: Performs any setup required before building.
//   - Rollback: Reverts changes if the build fails or is canceled.
//   - Collect: Gathers or stages files/metadata necessary for the build.
//   - Cleanup: Releases resources or performs cleanup actions after the build.
type Builder interface {
	Build(ctx context.Context) error
	Prepare() error
	Rollback() error
	Collect(ctx context.Context) error
	Cleanup()
}

// ModuleConfig provides access to parsed module configuration data.
//
// Typically sourced from a YAML or similar configuration file.
//
// Methods:
//   - GetVariables: Returns key-value pairs for variable substitution.
//   - GetRun: Returns the ordered map of run commands.
//   - GetStages: Returns the defined execution stages.
//   - GetIgnore: Returns file paths or patterns to ignore.
//   - GetChanges: Returns changelog-related metadata.
//   - IsLastVersion: Indicates whether the module represents the latest version.
type ModuleConfig interface {
	GetVariables() map[string]string
	GetRun() map[string][]string
	GetStages() []types.Stage
	GetIgnore() []string
	GetChanges() *types.Changes
	IsLastVersion() bool
}

// BuildLogger provides structured logging during the build process.
//
// Methods:
//   - Info: Logs informational messages with optional formatting arguments.
//   - Error: Logs error messages with the associated error and optional context.
//   - Cleanup: Finalizes or flushes any pending logs or resources.
type BuildLogger interface {
	Info(message string, args ...interface{})
	Error(message string, err error, args ...interface{})
	Cleanup()
}

// Runnable defines hooks for executing logic before and after a build stage.
//
// Methods:
//   - PreRun: Executes logic before the main run phase; receives context and logger.
//   - PostRun: Executes logic after the run phase completes.
type Runnable interface {
	PreRun(ctx context.Context, wg *sync.WaitGroup, logger BuildLogger)
	PostRun(ctx context.Context, wg *sync.WaitGroup, logger BuildLogger)
}

// Prompter abstracts user input collection and validation.
//
// Methods:
//   - Input: Prompts the user with a title and validates input.
//   - GetValue: Returns the result of the most recent input.
type Prompter interface {
	GetValue() string
	Input(title string, validator func(string) error) error
}
