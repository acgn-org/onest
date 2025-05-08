package telegram

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"golang.org/x/net/http/httpproxy"
	"net/url"
	"strconv"
)

func ProxyFromEnvironment(logger log.FieldLogger) *client.AddProxyRequest {
	var option client.AddProxyRequest

	env := httpproxy.FromEnvironment()
	var proxyUrlStr string
	if env.HTTPSProxy != "" {
		proxyUrlStr = env.HTTPSProxy
	} else {
		proxyUrlStr = env.HTTPProxy
	}
	if proxyUrlStr == "" {
		return &option
	}

	proxyUrl, err := url.Parse(proxyUrlStr)
	if err != nil {
		logger.WithError(err).WithField("url", proxyUrlStr).Errorln("invalid proxy url")
		return &option
	} else if proxyUrl.Scheme != "http" {
		logger.Warnf("unsupported scheme '%s'", proxyUrl.Scheme)
		return &option
	}

	var portStr = proxyUrl.Port()
	var port int32
	if portStr == "" {
		option.Port = 80
	} else {
		port64, err := strconv.ParseInt(portStr, 10, 32)
		if err != nil {
			logger.WithError(err).WithField("port", portStr).Errorln("invalid proxy port")
			return &option
		}
		port = int32(port64)
	}

	option.Enable = true
	option.Server = proxyUrl.Hostname()
	option.Port = port

	passwd, _ := proxyUrl.User.Password()
	option.Type = &client.ProxyTypeHttp{
		Username: proxyUrl.User.Username(),
		Password: passwd,
		HttpOnly: false,
	}

	return &option
}
