package module

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pixel365/bx/internal/types"
)

func TestCopyWorkers(t *testing.T) {
	var mu sync.Mutex
	var called []types.Path

	copyFileFunc = func(ctx context.Context, errCh chan<- error, path types.Path) {
		mu.Lock()
		called = append(called, path)
		mu.Unlock()
	}

	filesCh := make(chan types.Path, 3)
	errCh := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	copyWorkers(ctx, &wg, filesCh, errCh, 2)

	filesCh <- types.Path{From: "a.txt", To: "x"}
	filesCh <- types.Path{From: "b.txt", To: "y"}
	close(filesCh)

	wg.Wait()

	assert.Len(t, called, 2)
}

func TestErrorWorker(t *testing.T) {
	errCh := make(chan error, 2)
	var once sync.Once
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	var capturedErr error
	go errorWorker(errCh, cancel, &once, &capturedErr)

	expectedErr := errors.New("fail")
	errCh <- expectedErr
	time.Sleep(50 * time.Millisecond)

	assert.ErrorIs(t, capturedErr, expectedErr)
}

func TestLogWorker(t *testing.T) {
	logCh := make(chan string, 2)
	mock := &FakeBuildLogger{}

	go logWorker(logCh, mock)
	logCh <- "hello"
	logCh <- "world"
	close(logCh)

	time.Sleep(50 * time.Millisecond)

	assert.Len(t, mock.Logs, 2)
	assert.Equal(t, "hello", mock.Logs[0])
	assert.Equal(t, "world", mock.Logs[1])
}

func TestCleanupWorker(t *testing.T) {
	var stageWg, copyWg sync.WaitGroup
	stageWg.Add(1)
	copyWg.Add(1)

	filesCh := make(chan types.Path, 1)
	logCh := make(chan string, 1)
	errCh := make(chan error, 1)

	var canceled bool
	cancel := func() { canceled = true }

	var once sync.Once
	go cleanupWorker(&stageWg, &copyWg, &once, cancel, filesCh, logCh, errCh)

	stageWg.Done()
	copyWg.Done()

	time.Sleep(50 * time.Millisecond)

	select {
	case _, ok := <-filesCh:
		if ok {
			t.Error("filesCh should be closed")
		}
	default:
		t.Error("filesCh not closed")
	}

	assert.True(t, canceled)
}
