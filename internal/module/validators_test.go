package module

import (
	"testing"

	"github.com/pixel365/bx/internal/validators"

	"github.com/pixel365/bx/internal/types"
)

func TestValidateArgument(t *testing.T) {
	type args struct {
		arg string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty argument", args{""}, false},
		{"invalid argument", args{"*"}, false},
		{"valid argument", args{"--name"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validators.ValidateArgument(tt.args.arg); got != tt.want {
				t.Errorf("ValidateArgument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateStages(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{m: &Module{}}, "empty stages", true},
		{args{m: &Module{Stages: []types.Stage{{Name: ""}}}}, "empty stage name", true},
		{
			args{m: &Module{Stages: []types.Stage{{Name: "testing"}}}},
			"empty stage to",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{Name: "testing", To: "testing", ActionIfFileExists: ""},
					},
				},
			},
			"empty stage action",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: types.ReplaceIfNewer,
							From:               []string{""},
						},
					},
				},
			},
			"empty stage from",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: types.ReplaceIfNewer,
							From:               []string{"testing"},
						},
					},
				},
			},
			"valid stages",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: types.ReplaceIfNewer,
							From:               []string{"testing"},
							Filter:             nil,
						},
					},
				},
			},
			"nil filter",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: types.ReplaceIfNewer,
							From:               []string{"testing"},
							Filter:             []string{},
						},
					},
				},
			},
			"empty filter",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: types.ReplaceIfNewer,
							From:               []string{"testing"},
							Filter:             []string{""},
						},
					},
				},
			},
			"empty filter value",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: types.ReplaceIfNewer,
							From:               []string{"testing"},
							Filter:             []string{"**/*.php"},
						},
					},
				},
			},
			"not empty filter value",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateStages(tt.args.m.Stages); (err != nil) != tt.wantErr {
				t.Errorf("ValidateStages() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateIgnore(t *testing.T) {
	type args struct {
		m []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty ignore list", args{m: []string{}}, false},
		{"empty ignore value", args{m: []string{""}}, true},
		{"valid list", args{m: []string{"**/*.log"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRules(tt.args.m, "ignore"); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVariables(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{m: &Module{}}, "empty variables list", false},
		{
			args{m: &Module{Variables: map[string]string{"": "value"}}},
			"empty variable key",
			true,
		},
		{
			args{m: &Module{Variables: map[string]string{"key": ""}}},
			"empty variable value",
			true,
		},
		{args{m: &Module{Variables: map[string]string{"key": "value"}}}, "valid", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateVariables(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateVariables() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRelease(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{Release: []string{}},
				},
			},
			"empty",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{Release: []string{"testing", ""}},
				},
			},
			"empty stage",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{
						Release: []string{"testing", "testing"},
					},
				},
			},
			"duplicate stage",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{
						Release: []string{"testing"},
					},
				},
			},
			"valid stage",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{
						Release: []string{"test"},
					},
				},
			},
			"invalid stage",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRelease(tt.args.m.Builds.Release, tt.args.m.FindStage); (err != nil) != tt.wantErr {
				t.Errorf("ValidateLastVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateLastVersion(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{Release: []string{"testing"}},
				},
			},
			"empty lastVersion list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{Release: []string{"testing"}, LastVersion: []string{""}},
				},
			},
			"empty stage in lastVersion list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{
						Release:     []string{"testing"},
						LastVersion: []string{"testing", "testing"},
					},
				},
			},
			"duplicate stage in lastVersion list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{
						Release:     []string{"testing"},
						LastVersion: []string{"testing"},
					},
				},
			},
			"valid stage in lastVersion list",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []types.Stage{{Name: "testing"}},
					Builds: types.Builds{
						Release:     []string{"testing"},
						LastVersion: []string{"test"},
					},
				},
			},
			"invalid stage in lastVersion list",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateLastVersion(tt.args.m.Builds.LastVersion, tt.args.m.FindStage); (err != nil) != tt.wantErr {
				t.Errorf("ValidateLastVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateStagesInBuilds(t *testing.T) {
	m := &Module{Stages: []types.Stage{{Name: "testing"}}}
	type args struct {
		find   func(string) (types.Stage, error)
		name   string
		stages []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty stages",
			args{
				find:   m.FindStage,
				name:   "release",
				stages: m.Builds.Release,
			}, true},
		{"empty stage",
			args{
				find:   m.FindStage,
				name:   "release",
				stages: []string{""},
			}, true},
		{"valid stage",
			args{
				find:   m.FindStage,
				name:   "release",
				stages: []string{"testing"},
			}, false},
		{"duplicate stage",
			args{
				find:   m.FindStage,
				name:   "release",
				stages: []string{"testing", "testing"},
			}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateStagesList(tt.args.stages, tt.args.name, tt.args.find); (err != nil) != tt.wantErr {
				t.Errorf("validateStagesList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRun(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{m: &Module{}}, "empty run", false},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {}}},
		}, "empty stage values", true},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {"unknown"}}},
		}, "unknown stage", true},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {"testing", "testing"}}},
		}, "duplicated stage", true},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {"testing"}}},
		}, "valid", false},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{"": {"testing"}}},
		}, "empty key", true},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{"some key": {"testing"}}},
		}, "key with spaces", true},
		{args{m: &Module{
			Stages: []types.Stage{{Name: "testing"}},
			Run:    map[string][]string{}},
		}, "empty run", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRun(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
