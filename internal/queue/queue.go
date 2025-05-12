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
	downloading map[int64]*DownloadTask
)

func init() {
	logger := logfield.New(logfield.ComQueue).WithAction("init")

	downloadRepo := database.NewRepository[repository.DownloadRepository]()

	// resume downloads

	downloadingSlice, err := downloadRepo.GetDownloading()
	if err != nil {
		logger.Fatalln("load downloading failed:", err)
	}
	downloading = make(map[int64]*DownloadTask, len(downloadingSlice))

	for _, repo := range downloadingSlice {
		err := startDownload(repo)
		if err != nil {
			logger.Errorln("resume download failed:", err)
		}
	}

	// run supervisor

	supervisor()
}
