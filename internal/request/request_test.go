package request

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	client2 "github.com/pixel365/bx/internal/client"

	module2 "github.com/pixel365/bx/internal/module"
)

func Test_Authorization(t *testing.T) {
	client := &client2.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		resp.Header.Set("Set-Cookie", "BITRIX_SM_LOGIN=testuser")

		return resp, nil
	}}

	want := []*http.Cookie{
		{Name: "BITRIX_SM_LOGIN", Value: "testuser"},
	}

	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name    string
		client  client2.HTTPClient
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
			got, err := Authenticate(tt.client, tt.args.login, tt.args.password)
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

func Test_UploadZIP(t *testing.T) {
	ctx := context.Background()
	client := &client2.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		resp.Header.Set("Set-Cookie", "BITRIX_SM_LOGIN=testuser")

		return resp, nil
	}}

	cookies := []*http.Cookie{
		{Name: "BITRIX_SM_LOGIN", Value: "testuser"},
	}

	type args struct {
		module  *module2.Module
		cookies []*http.Cookie
	}
	tests := []struct {
		name    string
		client  client2.HTTPClient
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
			if err := UploadZIP(ctx, tt.client, tt.args.module, tt.args.cookies); (err != nil) != tt.wantErr {
				t.Errorf("UploadZIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_UploadZIP_InvalidZipPath(t *testing.T) {
	ctx := context.Background()
	origGetSession := getSessionFunc
	defer func() { getSessionFunc = origGetSession }()
	getSessionFunc = func(c client2.HTTPClient, module *module2.Module, cookies []*http.Cookie) string {
		return "fake-session-id"
	}

	client := &client2.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader("")),
		}
		resp.Header.Set("Set-Cookie", "BITRIX_SM_LOGIN=testuser")

		return resp, nil
	}}

	t.Run("", func(t *testing.T) {
		err := UploadZIP(
			ctx,
			client,
			&module2.Module{Name: "fake-name"},
			[]*http.Cookie{{Name: "foo", Value: "bar"}},
		)
		if err == nil {
			t.Error("expected error")
		}
	})
}

func Test_SessionId(t *testing.T) {
	client := &client2.MockHttpClient{DoFunc: func(req *http.Request) (*http.Response, error) {
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
		client client2.HTTPClient
		name   string
		want   string
		args   args
	}{
		{client, "success", "123456", args{module: module, cookies: cookies}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sessionId(tt.client, tt.args.module, tt.args.cookies); got != tt.want {
				t.Errorf("sessionId() = %v, want %v", got, tt.want)
			}
		})
	}
}
