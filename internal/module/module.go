package module

import (
	"sync"

	"github.com/pixel365/bx/internal/types"

	"github.com/pixel365/bx/internal/callback"
)

type Module struct {
	Variables      map[string]string   `yaml:"variables,omitempty"`
	Run            map[string][]string `yaml:"run,omitempty"`
	changes        *types.Changes      `yaml:"-"`
	Log            *types.Log          `yaml:"log,omitempty"`
	Repository     string              `yaml:"repository,omitempty"`
	Label          types.VersionLabel  `yaml:"label,omitempty"`
	Version        string              `yaml:"version"`
	Description    string              `yaml:"description,omitempty"`
	Name           string              `yaml:"name"`
	Account        string              `yaml:"account"`
	BuildDirectory string              `yaml:"buildDirectory,omitempty"`
	Changelog      types.Changelog     `yaml:"changelog,omitempty"`
	Builds         types.Builds        `yaml:"builds"`
	Stages         []types.Stage       `yaml:"stages"`
	Ignore         []string            `yaml:"ignore"`
	Callbacks      []callback.Callback `yaml:"callbacks,omitempty"`
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
