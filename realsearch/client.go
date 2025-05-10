package realsearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func NewClient(c *Config) (*Client, error) {
	if c.HttpClient == nil {
		c.HttpClient = http.DefaultClient
	}

	u, err := url.Parse(c.BaseUrl)
	if err != nil {
		return nil, err
	} else if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid url: %s", c.BaseUrl)
	}
	u.Path = path.Join(u.Path, "/api/public")
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

type ResponseSuccess struct {
	Data interface{} `json:"data"`
}

type ResponseError struct {
	Code uint8  `json:"code"`
	Msg  string `json:"msg"`
}

func (resp ResponseError) Error() string {
	return fmt.Sprintf("real search api error: code: %d, msg: %s", resp.Code, resp.Msg)
}

func (c Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.baseUrl+path, body)
}

func (c Client) Do(req *http.Request, data interface{}) error {
	httpRes, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer httpRes.Body.Close()

	if httpRes.StatusCode > 299 {
		var errResp ResponseError
		err := json.NewDecoder(httpRes.Body).Decode(&errResp)
		if err == nil {
			return &errResp
		}
		return fmt.Errorf("real search api internal error with http status %d", httpRes.StatusCode)
	}

	var res = ResponseSuccess{
		Data: data,
	}
	return json.NewDecoder(httpRes.Body).Decode(&res)
}
