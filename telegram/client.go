package telegram

import (
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"path/filepath"
)

type Config struct {
	Logger  log.FieldLogger
	Version string

	DataFolder string
	ApiId      int32
	ApiHash    string
}

func New(c *Config, opts ...client.Option) (*Telegram, error) {
	authorizer := client.ClientAuthorizer(&client.SetTdlibParametersRequest{
		DatabaseDirectory:   filepath.Join(c.DataFolder, "database"),
		FilesDirectory:      filepath.Join(c.DataFolder, "files"),
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		ApiId:               c.ApiId,
		ApiHash:             c.ApiHash,
		SystemLanguageCode:  "en-US",
		DeviceModel:         "onest",
		ApplicationVersion:  c.Version,
	})
	c.Logger.Infoln("authorizing...")
	go client.CliInteractor(authorizer)
	_client, err := client.NewClient(authorizer, opts...)
	if err != nil {
		return nil, err
	}

	if c.Logger == nil {
		c.Logger = log.StandardLogger()
	}

	user, err := _client.GetMe()
	if err != nil {
		return nil, err
	}
	c.Logger.Debugf("user GetMe: %+v", user)

	return &Telegram{
		logger: c.Logger,
		client: _client,
	}, nil
}

type Telegram struct {
	logger log.FieldLogger
	client *client.Client
}
