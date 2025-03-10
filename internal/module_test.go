package internal

import (
	"context"
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
