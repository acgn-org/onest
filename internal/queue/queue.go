package queue

import (
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/repository"
	"sync"
)

var (
	lock sync.RWMutex
	// Telegram.MessageID => DownloadTask
	downloading map[uint]*DownloadTask
)

func init() {
	logger := logfield.New(logfield.ComQueue).WithAction("init")

	downloadRepo := database.NewRepository[repository.DownloadRepository]()

	// resume downloads

	downloadingSlice, err := downloadRepo.GetDownloading()
	if err != nil {
		logger.Fatalln("load downloading failed:", err)
	}
	downloading = make(map[uint]*DownloadTask, len(downloadingSlice))

	for _, repo := range downloadingSlice {
		err := startDownload(repo)
		if err != nil {
			logger.Errorln("resume download failed:", err)
		}
	}

	// run supervisor

	supervisor()
}

func UpdatePriority(id uint, priority int32) {
	lock.Lock()
	defer lock.Unlock()

	download, ok := downloading[id]
	if !ok {
		return
	}
	download.priority = priority
	err := download.UpdateOrDownload()
	if err != nil {
		logfield.New(logfield.ComQueue).WithAction("update").Warnf("update priority of task %d to %d failed: %v", id, priority, err)
	}
}

func RemoveTasks(ids ...uint) {
	lock.Lock()
	defer lock.Unlock()

	for _, id := range ids {
		task, ok := downloading[id]
		if !ok {
			continue
		}

		if err := task.Terminate(); err != nil {
			logfield.New(logfield.ComQueue).WithAction("remove").Warnf("terminate task %d with error: %v", id, err)
		}
		delete(downloading, id)
	}
}
