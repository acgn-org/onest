package telegram

import (
	"context"
	"errors"
	"github.com/acgn-org/onest/tools"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

// WithRetry retry if getting 429 errors
func (t Telegram) WithRetry(ctx context.Context, fn func() (err error)) error {
start:
	fnErr := fn()
	if fnErr == nil {
		return nil
	}
	var responseErr *client.ResponseError
	const frequencyErrorPrefix = "Too Many Requests: retry after "
	if errors.As(fnErr, &responseErr) &&
		responseErr.Err.Code == 429 && strings.HasPrefix(responseErr.Err.Message, frequencyErrorPrefix) {
		// decode cool down duration
		coolDown := strings.TrimPrefix(responseErr.Err.Message, frequencyErrorPrefix)
		coolDownSeconds, err := strconv.ParseInt(coolDown, 10, 64)
		if err == nil {
			t.logger.Warnf("reached telegram flood control, will recover after %d seconds", coolDownSeconds)
			select {
			case <-time.After(time.Duration(coolDownSeconds) * time.Second):
			case <-ctx.Done():
				return ctx.Err()
			}
			goto start
		} else {
			t.logger.Debugf("decode cool down duration failed: %v", err)
		}
	}
	return fnErr
}

func (t Telegram) GetListener() *client.Listener {
	return t.client.GetListener()
}

func (t Telegram) GetHistory(ctx context.Context, chatID, fromMessageID int64, limit int32) (*client.Messages, error) {
	var messages *client.Messages
	return messages, t.WithRetry(ctx, func() (err error) {
		messages, err = t.client.GetChatHistory(&client.GetChatHistoryRequest{
			ChatId:        chatID,
			FromMessageId: fromMessageID,
			Offset:        0,
			Limit:         limit,
			OnlyLocal:     false,
		})
		return err
	})
}

func (t Telegram) GetChat(ctx context.Context, id int64) (*client.Chat, error) {
	var chat *client.Chat
	return chat, t.WithRetry(ctx, func() (err error) {
		chat, err = t.client.GetChat(&client.GetChatRequest{
			ChatId: id,
		})
		return err
	})
}

func (t Telegram) GetMessage(ctx context.Context, chatId, messageId int64) (*client.Message, error) {
	var message *client.Message
	return message, t.WithRetry(ctx, func() (err error) {
		message, err = t.client.GetMessage(&client.GetMessageRequest{
			ChatId:    chatId,
			MessageId: messageId,
		})
		return err
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
	return t.client.GetFile(&client.GetFileRequest{
		FileId: fileID,
	})
}

func (t Telegram) DownloadFile(fileID, priority int32, synchronous bool) (*client.File, error) {
	t.logger.Debugf("download file %d with priotiry %d", fileID, priority)
	var file *client.File
	return file, t.WithRetry(context.Background(), func() (err error) {
		file, err = t.client.DownloadFile(&client.DownloadFileRequest{
			FileId:      fileID,
			Priority:    priority,
			Offset:      0,
			Limit:       0,
			Synchronous: synchronous,
		})
		return err
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
