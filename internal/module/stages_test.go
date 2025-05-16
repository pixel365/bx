package module

import (
	"context"
	"errors"
	"fmt"
	"testing"

	errors2 "github.com/pixel365/bx/internal/errors"
	"github.com/pixel365/bx/internal/helpers"
	"github.com/pixel365/bx/internal/types"
)

func TestHandleStages_NoCustomCommandMode(t *testing.T) {
	ctx := context.Background()
	m := &Module{
		Stages: []types.Stage{
			{Name: "some-fake-stage"},
		},
	}
	t.Run("nil context", func(t *testing.T) {
		err := HandleStages(
			ctx,
			[]string{"some-fake-stage"},
			m,
			&FakeBuildLogger{},
			false,
		)
		if !errors.Is(err, errors2.ErrNilModule) {
			t.Errorf("HandleStages() error = %v, want %v", err, errors2.ErrNilModule)
		}
	})
}

func TestCheckStages(t *testing.T) {
	err := CheckStages(nil)
	if !errors.Is(err, errors2.ErrNilModule) {
		t.Errorf("CheckStages() error = %v, want %v", err, errors2.ErrNilModule)
	}
}

func TestCheckStages_NoErrors(t *testing.T) {
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
	if err != nil {
		t.Errorf("CheckStages() error = %v, want nil", err)
	}
}

func TestCheckStages_WithErrors(t *testing.T) {
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
	if err == nil {
		t.Errorf("CheckStages() error = %v, want error", err)
	} else {
		expectedMsg := "errors: [failed stage: fail]"
		if err.Error() != expectedMsg {
			t.Errorf("CheckStages() error = %v, want %v", err, expectedMsg)
		}
	}
}
