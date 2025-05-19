package client

import (
	"net/http"
	"net/url"
)

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (m *MockHttpClient) SetCookies(_ *url.URL, _ []*http.Cookie) {}
