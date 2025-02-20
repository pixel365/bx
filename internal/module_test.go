package internal

import (
	"context"
	"testing"
)

func TestModule_IsValid(t *testing.T) {
	type fields struct {
		Ctx            context.Context
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
			Name:           "test",
			Version:        "1.0.0",
			Account:        "tester",
			Repository:     "tester",
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
			Repository:     "tester",
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
				Name:           tt.fields.Name,
				Version:        tt.fields.Version,
				Account:        tt.fields.Account,
				Repository:     tt.fields.Repository,
				BuildDirectory: tt.fields.BuildDirectory,
				LogDirectory:   tt.fields.LogDirectory,
				Stages:         tt.fields.Stages,
				Ignore:         tt.fields.Ignore,
			}
			if err := m.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
