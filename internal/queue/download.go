package queue

import (
	"container/list"
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/acgn-org/onest/tools"
	"github.com/zelenin/go-tdlib/client"
	"regexp"
	"time"
)

func GetDownloading() ([]repository.DownloadTask, error) {
	lock.RLock()

	taskIds := make([]uint, 0, len(downloading))
	for k := range downloading {
		taskIds = append(taskIds, k)
	}

	lock.RUnlock()

	tasks, err := database.NewRepository[repository.DownloadRepository]().GetDownloadTaskByID(taskIds...)
	if err != nil {
		return nil, err
	}

	if len(tasks) != 0 {
		MigrateDownloadTaskInfo(tasks)
	}
	return tasks, nil
}

func MigrateDownloadTaskInfo(tasks []repository.DownloadTask) {
	lock.RLock()
	defer lock.RUnlock()

	for i, task := range tasks {
		taskQueue, ok := downloading[task.ID]
		if !ok {
			continue
		}

		state := taskQueue.state.Load()
		if state != nil {
			tasks[i].File = state.File
		}

		tasks[i].FatalError = taskQueue.log.isFatal.Load()
		errorState := taskQueue.log.error.Load()
		tasks[i].Error = errorState.Err
		tasks[i].ErrorAt = errorState.At.Unix()
	}
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

	items, err := itemRepo.GetForUpdates(int32(time.Now().Add(-time.Duration(config.Telegram.Get().ScanThresholdDays) * time.Hour * 24).Unix()))
	if err != nil {
		return 0, err
	}
	var created int
	for _, item := range items {
		logger := logger.WithField("item", item.Name)
		savepoint := fmt.Sprintf("sp%d", item.ID)
		var latest *client.Message
		var fromMessageID int64 = 0

		itemRegexp, err := regexp.Compile(item.Regexp)
		if err != nil {
			logger.Errorf("compile regexp '%s' failed: %v", item.Regexp, err)
			continue
		}

		if err := itemRepo.DB.SavePoint(savepoint).Error; err != nil {
			logger.Errorln("save transaction point failed:", err)
			continue
		}

		// fetch all new messages, list => *client.Message
		messageList := list.New()
	fetchMessage:
		messages, err := source.Telegram.GetHistory(item.ChannelID, fromMessageID)
		if err != nil {
			logger.Errorf("get chat %d history failed: %v", item.ChannelID, err)
			continue
		}
		fromMessageID = messages.Messages[0].Id
		for i := len(messages.Messages) - 1; i >= 0; i-- {
			msg := messages.Messages[i]
			if msg.Id <= item.Process {
				break
			}
			if latest == nil || latest.Id < msg.Id {
				latest = msg
			}
			messageList.PushFront(msg)
		}
		if fromMessageID > item.Process {
			goto fetchMessage
		}

		if latest == nil {
			// no new message found
			continue
		}

		// update item process
		item.Process, item.DateEnd = latest.Id, latest.Date
		if err := itemRepo.UpdateProcess(item.ID, item.Process, item.DateEnd); err != nil {
			logger.Errorln("update process failed:", err)
			itemRepo.DB.RollbackTo(savepoint)
			continue
		}

		// match messages
		el := messageList.Front()
		for el != nil {
			next := el.Next()
			msg := el.Value.(*client.Message)
			videoContent, ok := msg.Content.(*client.MessageVideo)
			if !ok || tools.ConvertPatternRegexp(videoContent.Caption.Text, itemRegexp, item.MatchPattern) != item.MatchContent {
				messageList.Remove(el)
			}
			el = next
		}

		// create download models
		if messageList.Len() > 0 {
			created += messageList.Len()
			messages := make([]*client.Message, 0, messageList.Len())
			for el := messageList.Front(); el != nil; el = el.Next() {
				messages = append(messages, el.Value.(*client.Message))
			}
			if _, err := downloadRepo.CreateWithMessages(item.ID, item.Priority, messages); err != nil {
				logger.Errorln("save download tasks to database failed:", err)
				itemRepo.DB.RollbackTo(savepoint)
				continue
			}
		}
	}

	return created, itemRepo.Commit().Error
}
