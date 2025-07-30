package module

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pixel365/bx/internal/types/changelog"

	"github.com/pixel365/bx/internal/types"
)

func TestModule_GetVariables(t *testing.T) {
	t.Parallel()
	mod := Module{
		Variables: map[string]string{
			"foo": "bar",
		},
	}
	val := mod.GetVariables()
	assert.Len(t, val, 1)
}

func TestModule_GetRun(t *testing.T) {
	t.Parallel()
	mod := Module{
		Run: map[string][]string{
			"run1": {"stage"},
		},
	}
	val := mod.GetRun()
	assert.Len(t, val, 1)
}

func TestModule_GetStages(t *testing.T) {
	t.Parallel()
	var s []types.Stage
	s = append(s, types.Stage{})

	mod := Module{Stages: s}
	val := mod.GetStages()
	assert.Len(t, val, 1)
}

func TestModule_GetIgnore(t *testing.T) {
	t.Parallel()
	mod := Module{
		Ignore: []string{"ignore"},
	}
	val := mod.GetIgnore()
	assert.Len(t, val, 1)
}

func TestModule_GetChanges(t *testing.T) {
	t.Parallel()
	mod := Module{
		Repository: "../../",
	}
	changes := mod.GetChanges()
	assert.Nil(t, changes)
}

func TestModule_GetChanges2(t *testing.T) {
	t.Parallel()
	mod := Module{
		Repository: "../../",
	}

	origChangesListFunc := changesListFunc
	defer func() {
		changesListFunc = origChangesListFunc
	}()

	changesListFunc = func(_ string, _ changelog.Changelog) (*types.Changes, error) {
		return &types.Changes{}, nil
	}

	changes := mod.GetChanges()
	assert.NotNil(t, changes)
}

func TestGetChanges_empty_repository(t *testing.T) {
	t.Parallel()
	mod := Module{}
	changes := mod.GetChanges()
	assert.Nil(t, changes)
}

func TestModule_IsLastVersion(t *testing.T) {
	t.Parallel()
	mod := Module{LastVersion: true}
	assert.True(t, mod.IsLastVersion())
}

func TestModule_SourceCount(t *testing.T) {
	t.Parallel()
	mod := Module{
		Stages: []types.Stage{
			{
				From: []string{"stage"},
			},
		},
	}

	count := mod.SourceCount()
	assert.Equal(t, 1, count)
}
