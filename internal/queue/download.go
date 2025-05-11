package queue

import (
	"errors"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
)

func CleanDownload() error {
	lock.Lock()
	defer lock.Unlock()

	downloading = make(map[int64]*DownloadTask)

	if err := source.Telegram.RemoveDownloads(); err != nil {
		return err
	}
	if err := source.Telegram.CleanDownloadDirectory(); err != nil {
		return err
	}
	return nil
}

func startDownload(model repository.Download) error {
	task, ok := downloading[model.MsgID]
	if ok {
		return task.UpdateOrDownload()
	}

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

	return errors.New("not implemented")
}
