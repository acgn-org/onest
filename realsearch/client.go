package realsearch

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewClient(c *Config) (*Client, error) {
	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	u, err := url.Parse(c.BaseUrl)
	if err != nil {
		return nil, err
	}
	u.Path = "/api/public/"
	u.RawQuery = ""
	u.Fragment = ""

	return &Client{
		baseUrl:    u.String(),
		httpClient: c.HttpClient,
	}, nil
}

type Config struct {
	HttpClient *http.Client
	BaseUrl    string
}

type Client struct {
	baseUrl    string
	httpClient *http.Client
}

func (c Client) NewProxy() *Proxy {
	u, _ := url.Parse(c.baseUrl)
	return &Proxy{
		httputil.NewSingleHostReverseProxy(u),
	}
}
