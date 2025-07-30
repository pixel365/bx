package module

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/types/changelog"

	"github.com/pixel365/bx/internal/callback"
	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"
)

var fakeCondType types.ChangelogConditionType = "fake"

func TestModule_IsValid(t *testing.T) {
	t.Parallel()
	emptyRun := make(map[string][]string)
	run := make(map[string][]string)
	run["test"] = []string{"test"}

	type fields struct {
		Variables      map[string]string
		Log            *types.Log
		Run            map[string][]string
		Label          types.VersionLabel
		Account        string
		Repository     string
		BuildDirectory string
		LogDirectory   string
		Version        string
		Name           string
		Builds         types.Builds
		Stages         []types.Stage
		Ignore         []string
		Changelog      changelog.Changelog
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
			Builds: types.Builds{
				Release: []string{"test"},
			},
			Run: run,
		}, false},
		{"invalid", fields{
			Name:           "test",
			Version:        "",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore:    []string{},
			Variables: nil,
			Run:       run,
		}, true},
		{"repository does not exists", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "repository",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
			Run:    run,
		}, true},
		{"valid sort value", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release: []string{"test"},
			},
			Run: run,
		}, false},
		{"valid stage filter", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release: []string{"test"},
			},
			Run: run,
		}, false},
		{"valid label", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Label:          types.Stable,
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release: []string{"test"},
			},
			Run: run,
		}, false},
		{"invalid label", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Label:          types.VersionLabel("invalid label"),
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release: []string{"test"},
			},
			Run: run,
		}, true},
		{"invalid stage in last version", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Label:          types.Alpha,
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release:     []string{"test"},
				LastVersion: []string{"invalid stage"},
			},
			Run: run,
		}, true},
		{"invalid stage in release", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Label:          types.Alpha,
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release:     []string{"invalid release stage"},
				LastVersion: []string{"test"},
			},
			Run: run,
		}, true},
		{"empty run", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Label:          types.Alpha,
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release:     []string{"test"},
				LastVersion: []string{"test"},
			},
			Run: emptyRun,
		}, true},
		{"invalid log", fields{
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Label:          types.Alpha,
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: changelog.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release:     []string{"test"},
				LastVersion: []string{"test"},
			},
			Run: run,
			Log: &types.Log{Dir: ""},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				Variables:      nil,
				Name:           tt.fields.Name,
				Version:        tt.fields.Version,
				Label:          tt.fields.Label,
				Account:        tt.fields.Account,
				BuildDirectory: tt.fields.BuildDirectory,
				Stages:         tt.fields.Stages,
				Ignore:         tt.fields.Ignore,
				Repository:     tt.fields.Repository,
				Builds:         tt.fields.Builds,
				Run:            tt.fields.Run,
				Log:            tt.fields.Log,
			}
			err := m.IsValid()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestModule_IsValid_EmptyName(t *testing.T) {
	t.Parallel()
	module := &Module{Name: ""}
	err := module.IsValid()
	assert.ErrorIs(t, err, errors2.ErrEmptyModuleName)
}

func TestModule_IsValid_SpacesInName(t *testing.T) {
	t.Parallel()
	module := &Module{Name: "some name"}
	err := module.IsValid()
	assert.ErrorIs(t, err, errors2.ErrNameContainsSpace)
}

func TestModule_IsValid_EmptyAccount(t *testing.T) {
	t.Parallel()
	module := &Module{Name: "name", Version: "1.0.0", Account: ""}
	err := module.IsValid()
	assert.ErrorIs(t, err, errors2.ErrEmptyAccountName)
}

func TestModule_NormalizeStages(t *testing.T) {
	t.Parallel()
	type fields struct {
		Variables      map[string]string
		Name           string
		Version        string
		Account        string
		Repository     string
		BuildDirectory string
		LogDirectory   string
		Stages         []types.Stage
		Ignore         []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{
			Variables: map[string]string{
				"foo": "bar",
			},
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "tester",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "{foo}",
					To:                 "{foo}",
					ActionIfFileExists: types.Replace,
					From:               []string{"./{foo}"},
				},
			},
			Ignore: []string{},
		}, false},
		{"invalid", fields{
			Variables: map[string]string{
				"": "bar",
			},
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "tester",
			BuildDirectory: "tester",
			Stages: []types.Stage{
				{
					Name:               "{foo}",
					To:                 "{foo}",
					ActionIfFileExists: types.Replace,
					From:               []string{"./{foo}"},
				},
			},
			Ignore: []string{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				Variables:      tt.fields.Variables,
				Name:           tt.fields.Name,
				Version:        tt.fields.Version,
				Account:        tt.fields.Account,
				BuildDirectory: tt.fields.BuildDirectory,
				Stages:         tt.fields.Stages,
				Ignore:         tt.fields.Ignore,
			}
			err := m.NormalizeStages()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestModule_PasswordEnv(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name string
	}
	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{"success without dots", "TEST_PASSWORD", fields{"test"}},
		{"success with dots", "TEST_TEST_PASSWORD", fields{"test.test"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				Name: tt.fields.Name,
			}
			got := m.PasswordEnv()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestModule_ValidateChangelog(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name           string
		Version        string
		Account        string
		BuildDirectory string
		LogDirectory   string
		Repository     string
		Stages         []types.Stage
		Changelog      changelog.Changelog
	}

	mod := fields{
		Changelog:      changelog.Changelog{},
		Name:           "test",
		Version:        "1.0.0",
		Account:        "tester",
		BuildDirectory: "tester",
		Repository:     ".",
		Stages: []types.Stage{
			{
				Name:               "test",
				To:                 "test",
				ActionIfFileExists: types.Replace,
				From:               []string{"./tes"},
			},
		},
	}

	tests := []struct {
		name    string
		fields  changelog.Changelog
		wantErr bool
	}{
		{"empty", changelog.Changelog{}, false},
		{"empty from type", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  "",
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		}, true},
		{"empty from value", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		}, true},
		{"empty to type", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  "",
				Value: "v2.0.0",
			},
		}, true},
		{"empty to value", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "",
			},
		}, true},
		{"valid without conditions", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		}, false},
		{"empty condition", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
					"",
				},
			},
		}, true},
		{"invalid regex in condition", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
					`(`,
				},
			},
		}, true},
		{"invalid changelog sort", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
				},
			},
			Sort: "sort",
		}, true},
		{"fully valid", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: types.Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
				},
			},
		}, false},
		{"invalid condition type", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type: fakeCondType,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
				},
			},
		}, true},
		{"empty condition values", changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type:  types.Include,
				Value: []string{},
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				Changelog:      tt.fields,
				Name:           mod.Name,
				Version:        mod.Version,
				Account:        mod.Account,
				BuildDirectory: mod.BuildDirectory,
				Repository:     mod.Repository,
				Stages:         mod.Stages,
			}
			err := m.ValidateChangelog()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestModule_ValidateChangelog_empty_repository(t *testing.T) {
	t.Parallel()
	m := &Module{
		Repository: "",
		Changelog: changelog.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
			Condition: types.TypeValue[types.ChangelogConditionType, []string]{
				Type:  types.Include,
				Value: []string{},
			},
		},
	}
	err := m.ValidateChangelog()
	require.NoError(t, err)
}

func TestModule_FindStage(t *testing.T) {
	t.Parallel()
	type fields struct {
		Name           string
		Version        string
		Account        string
		BuildDirectory string
		LogDirectory   string
		Stages         []types.Stage
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    types.Stage
		wantErr bool
	}{
		{"valid", fields{
			Name:           "test",
			Version:        "1.0.0",
			Account:        "test",
			BuildDirectory: "./build",
			Stages: []types.Stage{
				{
					Name:               "stage_1",
					To:                 "to",
					ActionIfFileExists: types.Replace,
					From:               []string{"from"},
				},
			},
		},
			args{name: "stage_1"},
			types.Stage{
				Name:               "stage_1",
				To:                 "to",
				ActionIfFileExists: types.Replace,
				From:               []string{"from"},
			},
			false,
		},
		{"invalid", fields{
			Name:           "test",
			Version:        "1.0.0",
			Account:        "test",
			BuildDirectory: "./build",
			Stages: []types.Stage{
				{
					Name:               "stage_1",
					To:                 "to",
					ActionIfFileExists: types.Replace,
					From:               []string{"from"},
				},
			},
		},
			args{name: "stage_2"},
			types.Stage{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				BuildDirectory: tt.fields.BuildDirectory,
				Version:        tt.fields.Version,
				Account:        tt.fields.Account,
				Name:           tt.fields.Name,
				Stages:         tt.fields.Stages,
			}
			got, err := m.FindStage(tt.args.name)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestModule_ZipPath(t *testing.T) {
	t.Parallel()
	mod := Module{}
	_, err := mod.ZipPath()
	require.Error(t, err)
}

func TestModule_StageCallback(t *testing.T) {
	t.Parallel()
	mod := Module{}
	_, err := mod.StageCallback("stage_1")
	assert.ErrorIs(t, err, errors2.ErrStageCallbackNotFound)
}

func TestModule_StageCallback_found(t *testing.T) {
	t.Parallel()
	mod := Module{
		Stages: []types.Stage{
			{
				Name:               "stage_1",
				To:                 "to",
				ActionIfFileExists: types.Replace,
				From:               []string{"from"},
			},
		},
		Callbacks: []callback.Callback{
			{
				Stage: "stage_1",
			},
		},
	}
	_, err := mod.StageCallback("stage_1")
	require.NoError(t, err)
}
