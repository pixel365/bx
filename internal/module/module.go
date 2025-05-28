package module

import (
	"sync"

	"github.com/pixel365/bx/internal/types/changelog"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/callback"
)

type Module struct {
	Variables      map[string]string   `yaml:"variables,omitempty"`
	Run            map[string][]string `yaml:"run,omitempty"`
	changes        *types.Changes      `yaml:"-"`
	Log            *types.Log          `yaml:"log,omitempty"`
	Name           string              `yaml:"name"`
	Version        string              `yaml:"version"`
	Description    string              `yaml:"description,omitempty"`
	Repository     string              `yaml:"repository,omitempty"`
	Account        string              `yaml:"account"`
	BuildDirectory string              `yaml:"buildDirectory,omitempty"`
	Label          types.VersionLabel  `yaml:"label,omitempty"`
	Builds         types.Builds        `yaml:"builds"`
	Ignore         []string            `yaml:"ignore"`
	Stages         []types.Stage       `yaml:"stages"`
	Callbacks      []callback.Callback `yaml:"callbacks,omitempty"`
	Changelog      changelog.Changelog `yaml:"changelog,omitempty"`
	mu             sync.Mutex          `yaml:"-"`
	LastVersion    bool                `yaml:"-"`
}

func (m *Module) GetVersion() string {
	if m.LastVersion {
		return ".last_version"
	}
	return m.Version
}

func (m *Module) GetLabel() types.VersionLabel {
	switch m.Label {
	case types.Alpha, types.Beta, types.Stable:
		return m.Label
	default:
		return types.Alpha
	}
}
