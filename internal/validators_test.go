package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidateModuleName_NotExisting(t *testing.T) {
	t.Run("TestValidateModuleName_NotExisting", func(t *testing.T) {
		if err := ValidateModuleName("not_exists", "./"); err != nil {
			t.Error(err)
		}
	})
}

func TestValidateModuleName_Existing(t *testing.T) {
	t.Run("TestValidateModuleName_Existing", func(t *testing.T) {
		name := fmt.Sprintf("%s_%d", "testing", time.Now().Unix())
		filePath, err := filepath.Abs(fmt.Sprintf("%s/%s.yaml", ".", name))
		if err != nil {
			t.Error()
		}

		err = os.WriteFile(filePath, []byte(DefaultYAML()), 0600)
		if err != nil {
			t.Error(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Error(err)
			}
		}(filePath)

		err = ValidateModuleName(name, ".")
		if err == nil {
			t.Errorf("error expected")
		}
	})
}

func TestValidateVersion(t *testing.T) {
	type args struct {
		version string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "1", args: args{version: "1.0.0"}, wantErr: false},
		{name: "2", args: args{version: "v1.0.0"}, wantErr: true},
		{name: "3", args: args{version: "3.0.10"}, wantErr: false},
		{name: "4", args: args{version: ""}, wantErr: true},
		{name: "5", args: args{version: "some version"}, wantErr: true},
		{name: "6", args: args{version: "111.000.123"}, wantErr: false},
		{name: "7", args: args{version: "111.00x0.123"}, wantErr: true},
		{name: "8", args: args{version: "111.00x0.123"}, wantErr: true},
		{name: "9", args: args{version: "1x11.00x0.123"}, wantErr: true},
		{name: "10", args: args{version: "1x11.00x0.123x"}, wantErr: true},
		{name: "11", args: args{version: "1..1.1"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateVersion(tt.args.version); (err != nil) != tt.wantErr {
				t.Errorf("ValidateVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"valid", args{password: "123456"}, false},
		{"empty", args{password: ""}, true},
		{"only spaces", args{password: "    "}, true},
		{"short", args{password: "123"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePassword(tt.args.password); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
			if got := ValidateArgument(tt.args.arg); got != tt.want {
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
		{args{m: &Module{Stages: []Stage{{Name: ""}}}}, "empty stage name", true},
		{args{m: &Module{Stages: []Stage{{Name: "testing"}}}}, "empty stage to", true},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing", To: "testing", ActionIfFileExists: ""}},
				},
			},
			"empty stage action",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: ReplaceIfNewer,
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
					Stages: []Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: ReplaceIfNewer,
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
					Stages: []Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: ReplaceIfNewer,
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
					Stages: []Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: ReplaceIfNewer,
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
					Stages: []Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: ReplaceIfNewer,
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
					Stages: []Stage{
						{
							Name:               "testing",
							To:                 "testing",
							ActionIfFileExists: ReplaceIfNewer,
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
			if err := ValidateStages(tt.args.m); (err != nil) != tt.wantErr {
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
		{args{m: &Module{Variables: map[string]string{"": "value"}}}, "empty variable key", true},
		{args{m: &Module{Variables: map[string]string{"key": ""}}}, "empty variable value", true},
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

func TestValidateCallbacks(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{m: &Module{}}, "empty callback list", false},
		{args{m: &Module{Callbacks: []Callback{{Stage: ""}}}}, "empty callback stage", true},
		{
			args{m: &Module{Callbacks: []Callback{{Stage: "testing"}}}},
			"empty callback pre/post type",
			true,
		},
		// TODO: we need more tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCallbacks(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCallbacks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBuilds(t *testing.T) {
	type args struct {
		m *Module
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{m: &Module{}}, "empty build list", true},
		{
			args{m: &Module{Builds: Builds{Release: []string{""}}}},
			"empty stage in release list",
			true,
		},
		{
			args{m: &Module{Builds: Builds{Release: []string{"testing", "testing"}}}},
			"duplicate stage in release list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"test"}},
				},
			},
			"invalid stage in release list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}},
				},
			},
			"valid stage in release list",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}},
				},
			},
			"empty lastVersion list",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}, LastVersion: []string{""}},
				},
			},
			"empty stage in lastVersion list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{
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
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}, LastVersion: []string{"testing"}},
				},
			},
			"valid stage in lastVersion list",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}, LastVersion: []string{"test"}},
				},
			},
			"invalid stage in lastVersion list",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateBuilds(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateBuilds() error = %v, wantErr %v", err, tt.wantErr)
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
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}},
				},
			},
			"empty lastVersion list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}, LastVersion: []string{""}},
				},
			},
			"empty stage in lastVersion list",
			true,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{
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
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}, LastVersion: []string{"testing"}},
				},
			},
			"valid stage in lastVersion list",
			false,
		},
		{
			args{
				m: &Module{
					Stages: []Stage{{Name: "testing"}},
					Builds: Builds{Release: []string{"testing"}, LastVersion: []string{"test"}},
				},
			},
			"invalid stage in lastVersion list",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateLastVersion(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateLastVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateStagesInBuilds(t *testing.T) {
	m := &Module{Stages: []Stage{{Name: "testing"}}}
	type args struct {
		find   func(string) (Stage, error)
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
			Stages: []Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {}}},
		}, "empty stage values", true},
		{args{m: &Module{
			Stages: []Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {"unknown"}}},
		}, "unknown stage", true},
		{args{m: &Module{
			Stages: []Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {"testing", "testing"}}},
		}, "duplicated stage", true},
		{args{m: &Module{
			Stages: []Stage{{Name: "testing"}},
			Run:    map[string][]string{"testing": {"testing"}}},
		}, "valid", false},
		{args{m: &Module{
			Stages: []Stage{{Name: "testing"}},
			Run:    map[string][]string{"": {"testing"}}},
		}, "empty key", true},
		{args{m: &Module{
			Stages: []Stage{{Name: "testing"}},
			Run:    map[string][]string{"some key": {"testing"}}},
		}, "key with spaces", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateRun(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("ValidateRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
