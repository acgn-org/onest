package source

import (
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/telegram"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
)

var Telegram *telegram.Telegram

func init() {
	loggerCom := logfield.New(logfield.ComSource)
	logger := loggerCom.WithAction("init:telegram")

	// options
	var opts = make([]client.Option, 0, 2)
	if log.StandardLogger().Level != log.TraceLevel {
		_, err := client.SetLogStream(&client.SetLogStreamRequest{
			LogStream: &client.LogStreamEmpty{},
		})
		if err != nil {
			logger.Warnln("set log stream failed:", err)
		}
		opts = append(opts, client.WithLogVerbosity(&client.SetLogVerbosityLevelRequest{
			NewVerbosityLevel: 0,
		}))
	}
	proxyRequest, ok := telegram.ProxyFromEnvironment(logger)
	if ok {
		opts = append(opts, client.WithProxy(proxyRequest))
	}

	// connect client
	var err error
	Telegram, err = telegram.New(&telegram.Config{
		Logger:     loggerCom.WithSubComponent(logfield.ComTelegram),
		Version:    config.VERSION,
		DataFolder: config.Telegram.DataFolder,
		ApiId:      config.Telegram.ApiId,
		ApiHash:    config.Telegram.ApiHash,
	}, opts...)
	if err != nil {
		logger.Fatalln("failed to create Telegram client:", err)
	}
}
