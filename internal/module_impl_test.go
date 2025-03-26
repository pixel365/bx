package internal

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

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
		Stages         []Stage
		Ignore         []string
		Changelog      Changelog
		Builds         Builds
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
			Stages: []Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
			Builds: Builds{
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
			Stages: []Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: Replace,
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
			Stages: []Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: Replace,
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
			Stages: []Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: Replace,
					From:               []string{"./tester"},
				},
			},
			Ignore: []string{},
			Changelog: Changelog{
				Sort: Asc,
			},
			Builds: Builds{
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
			Stages: []Stage{
				{
					Name:               "test",
					To:                 "tester",
					ActionIfFileExists: Replace,
					From:               []string{"./tester"},
					Filter:             []string{"**/*.php"},
				},
			},
			Ignore: []string{},
			Changelog: Changelog{
				Sort: Asc,
			},
			Builds: Builds{
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
		Stages         []Stage
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
			Stages: []Stage{
				{
					Name:               "{foo}",
					To:                 "{foo}",
					ActionIfFileExists: Replace,
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
			Stages: []Stage{
				{
					Name:               "{foo}",
					To:                 "{foo}",
					ActionIfFileExists: Replace,
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
		Changelog      Changelog
		Name           string
		Version        string
		Account        string
		BuildDirectory string
		LogDirectory   string
		Repository     string
		Stages         []Stage
	}

	mod := fields{
		Ctx:            context.TODO(),
		Changelog:      Changelog{},
		Name:           "test",
		Version:        "1.0.0",
		Account:        "tester",
		BuildDirectory: "tester",
		LogDirectory:   "tester",
		Repository:     ".",
		Stages: []Stage{
			{
				Name:               "test",
				To:                 "test",
				ActionIfFileExists: Replace,
				From:               []string{"./tes"},
			},
		},
	}

	tests := []struct {
		name    string
		fields  Changelog
		wantErr bool
	}{
		{"empty", Changelog{}, true},
		{"empty from type", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  "",
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
		}, true},
		{"empty from value", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
		}, true},
		{"empty to type", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  "",
				Value: "v2.0.0",
			},
		}, true},
		{"empty to value", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "",
			},
		}, true},
		{"valid without conditions", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
		}, false},
		{"empty condition", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
			Condition: TypeValue[ChangelogConditionType, []string]{
				Type: Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
					"",
				},
			},
		}, true},
		{"invalid regex in condition", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
			Condition: TypeValue[ChangelogConditionType, []string]{
				Type: Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
					`(`,
				},
			},
		}, true},
		{"invalid changelog sort", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
			Condition: TypeValue[ChangelogConditionType, []string]{
				Type: Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
				},
			},
			Sort: "sort",
		}, true},
		{"fully valid", Changelog{
			From: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v1.0.0",
			},
			To: TypeValue[ChangelogType, string]{
				Type:  Tag,
				Value: "v2.0.0",
			},
			Condition: TypeValue[ChangelogConditionType, []string]{
				Type: Include,
				Value: []string{
					`^feat: ([\W\w]+)$`,
					`^fix: ([\W\w]+)$`,
				},
			},
		}, false},
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

func TestModule_FindStage(t *testing.T) {
	type fields struct {
		Ctx            context.Context
		Name           string
		Version        string
		Account        string
		BuildDirectory string
		LogDirectory   string
		Stages         []Stage
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Stage
		wantErr bool
	}{
		{"valid", fields{
			Ctx:            context.TODO(),
			Name:           "test",
			Version:        "1.0.0",
			Account:        "test",
			BuildDirectory: "./build",
			LogDirectory:   "./logs",
			Stages: []Stage{
				{
					Name:               "stage_1",
					To:                 "to",
					ActionIfFileExists: Replace,
					From:               []string{"from"},
				},
			},
		},
			args{name: "stage_1"},
			Stage{
				Name:               "stage_1",
				To:                 "to",
				ActionIfFileExists: Replace,
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
			Stages: []Stage{
				{
					Name:               "stage_1",
					To:                 "to",
					ActionIfFileExists: Replace,
					From:               []string{"from"},
				},
			},
		},
			args{name: "stage_2"},
			Stage{},
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
	if !errors.Is(err, StageCallbackNotFoundError) {
		t.Errorf("StageCallback() error = %v, wantErr %v", err, StageCallbackNotFoundError)
	}
}

func TestModule_GetChanges(t *testing.T) {
	mod := Module{}
	changes := mod.GetChanges()
	if changes != nil {
		t.Errorf("GetChanges() error = %v, wantErr %v", changes, nil)
	}
}
