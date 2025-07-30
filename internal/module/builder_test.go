package module

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/interfaces"

	errors2 "github.com/pixel365/bx/internal/errors"
)

func TestNewModuleBuilder(t *testing.T) {
	t.Parallel()
	builder := NewModuleBuilder(nil, nil)
	assert.NotNil(t, builder)
}

func TestModuleBuilder_Build(t *testing.T) {
	t.Parallel()
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
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.fields.builder.Build(ctx)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}

	t.Cleanup(func() {
		cancel()
	})
}

func TestModuleBuilder_Prepare(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			err := tt.fields.builder.Prepare()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.ErrorIs(t, err, errors2.ErrNilModule)
		})
	}
}

func TestModuleBuilder_Cleanup(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			assert.NotPanics(t, func() {
				tt.fields.builder.Cleanup()
			})
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
	t.Parallel()
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
			t.Parallel()
			err := tt.fields.builder.Collect(context.Background())
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.ErrorIs(t, err, errors2.ErrNilModule)
		})
	}
}
