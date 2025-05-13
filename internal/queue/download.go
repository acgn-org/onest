package queue

import (
	"container/list"
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/zelenin/go-tdlib/client"
)

func clean() error {
	downloading = make(map[int64]*DownloadTask)

	if err := source.Telegram.RemoveDownloads(); err != nil {
		return err
	}
	if err := source.Telegram.CleanDownloadDirectory(); err != nil {
		return err
	}
	return nil
}

func CleanDownload() error {
	lock.Lock()
	defer lock.Unlock()
	return clean()
}

// create task and start
func startDownload(model repository.Download) error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	if !model.Downloading {
		if err := downloadRepo.SetDownloading(model.ID); err != nil {
			return err
		}
	}

	task, err := NewTask(model)
	downloading[model.MsgID] = task
	if err != nil {
		return err
	}

	return downloadRepo.Commit().Error
}

func AddDownloadQueue(model repository.Download) error {
	lock.Lock()
	defer lock.Unlock()

	task, ok := downloading[model.MsgID]
	if ok {
		return task.UpdateOrDownload()
	}

	if int(config.Telegram.Get().MaxParallelDownload) <= len(downloading) {
		// skip and wait for trigger from supervisor
		return nil
	}

	return startDownload(model)
}

func ScanAndCreateNewDownloadTasks() error {
	itemRepo := database.BeginRepository[repository.ItemRepository]()
	defer itemRepo.Rollback()

	downloadRepo := repository.DownloadRepository{Repository: itemRepo.Repository}

	logger := logfield.New(logfield.ComQueue).WithAction("add downloads with message")

	items, err := itemRepo.GetAllForUpdates()
	if err != nil {
		return err
	}
	for _, item := range items {
		logger := logger.WithField("item", item.Name)
		savepoint := fmt.Sprintf("%d", item.ID)
		var latest *client.Message
		var fromMessageID int64 = 0

		if err := itemRepo.DB.SavePoint(savepoint).Error; err != nil {
			logger.Errorln("save transaction point failed:", err)
			continue
		}

		// fetch all new messages, list => []]client.Message
		messageList := list.New()
	fetchMessage:
		messages, err := source.Telegram.GetHistory(item.ChannelID, fromMessageID)
		if err != nil {
			logger.Errorf("get chat %d history failed: %v", item.ID, err)
			continue
		}
		fromMessageID = messages.Messages[0].Id
		for i, msg := range messages.Messages {
			if msg.Id <= item.Process {
				messages.Messages = messages.Messages[i+1:]
				break
			}
			if latest == nil || latest.Id < msg.Id {
				latest = msg
			}
		}
		messageList.PushFront(messages.Messages)
		if fromMessageID > item.Process {
			goto fetchMessage
		}

		if latest == nil {
			// no new message found
			continue
		}

		// update item
		item.Process, item.DateEnd = latest.Id, latest.Date
		if err := itemRepo.UpdateProcess(item.ID, item.Process, item.DateEnd); err != nil {
			logger.Errorln("update process failed:", err)
			itemRepo.DB.RollbackTo(savepoint)
			continue
		}

		// create download models
		el := messageList.Front()
	createDownloadTask:
		if _, err := downloadRepo.CreateWithMessages(item.ID, item.Priority, el.Value.([]*client.Message)); err != nil {
			logger.Errorln("save download tasks to database failed:", err)
			itemRepo.DB.RollbackTo(savepoint)
			continue
		}
		if el.Next() != nil {
			el = el.Next()
			goto createDownloadTask
		}
	}

	return itemRepo.Commit().Error
}
