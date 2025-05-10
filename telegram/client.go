package telegram

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"os"
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
	var databaseDirectory = filepath.Join(c.DataFolder, "database")
	var filesDirectory = filepath.Join(c.DataFolder, "files")

	authorizer := client.ClientAuthorizer(&client.SetTdlibParametersRequest{
		DatabaseDirectory:   databaseDirectory,
		FilesDirectory:      filesDirectory,
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
		logger:            c.Logger,
		client:            _client,
		databaseDirectory: databaseDirectory,
		filesDirectory:    filesDirectory,
	}, nil
}

type Telegram struct {
	logger log.FieldLogger
	client *client.Client

	databaseDirectory string
	filesDirectory    string
}

func (t Telegram) RemoveDownloads() error {
	_, err := t.client.RemoveAllFilesFromDownloads(&client.RemoveAllFilesFromDownloadsRequest{
		OnlyActive:      false,
		OnlyCompleted:   false,
		DeleteFromCache: true,
	})
	if err != nil {
		return err
	}

	t.logger.Debugln("removed all downloads")

	return nil
}

func (t Telegram) CleanDownloadDirectory() error {
	dirContents, err := os.ReadDir(t.filesDirectory)
	if err != nil {
		return fmt.Errorf("could not read directory contents: %w", err)
	}
	for _, entry := range dirContents {
		pathname := filepath.Join(t.filesDirectory, entry.Name())
		err = os.RemoveAll(pathname)
		if err != nil {
			return fmt.Errorf("could not remove file '%s': %w", pathname, err)
		}
	}

	t.logger.Debugln("removed all files in download directory")

	return nil
}
