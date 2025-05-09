package telegram

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"golang.org/x/net/http/httpproxy"
	"net/url"
	"strconv"
)

func ProxyFromEnvironment(logger log.FieldLogger) (*client.AddProxyRequest, bool) {
	env := httpproxy.FromEnvironment()
	var proxyUrlStr string
	if env.HTTPProxy != "" {
		proxyUrlStr = env.HTTPProxy
	} else {
		proxyUrlStr = env.HTTPSProxy
	}
	if proxyUrlStr == "" {
		return nil, false
	}

	proxyUrl, err := url.Parse(proxyUrlStr)
	if err != nil {
		logger.WithError(err).WithField("url", proxyUrlStr).Errorln("invalid proxy url")
		return nil, false
	} else if proxyUrl.Scheme != "http" {
		logger.Warnf("unsupported scheme '%s'", proxyUrl.Scheme)
		return nil, false
	}

	var portStr = proxyUrl.Port()
	var port int32
	if portStr == "" {
		port = 80
	} else {
		port64, err := strconv.ParseInt(portStr, 10, 32)
		if err != nil {
			logger.WithError(err).WithField("port", portStr).Errorln("invalid proxy port")
			return nil, false
		}
		port = int32(port64)
	}

	passwd, _ := proxyUrl.User.Password()

	return &client.AddProxyRequest{
		Enable: true,
		Server: proxyUrl.Hostname(),
		Port:   port,
		Type: &client.ProxyTypeHttp{
			Username: proxyUrl.User.Username(),
			Password: passwd,
			HttpOnly: false,
		},
	}, true
}
