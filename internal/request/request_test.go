package request

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"

	module2 "github.com/pixel365/bx/internal/module"
)

type mockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestNewClient(t *testing.T) {
	t.Run("new client", func(t *testing.T) {
		client := NewClient(nil, nil)
		if client == nil {
			t.Error("nil client")
		}
	})
}

func TestClient_Authorization(t *testing.T) {
	mockClient := &mockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		resp.Header.Set("Set-Cookie", "BITRIX_SM_LOGIN=testuser")

		return resp, nil
	}}

	client := NewClient(mockClient, nil)

	want := []*http.Cookie{
		{Name: "BITRIX_SM_LOGIN", Value: "testuser"},
	}

	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		want    []*http.Cookie
		wantErr bool
	}{
		{"empty login", client, args{"", "1234556"}, nil, true},
		{"empty password", client, args{"abc", ""}, nil, true},
		{"success", client, args{"testuser", "123456"}, want, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.client.Authenticate(tt.args.login, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("Authenticate() got = %v, want %v", got, tt.want)
				return
			}

			for i := range got {
				if got[i].Name != tt.want[i].Name || got[i].Value != tt.want[i].Value {
					t.Errorf("Authenticate() got = %v, want %v", got[i], tt.want[i])
				}
			}
		})
	}
}

func TestClient_UploadZIP(t *testing.T) {
	mockClient := &mockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		resp.Header.Set("Set-Cookie", "BITRIX_SM_LOGIN=testuser")

		return resp, nil
	}}

	client := NewClient(mockClient, nil)
	cookies := []*http.Cookie{
		{Name: "BITRIX_SM_LOGIN", Value: "testuser"},
	}

	type args struct {
		module  *module2.Module
		cookies []*http.Cookie
	}
	tests := []struct {
		name    string
		client  *Client
		args    args
		wantErr bool
	}{
		{"nil module", client, args{module: nil, cookies: cookies}, true},
		{"nil cookies", client, args{module: &module2.Module{}, cookies: nil}, true},
		{"empty cookies", client, args{module: &module2.Module{}, cookies: []*http.Cookie{}}, true},
		{
			"empty module name",
			client,
			args{
				module:  &module2.Module{},
				cookies: []*http.Cookie{{Name: "foo", Value: "bar"}},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.UploadZIP(tt.args.module, tt.args.cookies); (err != nil) != tt.wantErr {
				t.Errorf("UploadZIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_UploadZIP_InvalidZipPath(t *testing.T) {
	origGetSession := getSessionFunc
	defer func() { getSessionFunc = origGetSession }()
	getSessionFunc = func(c *Client, module *module2.Module, cookies []*http.Cookie) string {
		return "fake-session-id"
	}

	mockClient := &mockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		resp.Header.Set("Set-Cookie", "BITRIX_SM_LOGIN=testuser")

		return resp, nil
	}}
	client := NewClient(mockClient, nil)

	t.Run("", func(t *testing.T) {
		err := client.UploadZIP(
			&module2.Module{Name: "fake-name"},
			[]*http.Cookie{{Name: "foo", Value: "bar"}},
		)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func TestClient_SessionId(t *testing.T) {
	mockClient := &mockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(
				strings.NewReader(
					"<input type=\"hidden\" name=\"sessid\" id=\"sessid\" value=\"123456\" />",
				),
			),
		}

		return resp, nil
	}}

	jar, _ := cookiejar.New(nil)
	client := NewClient(mockClient, jar)
	cookies := []*http.Cookie{
		{Name: "BITRIX_SM_LOGIN", Value: "testuser"},
	}

	module := &module2.Module{}
	module.Name = "test"

	type args struct {
		module  *module2.Module
		cookies []*http.Cookie
	}
	tests := []struct {
		client *Client
		name   string
		want   string
		args   args
	}{
		{client, "success", "123456", args{module: module, cookies: cookies}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.SessionId(tt.args.module, tt.args.cookies); got != tt.want {
				t.Errorf("SessionId() = %v, want %v", got, tt.want)
			}
		})
	}
}
