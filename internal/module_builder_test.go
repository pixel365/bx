package internal

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func Test_makeZipFilePath(t *testing.T) {
	mod1 := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	mod2 := &Module{
		BuildDirectory: "testdata/build",
		Version:        "1.0.1",
	}

	cur, _ := os.Getwd()

	type args struct {
		module *Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{mod1}, fmt.Sprintf("%s/testdata/1.0.0.zip", cur), false},
		{"2", args{mod2}, fmt.Sprintf("%s/testdata/build/1.0.1.zip", cur), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeZipFilePath(tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeZipFilePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeZipFilePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeVersionDirectory(t *testing.T) {
	mod1 := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	mod2 := &Module{
		BuildDirectory: "testdata/build",
		Version:        "1.0.1",
	}

	cur, _ := os.Getwd()

	type args struct {
		module *Module
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"1", args{mod1}, fmt.Sprintf("%s/testdata/1.0.0", cur), false},
		{"2", args{mod2}, fmt.Sprintf("%s/testdata/build/1.0.1", cur), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeVersionDirectory(tt.args.module)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeVersionDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("makeVersionDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeVersionDescription(t *testing.T) {
	mod := &Module{
		BuildDirectory: "testdata",
		Version:        "1.0.0",
	}

	builder := &ModuleBuilder{
		module: mod,
		logger: nil,
	}

	type args struct {
		builder *ModuleBuilder
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{builder: builder}, "empty repository", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeVersionDescription(tt.args.builder); (err != nil) != tt.wantErr {
				t.Errorf("makeVersionDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewModuleBuilder(t *testing.T) {
	t.Run("new builder", func(t *testing.T) {
		builder := NewModuleBuilder(nil, nil)
		if builder == nil {
			t.Error("NewModuleBuilder() should not be nil")
		}
	})
}

func TestModuleBuilder_Build(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, NilModuleError) {
				t.Errorf("Build() error = %v, wantErr %v", err, NilModuleError)
			}
		})
	}
}

func TestModuleBuilder_Prepare(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Prepare()
			if (err != nil) != tt.wantErr {
				t.Errorf("Prepare() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, NilModuleError) {
				t.Errorf("Prepare() error = %v, wantErr %v", err, NilModuleError)
			}
		})
	}
}

func TestModuleBuilder_Cleanup(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.builder.Cleanup()
		})
	}
}

func TestModuleBuilder_Rollback(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Rollback()
			if (err != nil) != tt.wantErr {
				t.Errorf("Rollback() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, NilModuleError) {
				t.Errorf("Rollback() error = %v, wantErr %v", err, NilModuleError)
			}
		})
	}
}

func TestModuleBuilder_Collect(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Collect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, NilModuleError) {
				t.Errorf("Collect() error = %v, wantErr %v", err, NilModuleError)
			}
		})
	}
}
