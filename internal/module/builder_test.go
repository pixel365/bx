package module

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/pixel365/bx/internal/interfaces"

	errors2 "github.com/pixel365/bx/internal/errors"
)

func Test_makeVersionDescription(t *testing.T) {
	defer func() {
		stat, err := os.Stat("testdata")
		if err != nil {
			return
		}

		if stat.IsDir() {
			_ = os.RemoveAll(stat.Name())
		}
	}()

	type args struct {
		builder *ModuleBuilder
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
			},
			logger: nil,
		}}, "empty repository", false},
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
				Description:    "some description",
			},
			logger: nil,
		}}, "has description", false},
		{args{builder: &ModuleBuilder{
			module: &Module{
				BuildDirectory: "testdata",
				Version:        "1.0.0",
				Repository:     ".",
			},
			logger: nil,
		}}, "has repository", false},
		{args{builder: &ModuleBuilder{module: &Module{LastVersion: true}}}, "last version", false},
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
		builder interfaces.Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
		{
			fields{builder: NewModuleBuilder(&Module{}, nil)},
			"todo context",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Build(context.TODO())
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModuleBuilder_Prepare(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder interfaces.Builder
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

			if !errors.Is(err, errors2.ErrNilModule) {
				t.Errorf("Prepare() error = %v, wantErr %v", err, errors2.ErrNilModule)
			}
		})
	}
}

func TestModuleBuilder_Cleanup(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder interfaces.Builder
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
		builder interfaces.Builder
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

			if !errors.Is(err, errors2.ErrNilModule) {
				t.Errorf("Rollback() error = %v, wantErr %v", err, errors2.ErrNilModule)
			}
		})
	}
}

func TestModuleBuilder_Collect(t *testing.T) {
	builder := NewModuleBuilder(nil, nil)
	type fields struct {
		builder interfaces.Builder
	}
	tests := []struct {
		fields  fields
		name    string
		wantErr bool
	}{
		{fields{builder: builder}, "nil module", true},
		{fields{builder: NewModuleBuilder(&Module{}, nil)}, "empty build directory", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Collect(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !errors.Is(err, errors2.ErrNilModule) {
				t.Errorf("Collect() error = %v, wantErr %v", err, errors2.ErrNilModule)
			}
		})
	}
}

func Test_makeVersionFile(t *testing.T) {
	defer func() {
		stat, err := os.Stat("testdata")
		if err != nil {
			return
		}

		if stat.IsDir() {
			_ = os.RemoveAll(stat.Name())
		}
	}()

	type args struct {
		builder *ModuleBuilder
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		{args{builder: &ModuleBuilder{module: &Module{}}}, "empty module", true},
		{args{builder: &ModuleBuilder{module: &Module{LastVersion: true}}}, "last version", false},
		{args{builder: &ModuleBuilder{module: &Module{
			BuildDirectory: "testdata",
			Version:        "1.0.0",
		}}}, "valid version", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeVersionFile(tt.args.builder); (err != nil) != tt.wantErr {
				t.Errorf("makeVersionFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
