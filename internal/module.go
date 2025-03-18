package internal

import (
	"context"
	"sync"
)

type FileExistsAction string
type ChangelogType string
type ChangelogConditionType string
type ChangelogSort string
type BuildType string

const (
	Replace        FileExistsAction = "replace"
	Skip           FileExistsAction = "skip"
	ReplaceIfNewer FileExistsAction = "replace_if_newer"

	Commit ChangelogType = "commit"
	Tag    ChangelogType = "tag"

	Include ChangelogConditionType = "include"
	Exclude ChangelogConditionType = "exclude"

	Asc  ChangelogSort = "asc"
	Desc ChangelogSort = "desc"

	Release     BuildType = "release"
	LastVersion BuildType = "lastVersion"
)

type TypeValue[T1 any, T2 any] struct {
	Type  T1 `yaml:"type"`
	Value T2 `yaml:"value"`
}

type Changelog struct {
	From      TypeValue[ChangelogType, string]            `yaml:"from"`
	To        TypeValue[ChangelogType, string]            `yaml:"to"`
	Sort      ChangelogSort                               `yaml:"sort,omitempty"`
	Condition TypeValue[ChangelogConditionType, []string] `yaml:"condition,omitempty"`
}

type Stage struct {
	Name               string           `yaml:"name"`
	To                 string           `yaml:"to"`
	ActionIfFileExists FileExistsAction `yaml:"actionIfFileExists"`
	From               []string         `yaml:"from"`
	ConvertTo1251      bool             `yaml:"convertTo1251,omitempty"`
}

type Builds struct {
	Release     []string `yaml:"release"`
	LastVersion []string `yaml:"lastVersion,omitempty"`
}

type Module struct {
	Ctx            context.Context     `yaml:"-"`
	Variables      map[string]string   `yaml:"variables,omitempty"`
	Run            map[string][]string `yaml:"run,omitempty"`
	changes        *Changes            `yaml:"-"`
	Repository     string              `yaml:"repository,omitempty"`
	Account        string              `yaml:"account"`
	BuildDirectory string              `yaml:"buildDirectory,omitempty"`
	LogDirectory   string              `yaml:"logDirectory,omitempty"`
	Version        string              `yaml:"version"`
	Name           string              `yaml:"name"`
	Changelog      Changelog           `yaml:"changelog,omitempty"`
	Builds         Builds              `yaml:"builds"`
	Stages         []Stage             `yaml:"stages"`
	Ignore         []string            `yaml:"ignore"`
	Callbacks      []Callback          `yaml:"callbacks,omitempty"`
	LastVersion    bool                `yaml:"-"`
	mu             sync.Mutex          `yaml:"-"`
}

func (m *Module) GetVersion() string {
	if m.LastVersion {
		return ".last_version"
	}
	return m.Version
}
