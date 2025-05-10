package queue

import (
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"sync"
)

var (
	parallel = config.Telegram.Get().MaxParallelDownload

	lock   sync.RWMutex
	queued int64
	// Telegram.MessageID => DownloadTask
	downloading map[int64]*DownloadTask
)

func init() {
	logger := logfield.New(logfield.ComQueue).WithAction("init")

	downloadRepo := repository.DownloadRepository{
		DB: database.DB,
	}

	var err error
	queued, err = downloadRepo.CountQueued()
	if err != nil {
		logger.Fatalln("count download failed:", err)
	}

	// resume downloads

	downloadingSlice, err := downloadRepo.GetDownloading()
	if err != nil {
		logger.Fatalln("load downloading failed:", err)
	}
	downloading = make(map[int64]*DownloadTask, len(downloadingSlice))

	for _, repo := range downloadingSlice {
		err := StartDownload(repo)
		if err != nil {
			logger.Errorln("resume download failed:", err)
		}
	}
}

func Clean() error {
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
