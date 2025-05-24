package telegram

import (
	"context"
	"github.com/acgn-org/onest/tools"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"golang.org/x/time/rate"
	"path"
	"path/filepath"
	"time"
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
		rate:              rate.NewLimiter(rate.Every(time.Second)/2, 3),
		databaseDirectory: databaseDirectory,
		filesDirectory:    filesDirectory,
	}, nil
}

type Telegram struct {
	logger log.FieldLogger
	client *client.Client
	rate   *rate.Limiter

	databaseDirectory string
	filesDirectory    string
}

func (t Telegram) GetListener() *client.Listener {
	return t.client.GetListener()
}

func (t Telegram) GetHistory(chatID, fromMessageID int64) (*client.Messages, error) {
	_ = t.rate.Wait(context.Background())
	return t.client.GetChatHistory(&client.GetChatHistoryRequest{
		ChatId:        chatID,
		FromMessageId: fromMessageID,
		Offset:        0,
		Limit:         99,
		OnlyLocal:     false,
	})
}

func (t Telegram) GetChat(id int64) (*client.Chat, error) {
	_ = t.rate.Wait(context.Background())
	return t.client.GetChat(&client.GetChatRequest{
		ChatId: id,
	})
}

func (t Telegram) GetMessage(chatId, messageId int64) (*client.Message, error) {
	_ = t.rate.Wait(context.Background())
	return t.client.GetMessage(&client.GetMessageRequest{
		ChatId:    chatId,
		MessageId: messageId,
	})
}

func (t Telegram) GetMessageVideo(msg *client.Message) (*client.MessageVideo, bool) {
	msgVideo, ok := msg.Content.(*client.MessageVideo)
	if !ok {
		return nil, false
	}
	return msgVideo, true
}

func (t Telegram) GetFile(fileID int32) (*client.File, error) {
	_ = t.rate.Wait(context.Background())
	return t.client.GetFile(&client.GetFileRequest{
		FileId: fileID,
	})
}

func (t Telegram) DownloadFile(fileID, priority int32, synchronous bool) (*client.File, error) {
	_ = t.rate.Wait(context.Background())
	t.logger.Debugf("download file %d with priotiry %d", fileID, priority)
	return t.client.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileID,
		Priority:    priority,
		Offset:      0,
		Limit:       0,
		Synchronous: synchronous,
	})
}

func (t Telegram) CancelDownloadFile(fileID int32) error {
	_, err := t.client.CancelDownloadFile(&client.CancelDownloadFileRequest{
		FileId: fileID,
	})
	return err
}

func (t Telegram) RemoveFileFromDownloads(fileID int32) error {
	_, err := t.client.RemoveFileFromDownloads(&client.RemoveFileFromDownloadsRequest{
		FileId:          fileID,
		DeleteFromCache: true,
	})
	return err
}

func (t Telegram) RemoveAllDownloads() error {
	_ = t.rate.Wait(context.Background())
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
	if err := tools.CleanDirectory(path.Join(t.filesDirectory, "videos")); err != nil {
		return err
	}
	if err := tools.CleanDirectory(path.Join(t.filesDirectory, "temp")); err != nil {
		return err
	}

	t.logger.Debugln("removed all video and temp files in download directory")

	return nil
}
