package client

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	SetCookies(u *url.URL, cookies []*http.Cookie)
}

type httpClient struct {
	c *http.Client
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	return c.c.Do(req)
}

func (c *httpClient) SetCookies(u *url.URL, cookies []*http.Cookie) {
	c.c.Jar.SetCookies(u, cookies)
}

func NewClient(ttl time.Duration) HTTPClient {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:     jar,
		Timeout: ttl,
	}
	return &httpClient{
		c: client,
	}
}
