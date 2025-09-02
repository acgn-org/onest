package realsearch

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/acgn-org/onest/tools"
)

func (c Client) NewProxy() *Proxy {
	u, _ := url.Parse(c.baseUrl)
	proxy := httputil.NewSingleHostReverseProxy(u)

	proxy.Transport = c.httpClient.Transport
	proxy.BufferPool = tools.BufferHttpUtil{}

	rawDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		rawDirector(req)
		req.Host = req.URL.Host
	}
	return &Proxy{
		proxy,
	}
}

type Proxy struct {
	*httputil.ReverseProxy
}
