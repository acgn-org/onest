package config

import (
	"net"
	"net/http"
	"time"
)

type _RealSearch struct {
	Timeout uint8  `yaml:"timeout"`
	BaseUrl string `yaml:"base_url"`
}

var RealSearch = LoadScoped("realsearch", &_RealSearch{
	Timeout: 30,
	BaseUrl: "https://search.acgn.es/",
})

var RealSearchHttpClient *http.Client

func init() {
	timeout := time.Duration(RealSearch.Get().Timeout) * time.Second

	RealSearchHttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: timeout,
			}).DialContext,
			TLSHandshakeTimeout: timeout,
			IdleConnTimeout:     time.Minute,
		},
		Timeout: timeout,
	}
}
