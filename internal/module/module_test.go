package module

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pixel365/bx/internal/types"
)

func TestModule_GetVersion(t *testing.T) {
	t.Parallel()
	type fields struct {
		Version     string
		LastVersion bool
	}
	tests := []struct {
		name   string
		want   string
		fields fields
	}{
		{"release", "1.0.0", fields{"1.0.0", false}},
		{"last version", ".last_version", fields{"1.0.0", true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				Version:     tt.fields.Version,
				LastVersion: tt.fields.LastVersion,
			}
			got := m.GetVersion()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestModule_GetLabel(t *testing.T) {
	t.Parallel()
	type fields struct {
		Label types.VersionLabel
	}
	tests := []struct {
		name   string
		want   types.VersionLabel
		fields fields
	}{
		{"default", types.Alpha, fields{}},
		{"alpha", types.Alpha, fields{types.Alpha}},
		{"beta", types.Beta, fields{types.Beta}},
		{"stable", types.Stable, fields{types.Stable}},
		{"override", types.Alpha, fields{types.VersionLabel("override")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			m := &Module{
				Label: tt.fields.Label,
			}
			got := m.GetLabel()
			assert.Equal(t, tt.want, got)
		})
	}
}
