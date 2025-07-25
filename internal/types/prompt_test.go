package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestPrompt_Input(t *testing.T) {
	t.Parallel()

	type fields struct {
		Value string
	}
	type args struct {
		validator func(string) error
		title     string
	}
	tests := []struct {
		args    args
		name    string
		fields  fields
		wantErr bool
	}{
		{args{
			validator: func(string) error { return nil },
			title:     "",
		}, "empty title", fields{}, true},
		{args{
			validator: nil,
			title:     "title",
		}, "empty title", fields{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			p := NewPrompt()
			p.Value = tt.fields.Value
			err := p.Input(tt.args.title, tt.args.validator)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewPrompt(t *testing.T) {
	p := NewPrompt()
	assert.Empty(t, p.GetValue())
}
