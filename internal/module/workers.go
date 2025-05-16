package module

import (
	"context"
	"sync"

	"github.com/pixel365/bx/internal/interfaces"
	"github.com/pixel365/bx/internal/types"
)

// cleanupWorker coordinates the shutdown of all channels and workers.
//
// It waits for stage processing (stage worker goroutines) and copy workers to complete,
// then closes all communication channels and invokes cancellation once to signal global termination.
//
// Parameters:
//   - stagesWorkersWg: WaitGroup tracking stage processing goroutines.
//   - copyFilesWg: WaitGroup tracking file copy worker goroutines.
//   - once: sync.Once to ensure cancellation is called only once.
//   - cancel: context.CancelFunc to terminate the parent context.
//   - filesCh: channel of file copy tasks to be closed after stage completion.
//   - logCh: channel of log messages to be closed after copy workers finish.
//   - errCh: channel of errors to be closed after copy workers finish.
func cleanupWorker(
	stagesWorkersWg, copyFilesWg *sync.WaitGroup,
	once *sync.Once,
	cancel context.CancelFunc,
	filesCh chan types.Path,
	logCh chan string,
	errCh chan error,
) {
	stagesWorkersWg.Wait()
	close(filesCh)

	copyFilesWg.Wait()
	close(logCh)
	close(errCh)

	once.Do(cancel)
}

// copyWorkers launches a fixed number of worker goroutines to process file copy tasks.
//
// Each worker reads from the `filesCh` channel and invokes the `copyFileFunc` function
// to handle the file copying.
// Workers exit when the channel is closed.
//
// Parameters:
//   - ctx: context for cancellation.
//   - wg: WaitGroup used to signal completion of all worker goroutines.
//   - filesCh: channel carrying file copy tasks (type Path).
//   - errCh: channel used to report errors from the workers.
//   - workersCount: number of worker goroutines to spawn.
func copyWorkers(ctx context.Context, wg *sync.WaitGroup, filesCh chan types.Path,
	errCh chan<- error, workersCount int) {
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesCh {
				copyFileFunc(ctx, errCh, file)
			}
		}()
	}
}

// errorWorker monitors the error channel and triggers cancellation on the first non-nil error.
//
// It ensures that only the first error is captured and that cancellation is invoked only once
// using sync.Once.
// It is intended to run as a background goroutine.
//
// Parameters:
//   - ch: channel from which errors are read.
//   - cancel: function to cancel the shared context.
//   - once: sync.Once to ensure cancel is called only once.
//   - err: pointer to the shared error variable to capture the first error.
func errorWorker(ch <-chan error, cancel context.CancelFunc, once *sync.Once, err *error) {
	for e := range ch {
		if e != nil {
			once.Do(func() {
				*err = e
				cancel()
			})
		}
	}
}

// logWorker consumes log messages from the provided channel and sends them to the logger.
//
// It runs as a background goroutine, and processes log messages until the channel is closed.
//
// Parameters:
//   - ch: channel carrying log messages.
//   - logger: implementation of BuildLogger used to output logs.
func logWorker(ch <-chan string, logger interfaces.BuildLogger) {
	for msg := range ch {
		logger.Info(msg)
	}
}
