package queue

import (
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"sync/atomic"
	"time"
)

var ActivateTaskControl = make(chan struct{})

func supervisor() {
	instance := _Supervisor{
		logger:  logfield.New(logfield.ComQueueSupervisor),
		Cleaned: &atomic.Bool{},
	}
	go instance.WorkerTaskControl()
	go instance.WorkerListen()
}

type _Supervisor struct {
	logger  log.FieldLogger
	Cleaned *atomic.Bool
}

func (s _Supervisor) WorkerTaskControl() {
	var slowDown bool
	for {
		sleep := time.Second * 10
		if slowDown {
			sleep = time.Minute * 5
		}

		select {
		case <-time.After(sleep):
		case <-ActivateTaskControl:
		}

		slowDown = s.TaskControl()
	}
}

func (s _Supervisor) TaskControl() (slowDown bool) {
	lock.Lock()
	defer lock.Unlock()

	for key, task := range downloading {
		logger := s.logger.WithField("task", key)

		// remove tasks with fatal state
		if task.isFatal.Load() {
			err := task.writeFatalStateToDatabase()
			if err != nil {
				logger.Errorln("failed to write download task fatal state into database:", err)
			} else {
				delete(downloading, key)
			}
			continue
		}

		task.lock.Lock()

		if !task.errorAt.IsZero() && time.Since(task.errorAt) > time.Second*10 {
			// restart downloads with error
			task.errorAt = time.Time{}
			if err := task.doUpdateOrDownload(); err != nil {
				logger.Errorln("failed to restart task:", err)
			}
		} else if time.Since(task.stateUpdatedAt) > time.Second*15 {
			// proactively update stats
			if err := task.doUpdateOrDownload(); err != nil {
				logger.Errorln("failed to update task state:", err)
			}
		}

		// proceed downloads completed
		if task.state != nil && task.state.Local.IsDownloadingCompleted {
			if err := task.completeDownload(); err != nil {
				logger.Errorln("failed to complete download task:", err)
			} else {
				delete(downloading, key)
			}
		}

		task.lock.Unlock()
	}

	// maintain number of parallel downloads
	numToDownload := int(config.Telegram.Get().MaxParallelDownload) - len(downloading)
	if numToDownload > 0 {
		downloadRepo := database.NewRepository[repository.DownloadRepository]()

		models, err := downloadRepo.EarliestToDownload(numToDownload)
		if err != nil {
			s.logger.Errorln("load download task from database failed:", err)
		} else if len(models) != 0 {
			s.Cleaned.Store(false)
			for _, model := range models {
				if err := startDownload(model); err != nil {
					s.logger.Errorln("error occurred while start download task:", err)
				}
			}
		} else if len(downloading) == 0 {
			// clean up downloads
			if !s.Cleaned.Load() {
				err := clean()
				if err != nil {
					s.logger.Errorln("failed to clean up resources:", err)
				} else {
					s.Cleaned.Store(true)
				}
			}
			return true
		}
	}
	return false
}

func (s _Supervisor) WorkerListen() {
	listener := source.Telegram.GetListener()
	defer listener.Close()

	for {
		update := <-listener.Updates

		switch update.GetType() {

		case client.TypeUpdateFile:
			file := update.(*client.UpdateFile).File
			lock.Lock()
			for _, task := range downloading {
				task.lock.Lock()
				if task.state != nil && task.state.Id == file.Id {
					task.state = file
					if file.Local.IsDownloadingCompleted {
						err := task.completeDownload()
						if err != nil {
							s.logger.Errorln("failed to complete download task:", err)
						}
					}
				}
				task.lock.Unlock()
			}
			lock.Unlock()

		case client.TypeUpdateNewMessage:
			// match new downloads
			err := ScanAndCreateNewDownloadTasks()
			if err != nil {
				s.logger.Errorln("failed to create downloads with message:", err)
				continue
			}

		default:
			goto skip
		}

		select {
		case ActivateTaskControl <- struct{}{}:
		default:
		}

	skip:
	}
}
