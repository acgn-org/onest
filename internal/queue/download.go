package queue

import (
	"errors"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"time"
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

func setDownloadError(id uint, isFatal bool, msg string, date int64) error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	err := downloadRepo.UpdateDownloadError(id, isFatal, msg, date)
	if err != nil {
		return err
	}
	return downloadRepo.Commit().Error
}

func StartDownload(model repository.Download) error {
	lock.Lock()
	defer lock.Unlock()

	task, ok := downloading[model.MsgID]
	if ok { // resume download within queue
		task.lock.Lock()
		defer task.lock.Unlock()

		newState, err := source.Telegram.DownloadFile(task.VideoFile.Video.Id, task.priority)
		if err != nil {
			task.errorAt = time.Now()
			_ = setDownloadError(task.RepoID, false, err.Error(), task.errorAt.Unix())
		} else {
			task.state = newState
			task.stateUpdatedAt = time.Now()
		}
		return err
	}

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	if !model.Downloading {
		if err := downloadRepo.SetDownloading(model.ID); err != nil {
			return err
		}
	}

	task, err := retrieveAndStartDownload(model)
	if err != nil {
		_ = setDownloadError(model.ID, true, err.Error(), time.Now().Unix())
		return err
	}
	downloading[model.MsgID] = task

	return downloadRepo.Commit().Error
}

func AddDownloadQueue(repo repository.Download) error {
	// todo
	return errors.New("not implemented")
}
