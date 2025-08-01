package callback

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
					Type:   ExternalType,
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
		{
			"invalid post",
			fields{
				Stage: "stage name",
				Post: CallbackParameters{
					Type: ExternalType,
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Callback{
				Stage: tt.fields.Stage,
				Pre:   tt.fields.Pre,
				Post:  tt.fields.Post,
			}
			err := c.IsValid()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
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
		{"empty action", fields{Type: ExternalType, Action: ""}, true},
		{"invalid type", fields{Type: "some type"}, true},
		{
			"invalid method",
			fields{Type: ExternalType, Action: "https://example.com", Method: http.MethodHead},
			true,
		},
		{
			"empty method",
			fields{Type: ExternalType, Action: "https://example.com", Method: ""},
			true,
		},
		{
			"invalid action url",
			fields{Type: ExternalType, Action: "some invalid url", Method: http.MethodGet},
			true,
		},
		{
			"invalid action url scheme",
			fields{Type: ExternalType, Action: "ftp://example.com", Method: http.MethodGet},
			true,
		},
		{
			"valid",
			fields{
				Type:       ExternalType,
				Action:     "https://example.com",
				Method:     http.MethodGet,
				Parameters: []string{"key=value"},
			},
			false,
		},
		{
			"invalid parameters",
			fields{
				Type:       ExternalType,
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
			err := c.IsValid()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCallbackParameters_buildUrlAndBody(t *testing.T) {
	type fields struct {
		Type       string
		Action     string
		Method     string
		Parameters []string
	}
	tests := []struct {
		want1  io.Reader
		name   string
		want   string
		fields fields
	}{
		{
			nil,
			"valid get",
			"https://example.com?key=value",
			fields{
				Type:       ExternalType,
				Action:     "https://example.com",
				Method:     http.MethodGet,
				Parameters: []string{"key=value"},
			},
		},
		{
			nil,
			"valid get without params",
			"https://example.com",
			fields{
				Type:   ExternalType,
				Action: "https://example.com",
				Method: http.MethodGet,
			},
		},
		{
			strings.NewReader(url.Values{"key": []string{"value"}}.Encode()),
			"valid post",
			"https://example.com",
			fields{
				Type:       ExternalType,
				Action:     "https://example.com",
				Method:     http.MethodPost,
				Parameters: []string{"key=value"},
			},
		},
		{
			nil,
			"valid post without params",
			"https://example.com",
			fields{
				Type:   ExternalType,
				Action: "https://example.com",
				Method: http.MethodPost,
			},
		},
		{
			nil,
			"valid command",
			"ls -lsa",
			fields{
				Type:       CommandType,
				Action:     "ls",
				Parameters: []string{"-lsa"},
			},
		},
		{
			nil,
			"valid command without params",
			"ls",
			fields{
				Type:   CommandType,
				Action: "ls",
			},
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
			got, got1 := c.buildUrlAndBody()
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestValidateCallbacks(t *testing.T) {
	tests := []struct {
		name    string
		cb      []Callback
		wantErr bool
	}{
		{"empty callback list", []Callback{}, false},
		{
			"empty callback stage",
			[]Callback{{Stage: ""}},
			true,
		},
		{
			"empty callback pre/post type",
			[]Callback{{Stage: "testing"}},
			true,
		},
		// TODO: we need more tests
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCallbacks(tt.cb); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCallbacks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInvalidCallbackParametersRun(t *testing.T) {
	ctx := context.TODO()
	cbp := CallbackParameters{}
	err := cbp.Run(ctx)
	require.Error(t, err)
}

func TestCallback_PreRun(t *testing.T) {
	ctx := context.TODO()
	cb := Callback{}
	err := cb.PreRun(ctx)
	require.Error(t, err)
}

func TestCallback_PostRun(t *testing.T) {
	ctx := context.TODO()
	cb := Callback{}
	err := cb.PostRun(ctx)
	require.Error(t, err)
}
