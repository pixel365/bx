package module

import (
	"context"
	"errors"
	"testing"

	"github.com/pixel365/bx/internal/interfaces"

	errors2 "github.com/pixel365/bx/internal/errors"
)

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
			"cancelled context",
			true,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for _, tt := range tests {
		cancel()
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fields.builder.Build(ctx)
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
