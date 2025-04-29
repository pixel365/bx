package interfaces

import (
	"context"
	"sync"

	"github.com/pixel365/bx/internal/types"
)

type Builder interface {
	Build() error
	Prepare() error
	Rollback() error
	Collect() error
	Cleanup()
}

type ModuleConfig interface {
	GetVariables() map[string]string
	GetRun() map[string][]string
	GetStages() []types.Stage
	GetIgnore() []string
	GetChanges() *types.Changes
	IsLastVersion() bool
}

type BuildLogger interface {
	Info(message string, args ...interface{})
	Error(message string, err error, args ...interface{})
	Cleanup()
}

type Runnable interface {
	PreRun(ctx context.Context, wg *sync.WaitGroup, logger BuildLogger)
	PostRun(ctx context.Context, wg *sync.WaitGroup, logger BuildLogger)
}
