package module

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/interfaces"

	errors2 "github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/types"
)

func TestHandleStages_NoCustomCommandMode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := &Module{
		Stages: []types.Stage{
			{Name: "some-fake-stage"},
		},
	}

	err := HandleStages(
		ctx,
		[]string{"some-fake-stage"},
		m,
		&FakeBuildLogger{},
		false,
	)
	assert.ErrorIs(t, errors2.ErrNilModule, err)
}

func TestHandleStages(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := &Module{
		Stages: []types.Stage{
			{Name: "some-fake-stage"},
		},
	}

	handleStageFuncOrig := handleStageFunc
	handleStageFunc = func(ctx context.Context, filesCh chan<- types.Path, logCh chan<- string, errCh chan<- error,
		module *Module, stage types.Stage, rootDir string, cb func(string) (interfaces.Runnable, error)) {
	}
	defer func() {
		handleStageFunc = handleStageFuncOrig
	}()

	err := HandleStages(ctx, []string{""}, m, nil, true)
	require.NoError(t, err)
}

func TestCheckStages(t *testing.T) {
	t.Parallel()
	err := CheckStages(nil)
	assert.ErrorIs(t, errors2.ErrNilModule, err)
}

func TestCheckStages_NoErrors(t *testing.T) {
	t.Parallel()
	originalCheckPaths := helpers.CheckPaths
	checkPathsFunc = func(stage types.Stage, errCh chan<- error) {}
	defer func() { checkPathsFunc = originalCheckPaths }()

	m := &Module{
		Stages: []types.Stage{
			{Name: "stage1"},
			{Name: "stage2"},
		},
	}

	err := CheckStages(m)
	require.NoError(t, err)
}

func TestCheckStages_WithErrors(t *testing.T) {
	t.Parallel()
	originalCheckPaths := helpers.CheckPaths
	checkPathsFunc = func(stage types.Stage, errCh chan<- error) {
		if stage.Name == "fail" {
			errCh <- fmt.Errorf("failed stage: %s", stage.Name)
		}
	}
	defer func() { checkPathsFunc = originalCheckPaths }()

	m := &Module{
		Stages: []types.Stage{
			{Name: "ok"},
			{Name: "fail"},
		},
	}

	err := CheckStages(m)
	require.Error(t, err)

	expectedMsg := "errors: [failed stage: fail]"
	assert.Equal(t, expectedMsg, err.Error())
}

func Test_workersQty(t *testing.T) {
	t.Parallel()
	cnt := runtime.NumCPU() * 2

	assert.Equal(t, cnt, workersQty(cnt))

	n := cnt * 2
	assert.Equal(t, n, workersQty(n))
}
