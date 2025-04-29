package module

import (
	"testing"

	"github.com/pixel365/bx/internal/types"
)

func TestModule_GetVariables(t *testing.T) {
	mod := Module{
		Variables: map[string]string{
			"foo": "bar",
		},
	}
	val := mod.GetVariables()

	if len(val) != 1 {
		t.Error("GetVariables should return a single variable")
	}
}

func TestModule_GetRun(t *testing.T) {
	mod := Module{
		Run: map[string][]string{
			"run1": {"stage"},
		},
	}
	val := mod.GetRun()
	if len(val) != 1 {
		t.Error("GetRun should return a single variable")
	}
}

func TestModule_GetStages(t *testing.T) {
	var s []types.Stage
	s = append(s, types.Stage{})

	mod := Module{
		Stages: s,
	}
	val := mod.GetStages()
	if len(val) != 1 {
		t.Error("GetStages should return a single variable")
	}
}

func TestModule_GetIgnore(t *testing.T) {
	mod := Module{
		Ignore: []string{"ignore"},
	}
	val := mod.GetIgnore()
	if len(val) != 1 {
		t.Error("GetIgnore should return a single variable")
	}
}

func TestModule_GetChanges(t *testing.T) {
	mod := Module{
		Repository: "../../",
	}
	changes := mod.GetChanges()
	if changes != nil {
		t.Errorf("GetChanges() error = %v, wantErr %v", changes, nil)
	}
}

func TestModule_IsLastVersion(t *testing.T) {
	mod := Module{
		LastVersion: true,
	}
	if !mod.IsLastVersion() {
		t.Error("IsLastVersion should return true when last version is true")
	}
}
