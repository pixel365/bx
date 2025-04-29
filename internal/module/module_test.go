package module

import (
	"testing"
)

func TestModule_GetVersion(t *testing.T) {
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
			m := &Module{
				Version:     tt.fields.Version,
				LastVersion: tt.fields.LastVersion,
			}
			if got := m.GetVersion(); got != tt.want {
				t.Errorf("GetVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
