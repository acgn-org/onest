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
	logger := logfield.New(logfield.ComSource).WithAction("init:telegram")

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
	var telegramConfig = config.Telegram.Get()
	var err error
	Telegram, err = telegram.New(&telegram.Config{
		Logger:     logfield.New(logfield.ComTelegram),
		Version:    config.VERSION,
		DataFolder: telegramConfig.DataFolder,
		ApiId:      telegramConfig.ApiId,
		ApiHash:    telegramConfig.ApiHash,
	}, opts...)
	if err != nil {
		logger.Fatalln("failed to create Telegram client:", err)
	}
}
