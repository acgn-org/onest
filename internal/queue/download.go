package queue

import (
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
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
		queued++
		return nil
	}

	return startDownload(model)
}
