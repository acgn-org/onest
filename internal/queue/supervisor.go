package queue

import (
	"errors"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/repository"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func supervisor() {
	instance := _Supervisor{
		logger: logfield.New(logfield.ComQueueSupervisor),
	}
	go instance.WorkerTaskControl()
}

type _Supervisor struct {
	logger log.FieldLogger
}

func (s _Supervisor) WorkerTaskControl() {
	for {
		time.Sleep(time.Second * 10)
		s.TaskControl()
	}
}

func (s _Supervisor) TaskControl() {
	lock.Lock()
	defer lock.Unlock()

	for key, task := range downloading {
		// remove tasks with fatal state
		if task.isFatal.Load() {
			err := task.WriteFatalStateToDatabase()
			if err != nil {
				s.logger.WithField("id", task.RepoID).Errorln("failed to write download task fatal state into database:", err)
			} else {
				delete(downloading, key)
			}
			continue
		}

		// restart downloads with error
		if !task.errorAt.IsZero() && time.Since(task.errorAt) > time.Second*10 {
			task.errorAt = time.Time{}
			_ = task.UpdateOrDownload()
			continue
		}

		// proactively update stats
		if time.Since(task.stateUpdatedAt) > time.Second*15 {
			_ = task.UpdateOrDownload()
		}
	}

	// clean up downloads
	if len(downloading) == 0 {
		err := clean()
		if err != nil {
			s.logger.Errorln("failed to clean up resources:", err)
		}
	}

	// maintain parallel downloads
	numToDownload := int(config.Telegram.Get().MaxParallelDownload) - len(downloading)
	if numToDownload > 0 {
		downloadRepo := database.NewRepository[repository.DownloadRepository]()

		for range numToDownload { // todo customizable download order
			model, err := downloadRepo.FirstToDownload()
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					// no more
					break
				}
				s.logger.Errorln("load download task from database failed:", err)
				continue
			}

			if err := startDownload(*model); err != nil {
				s.logger.Errorln("error occurred while start download task:", err)
			}
		}
	}
}
