package queue

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/zelenin/go-tdlib/client"
)

func GetDownloading() ([]repository.DownloadTask, error) {
	lock.RLock()
	defer lock.RUnlock()

	tasks := make([]repository.DownloadTask, 0, len(downloading))
	for _, task := range downloading {
		var stateEncoded json.RawMessage
		if task.state != nil {
			var err error
			stateEncoded, err = task.state.MarshalJSON()
			if err != nil {
				return nil, err
			}
		}

		tasks = append(tasks, repository.DownloadTask{
			ID:         task.ID,
			MsgID:      task.MsgID,
			Priority:   task.priority,
			FatalError: task.isFatal.Load(),
			File:       stateEncoded,

			ItemID:      0,
			Text:        "",
			Size:        0,
			Date:        0,
			Downloading: false,
			Downloaded:  false,
			Error:       "",
			ErrorAt:     0,
		})
	}

	return tasks, database.NewRepository[repository.DownloadRepository]().GetDownloadTaskInfo(tasks)
}

func clean() error {
	downloading = make(map[uint]*DownloadTask)

	if err := source.Telegram.RemoveAllDownloads(); err != nil {
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
func startDownload(channelId int64, download repository.Download) error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	if !download.Downloading {
		if err := downloadRepo.SetDownloading(download.ID); err != nil {
			return err
		}
	}
	if err := downloadRepo.Commit().Error; err != nil {
		return err
	}

	task, err := NewTask(channelId, download)
	downloading[download.ID] = task
	if err != nil {
		return err
	}

	return nil
}

func AddDownloadQueue(channelId int64, model repository.Download) error {
	lock.Lock()
	defer lock.Unlock()

	task, ok := downloading[model.ID]
	if ok {
		return task.UpdateOrDownload()
	}

	if int(config.Telegram.Get().MaxParallelDownload) <= len(downloading) {
		// skip and wait for trigger from supervisor
		return nil
	}

	return startDownload(channelId, model)
}

func ScanAndCreateNewDownloadTasks() (int, error) {
	itemRepo := database.BeginRepository[repository.ItemRepository]()
	defer itemRepo.Rollback()

	downloadRepo := repository.DownloadRepository{Repository: itemRepo.Repository}

	logger := logfield.New(logfield.ComQueue).WithAction("add downloads with message")

	items, err := itemRepo.GetAllForUpdates()
	if err != nil {
		return 0, err
	}
	var created int
	for _, item := range items {
		logger := logger.WithField("item", item.Name)
		savepoint := fmt.Sprintf("sp%d", item.ID)
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
		created += messageList.Len()
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

	return created, itemRepo.Commit().Error
}
