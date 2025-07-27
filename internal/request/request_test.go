package request

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pixel365/bx/internal/types"

	client2 "github.com/pixel365/bx/internal/client"

	module2 "github.com/pixel365/bx/internal/module"
)

func Test_Authorization(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got, err := Authenticate(tt.client, tt.args.login, tt.args.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			for i := range got {
				assert.Equal(t, tt.want[i].Name, got[i].Name)
				assert.Equal(t, tt.want[i].Value, got[i].Value)
			}
		})
	}
}

func Test_UploadZIP(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			err := UploadZIP(ctx, tt.client, tt.args.module, tt.args.cookies)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_UploadZIP_InvalidZipPath(t *testing.T) {
	t.Parallel()

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

	err := UploadZIP(
		ctx,
		client,
		&module2.Module{Name: "fake-name"},
		[]*http.Cookie{{Name: "foo", Value: "bar"}},
	)
	require.Error(t, err)
}

func Test_SessionId(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			got := sessionId(tt.client, tt.args.module, tt.args.cookies)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVersions(t *testing.T) {
	t.Parallel()

	type args struct {
		client  client2.HTTPClient
		module  *module2.Module
		cookies []*http.Cookie
	}
	tests := []struct {
		want    types.Versions
		name    string
		args    args
		wantErr bool
	}{
		{nil, "nil module", args{client: nil, module: nil}, true},
		{
			nil,
			"nil cookies",
			args{client: nil, module: &module2.Module{}, cookies: nil},
			true,
		},
	}

	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Versions(ctx, tt.args.client, tt.args.module, tt.args.cookies)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
