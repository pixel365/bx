package types

import "testing"

func TestPrompt_Input(t *testing.T) {
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
			p := NewPrompt()
			p.Value = tt.fields.Value
			if err := p.Input(tt.args.title, tt.args.validator); (err != nil) != tt.wantErr {
				t.Errorf("Input() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewPrompt(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		p := NewPrompt()
		if p.GetValue() != "" {
			t.Errorf("NewPrompt.Value = %v, want %v", p.GetValue(), "")
		}
	})
}
