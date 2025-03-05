package internal

import (
	"net/http"
	"testing"
)

func TestCallback_IsValid(t *testing.T) {
	type fields struct {
		Stage string
		Pre   CallbackParameters
		Post  CallbackParameters
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"empty callback",
			fields{},
			true,
		},
		{
			"empty pre/post",
			fields{Stage: "stage name"},
			true,
		},
		{
			"invalid method",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "https://example.com",
					Method: http.MethodHead,
				},
			},
			true,
		},
		{
			"empty method",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "https://example.com",
					Method: "",
				},
			},
			true,
		},
		{
			"invalid action url",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "some invalid url",
					Method: http.MethodGet,
				},
			},
			true,
		},
		{
			"invalid type",
			fields{Stage: "stage name", Pre: CallbackParameters{Type: "some type"}},
			true,
		},
		{
			"empty action",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "",
					Method: http.MethodGet,
				},
			},
			true,
		},
		{
			"invalid scheme",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "ftp://example.com",
					Method: http.MethodGet,
				},
			},
			true,
		},
		{
			"empty scheme",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "example.com",
					Method: http.MethodGet,
				},
			},
			true,
		},
		{
			"valid",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "https://example.com",
					Method: http.MethodGet,
				},
			},
			false,
		},
		{
			"empty parameter",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "https://example.com",
					Method: http.MethodGet,
					Parameters: []string{
						"",
						"key=value",
					},
				},
			},
			true,
		},
		{
			"invalid parameter pair",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "https://example.com",
					Method: http.MethodGet,
					Parameters: []string{
						"key=value",
						"param",
					},
				},
			},
			true,
		},
		{
			"valid parameters",
			fields{
				Stage: "stage name",
				Pre: CallbackParameters{
					Type:   external,
					Action: "https://example.com",
					Method: http.MethodGet,
					Parameters: []string{
						"key=value",
						"key2=value2",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Callback{
				Stage: tt.fields.Stage,
				Pre:   tt.fields.Pre,
				Post:  tt.fields.Post,
			}
			if err := c.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCallbackParameters_IsValid(t *testing.T) {
	type fields struct {
		Type       string
		Action     string
		Method     string
		Parameters []string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, true},
		{"empty type", fields{Type: ""}, true},
		{"empty action", fields{Type: external, Action: ""}, true},
		{"invalid type", fields{Type: "some type"}, true},
		{
			"invalid method",
			fields{Type: external, Action: "https://example.com", Method: http.MethodHead},
			true,
		},
		{"empty method", fields{Type: external, Action: "https://example.com", Method: ""}, true},
		{
			"invalid action url",
			fields{Type: external, Action: "some invalid url", Method: http.MethodGet},
			true,
		},
		{
			"invalid action url scheme",
			fields{Type: external, Action: "ftp://example.com", Method: http.MethodGet},
			true,
		},
		{
			"valid",
			fields{
				Type:       external,
				Action:     "https://example.com",
				Method:     http.MethodGet,
				Parameters: []string{"key=value"},
			},
			false,
		},
		{
			"invalid parameters",
			fields{
				Type:       external,
				Action:     "https://example.com",
				Method:     http.MethodGet,
				Parameters: []string{"=value"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CallbackParameters{
				Type:       tt.fields.Type,
				Action:     tt.fields.Action,
				Method:     tt.fields.Method,
				Parameters: tt.fields.Parameters,
			}
			if err := c.IsValid(); (err != nil) != tt.wantErr {
				t.Errorf("IsValid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
