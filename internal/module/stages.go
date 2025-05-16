package module

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/fs"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/interfaces"
	"github.com/pixel365/bx/internal/types"
)

// HandleStages executes a sequence of stages defined in the provided module.
//
// For each stage name in the `stages` slice, the corresponding stage is resolved from the module `m`
// and processed concurrently via the `handleStage` function.
// Each stage may produce file copy tasks,
// which are sent to a shared channel (`filesCh`) and handled by a pool of worker goroutines.
// Log messages are sent asynchronously to a logging worker via `logCh`.
//
// The function manages synchronization using multiple WaitGroups and coordinates shutdown via
// context cancellation.
// If any error occurs in stage processing or file copying,
// the first error is captured and the context is canceled.
//
// Parameters:
//   - ctx: base context used for cancellation.
//   - stages: list of stage names to execute.
//   - m: the module containing the stage definitions.
//   - logger: implementation of BuildLogger for reporting progress and messages.
//   - customCommandMode: if true, skips versioned output directory creation.
//
// Returns:
//   - err: the first error encountered during execution, or nil if all stages completed successfully.
func HandleStages(
	ctx context.Context,
	stages []string,
	m *Module,
	logger interfaces.BuildLogger,
	customCommandMode bool,
) (err error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	dir := ""
	if !customCommandMode {
		dir, err = makeVersionDirectory(m)
		if err != nil {
			return
		}
	}

	minWorkers := runtime.NumCPU() * 2
	workersCount := m.SourceCount()
	if workersCount < minWorkers {
		workersCount = minWorkers
	}

	filesCh := make(chan types.Path, workersCount)
	errCh := make(chan error, 1)
	logCh := make(chan string, workersCount)

	var stagesWorkersWg sync.WaitGroup
	var copyFilesWg sync.WaitGroup
	var once sync.Once

	go errorWorker(errCh, cancel, &once, &err)

	go logWorker(logCh, logger)

	go copyWorkers(ctx, &copyFilesWg, filesCh, errCh, workersCount)

	for _, name := range stages {
		stage, _ := m.FindStage(name)
		stagesWorkersWg.Add(1)
		go func(stage types.Stage) {
			defer stagesWorkersWg.Done()
			handleStage(ctx, filesCh, logCh, errCh, m, stage, dir, m.StageCallback)
		}(stage)
	}

	go cleanupWorker(&stagesWorkersWg, &copyFilesWg, &once, cancel, filesCh, logCh, errCh)

	<-ctx.Done()
	return
}

// CheckStages validates the paths in the stages of the given module.
//
// This function iterates over the stages in the provided module and concurrently
// checks the paths defined in each stage using goroutines. If any errors are encountered,
// they are collected in a channel and returned as a combined error.
//
// Parameters:
//   - module (*Module): The module containing stages to be validated.
//
// Returns:
//   - error: Returns an error if any validation fails in any stage's paths.
//     If no errors are found, it returns nil.
func CheckStages(module *Module) error {
	if module == nil {
		return errors.ErrNilModule
	}

	var errs []error
	var wg sync.WaitGroup
	errCh := make(chan error, len(module.Stages)*5)

	for _, item := range module.Stages {
		wg.Add(1)
		go func(item types.Stage) {
			defer wg.Done()
			checkPathsFunc(item, errCh)
		}(item)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors: %v", errs)
	}

	return nil
}

// handleStage processes a single build stage concurrently and sends copy tasks to the shared file worker pool.
//
// The function performs the following steps:
//  1. Logs the start and completion of the stage via `logCh`.
//  2. Resolves a `Runnable` using the provided callback `cb` and executes its `PreRun` hook.
//  3. Validates the context and resolves the target directory for the stage output.
//  4. Create the target directory (recursively if needed).
//  5. For each input path in `stage.From`, spawns a goroutine that validates the context,
//     builds a copy `Path` struct, and sends it to `filesCh` to be processed by `copyWorkers`.
//  6. Wait for all spawned copy goroutines to finish.
//  7. If the runner was initialized, execute the `PostRun` hook.
//
// Errors from directory creation, pre/post run hooks, or file copy failures are sent to `errCh`.
// If the context is canceled at any point, the function exits early.
//
// Parameters:
//   - ctx: the context used for cancellation and timeout propagation.
//   - filesCh: a shared channel for submitting file copy tasks.
//   - logCh: a channel for emitting informational log messages.
//   - errCh: a channel for reporting execution errors.
//   - module: the current build module used in copy context.
//   - stage: the stage definition to process.
//   - rootDir: optional root output directory prefix; if empty, use `stage.To` as-is.
//   - cb: a callback that returns a `Runnable` with `PreRun` and `PostRun` methods for the stage.
func handleStage(
	ctx context.Context,
	filesCh chan<- types.Path,
	logCh chan<- string,
	errCh chan<- error,
	module *Module,
	stage types.Stage,
	rootDir string,
	cb func(string) (interfaces.Runnable, error),
) {
	logCh <- fmt.Sprintf("Handling stage %s", stage.Name)
	defer func() {
		logCh <- fmt.Sprintf("Finished stage %s", stage.Name)
	}()

	runner, cbErr := cb(stage.Name)

	if cbErr == nil {
		if err := runner.PreRun(ctx); err != nil {
			errCh <- fmt.Errorf("pre-run callback failed for stage %s: %w", stage.Name, err)
			return
		}
	}

	if err := helpers.CheckContext(ctx); err != nil {
		return
	}

	dirPath := stage.To
	if rootDir != "" {
		dirPath = filepath.Join(rootDir, stage.To)
	}

	dirPath = filepath.Clean(dirPath)
	dirPath, err := filepath.Abs(dirPath)
	if err != nil {
		errCh <- fmt.Errorf("failed to get absolute path for stage %s: %s", stage.Name, err)
		return
	}

	to, err := fs.MkDir(dirPath)
	if err != nil {
		errCh <- fmt.Errorf("failed to make stage `to` directory: %w", err)
		return
	}

	var wg sync.WaitGroup
	for _, from := range stage.From {
		if err = helpers.CheckContext(ctx); err != nil {
			return
		}

		fromCopy := from
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err = helpers.CheckContext(ctx); err != nil {
				return
			}

			path := types.Path{
				From:           fromCopy,
				To:             to,
				ActionIfExists: stage.ActionIfFileExists,
				Convert:        stage.ConvertTo1251,
			}

			if err := fs.PathProcessing(
				ctx,
				filesCh,
				module,
				path,
				stage.Filter,
			); err != nil {
				errCh <- fmt.Errorf("failed to copy from %s to %s: %w", path.From, path.To, err)
			}
		}()
	}

	wg.Wait()

	if runner != nil {
		if err = runner.PostRun(ctx); err != nil {
			errCh <- fmt.Errorf("post-run callback failed for stage %s: %w", stage.Name, err)
		}
	}
}
