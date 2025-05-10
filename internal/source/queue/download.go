package queue

import (
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"time"
)

func StartDownload(repo repository.Download) error {
	lock.Lock()
	defer lock.Unlock()

	task, ok := downloading[repo.MsgID]
	if ok { // resume download within queue
		task.lock.Lock()
		defer task.lock.Unlock()

		var err error
		task.state, err = source.Telegram.DownloadFile(task.VideoFile.Video.Id, task.priority)
		if err != nil {
			task.errorAt = time.Now()
			_ = setDownloadError(task.RepoID, false, err.Error(), task.errorAt.Unix())
		}
		return err
	}

	task, err := retrieveDownloadInfo(repo)
	if err != nil {
		_ = setDownloadError(repo.ID, true, err.Error(), time.Now().Unix())
		return err
	}
	downloading[repo.MsgID] = task

	return nil
}
