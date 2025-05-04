package realsearch

import "net/http/httputil"

type Proxy struct {
	*httputil.ReverseProxy
}
