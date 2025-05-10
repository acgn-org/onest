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
	parallel = config.Telegram.MaxParallelDownload
	lock     sync.RWMutex

	// active or pending is determined by tdlib
	queued  int64
	active  int64
	pending int64
)

func init() {
	logger := logfield.New(logfield.ComQueue).WithAction("init")

	downloadRepo := repository.DownloadRepository{
		DB: database.DB,
	}

	var err error
	queued, err = downloadRepo.CountNotDownloaded()
	if err != nil {
		logger.Fatalln("count download failed:", err)
	}
}

func Clean() error {
	lock.Lock()
	defer lock.Unlock()

	if err := source.Telegram.RemoveDownloads(); err != nil {
		return err
	}
	if err := source.Telegram.CleanDownloadDirectory(); err != nil {
		return err
	}
	return nil
}
