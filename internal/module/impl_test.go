package module

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/pixel365/bx/internal/callback"
	errors2 "github.com/pixel365/bx/internal/errors"

	"github.com/pixel365/bx/internal/types"
)

var fakeCondType types.ChangelogConditionType = "fake"

func TestModule_IsValid(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		Variables      map[string]string
		Name           string
		Version        string
		Account        string
		Repository     string
		BuildDirectory string
		LogDirectory   string
		Stages         []types.Stage
		Ignore         []string
		Changelog      types.Changelog
		Builds         types.Builds
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{
			Ctx:            context.Background(),
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
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
		}, false},
		{"invalid", fields{
			Ctx:            context.Background(),
			Name:           "test",
			Version:        "",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
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
		}, true},
		{"repository does not exists", fields{
			Ctx:            context.Background(),
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "repository",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
		}, true},
		{"valid sort value", fields{
			Ctx:            context.Background(),
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
			Stages: []types.Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: types.Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
			Changelog: types.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release: []string{"test"},
			},
		}, false},
		{"valid stage filter", fields{
			Ctx:            context.Background(),
			Variables:      nil,
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
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
			Changelog: types.Changelog{
				Sort: types.Asc,
			},
			Builds: types.Builds{
				Release: []string{"test"},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Module{
				Ctx:            tt.fields.Ctx,
				Variables:      nil,
				Name:           tt.fields.Name,
				Version:        tt.fields.Version,
				Account:        tt.fields.Account,
				BuildDirectory: tt.fields.BuildDirectory,
				LogDirectory:   tt.fields.LogDirectory,
				Stages:         tt.fields.Stages,
				Ignore:         tt.fields.Ignore,
				Repository:     tt.fields.Repository,
				Builds:         tt.fields.Builds,
			}
			if err := m.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModule_IsValid_EmptyName(t *testing.T) {
	module := &Module{Name: ""}
	t.Run("empty name", func(t *testing.T) {
		err := module.IsValid()
		if !errors.Is(err, errors2.EmptyModuleNameError) {
			t.Errorf("IsValid() error = %v, wantErr %v", err, errors2.EmptyModuleNameError)
		}
	})
}

func TestModule_IsValid_SpacesInName(t *testing.T) {
	module := &Module{Name: "some name"}
	t.Run("spaces in name", func(t *testing.T) {
		err := module.IsValid()
		if !errors.Is(err, errors2.NameContainsSpaceError) {
			t.Errorf("IsValid() error = %v, wantErr %v", err, errors2.NameContainsSpaceError)
		}
	})
}

func TestModule_IsValid_EmptyAccount(t *testing.T) {
	module := &Module{Name: "name", Version: "1.0.0", Account: ""}
	t.Run("empty account", func(t *testing.T) {
		err := module.IsValid()
		if !errors.Is(err, errors2.EmptyAccountNameError) {
			t.Errorf("IsValid() error = %v, wantErr %v", err, errors2.EmptyAccountNameError)
		}
	})
}

func TestModule_NormalizeStages(t *testing.T) {
	type fields struct {
		Ctx            context.Context
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
			Ctx: context.Background(),
			Variables: map[string]string{
				"foo": "bar",
			},
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "tester",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
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
			Ctx: context.Background(),
			Variables: map[string]string{
				"": "bar",
			},
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "tester",
			BuildDirectory: "tester",
			LogDirectory:   "tester",
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
			m := &Module{
				Ctx:            tt.fields.Ctx,
				Variables:      tt.fields.Variables,
				Name:           tt.fields.Name,
				Version:        tt.fields.Version,
				Account:        tt.fields.Account,
				BuildDirectory: tt.fields.BuildDirectory,
				LogDirectory:   tt.fields.LogDirectory,
				Stages:         tt.fields.Stages,
				Ignore:         tt.fields.Ignore,
			}
			if err := m.NormalizeStages(); (err != nil) != tt.wantErr {
				t.Errorf("NormalizeStages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModule_PasswordEnv(t *testing.T) {
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
			m := &Module{
				Name: tt.fields.Name,
			}
			if got := m.PasswordEnv(); got != tt.want {
				t.Errorf("PasswordEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_ValidateChangelog(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		Changelog      types.Changelog
		Name           string
		Version        string
		Account        string
		BuildDirectory string
		LogDirectory   string
		Repository     string
		Stages         []types.Stage
	}

	mod := fields{
		Ctx:            context.TODO(),
		Changelog:      types.Changelog{},
		Name:           "test",
		Version:        "1.0.0",
		Account:        "tester",
		BuildDirectory: "tester",
		LogDirectory:   "tester",
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
		fields  types.Changelog
		wantErr bool
	}{
		{"empty", types.Changelog{}, true},
		{"empty from type", types.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  "",
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		}, true},
		{"empty from value", types.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		}, true},
		{"empty to type", types.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  "",
				Value: "v2.0.0",
			},
		}, true},
		{"empty to value", types.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "",
			},
		}, true},
		{"valid without conditions", types.Changelog{
			From: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v1.0.0",
			},
			To: types.TypeValue[types.ChangelogType, string]{
				Type:  types.Tag,
				Value: "v2.0.0",
			},
		}, false},
		{"empty condition", types.Changelog{
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
		{"invalid regex in condition", types.Changelog{
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
		{"invalid changelog sort", types.Changelog{
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
		{"fully valid", types.Changelog{
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
		{"invalid condition type", types.Changelog{
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
		{"empty condition values", types.Changelog{
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
			m := &Module{
				Ctx:            mod.Ctx,
				Changelog:      tt.fields,
				Name:           mod.Name,
				Version:        mod.Version,
				Account:        mod.Account,
				BuildDirectory: mod.BuildDirectory,
				LogDirectory:   mod.LogDirectory,
				Repository:     mod.Repository,
				Stages:         mod.Stages,
			}
			if err := m.ValidateChangelog(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateChangelog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModule_ValidateChangelog_empty_repository(t *testing.T) {
	t.Run("empty repository", func(t *testing.T) {
		m := &Module{
			Repository: "",
			Changelog: types.Changelog{
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
		if !errors.Is(err, errors2.InvalidChangelogSettingsError) {
			t.Errorf(
				"ValidateChangelog() error = %v, wantErr %v",
				err,
				errors2.InvalidChangelogSettingsError,
			)
		}
	})
}

func TestModule_FindStage(t *testing.T) {
	type fields struct {
		Ctx            context.Context
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
			Ctx:            context.TODO(),
			Name:           "test",
			Version:        "1.0.0",
			Account:        "test",
			BuildDirectory: "./build",
			LogDirectory:   "./logs",
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
			Ctx:            context.TODO(),
			Name:           "test",
			Version:        "1.0.0",
			Account:        "test",
			BuildDirectory: "./build",
			LogDirectory:   "./logs",
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
			m := &Module{
				Ctx:            tt.fields.Ctx,
				BuildDirectory: tt.fields.BuildDirectory,
				Version:        tt.fields.Version,
				Account:        tt.fields.Account,
				Name:           tt.fields.Name,
				LogDirectory:   tt.fields.LogDirectory,
				Stages:         tt.fields.Stages,
			}
			got, err := m.FindStage(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindStage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindStage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestModule_ZipPath(t *testing.T) {
	mod := Module{}
	_, err := mod.ZipPath()
	if err == nil {
		t.Errorf("ZipPath() error = %v, wantErr %v", err, nil)
	}
}

func TestModule_StageCallback(t *testing.T) {
	mod := Module{}
	_, err := mod.StageCallback("stage_1")
	if !errors.Is(err, errors2.StageCallbackNotFoundError) {
		t.Errorf("StageCallback() error = %v, wantErr %v", err, errors2.StageCallbackNotFoundError)
	}
}

func TestModule_StageCallback_found(t *testing.T) {
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
	if err != nil {
		t.Errorf("StageCallback() error = %v, wantErr %v", err, nil)
	}
}
