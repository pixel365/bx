package module

import (
	"github.com/pixel365/bx/internal/repo"
	"github.com/pixel365/bx/internal/types"
)

var changesListFunc = repo.ChangesList

func (m *Module) GetVariables() map[string]string {
	return m.Variables
}

func (m *Module) GetRun() map[string][]string {
	return m.Run
}

func (m *Module) GetStages() []types.Stage {
	return m.Stages
}

func (m *Module) GetIgnore() []string {
	return m.Ignore
}

func (m *Module) GetChanges() *types.Changes {
	if m.Repository == "" {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if m.changes == nil {
		changes, err := changesListFunc(m.Repository, m.Changelog)
		if err != nil {
			return nil
		}

		m.changes = changes
	}

	return m.changes
}

func (m *Module) IsLastVersion() bool {
	return m.LastVersion
}

func (m *Module) SourceCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	i := 0
	for _, stage := range m.Stages {
		i += len(stage.From)
	}

	return i
}
