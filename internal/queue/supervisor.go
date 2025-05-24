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
			s.logger.Debugln("activate task control by signal")
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
		if task.log.isFatal.Load() {
			err := task._WriteFatalStateToDatabase()
			if err != nil {
				logger.Errorln("failed to write download task fatal state into database:", err)
			} else {
				delete(downloading, key)
			}
			continue
		}

		if state := task.state.Load(); state == nil || time.Since(state.UpdatedAt) > time.Second*10 {
			// proactively update stats, or restart downloads with error
			if err := task.UpdateOrDownload(false); err != nil {
				logger.Errorln("failed to update task state:", err)
			}
		}

		// proceed downloads completed
		if state := task.state.Load(); state != nil && state.File.Local.IsDownloadingCompleted {
			if err := task.CompleteDownload(); err != nil {
				logger.Errorln("failed to complete download task:", err)
			} else {
				delete(downloading, key)
			}
		}
	}

	// maintain number of parallel downloads
	numToDownload := int(config.Telegram.Get().MaxParallelDownload) - len(downloading)
	if numToDownload > 0 {
		downloadRepo := database.NewRepository[repository.DownloadRepository]()
		models, err := downloadRepo.GetForDownloadPreloadItem(numToDownload)
		if err != nil {
			s.logger.Errorln("load download task from database failed:", err)
		} else if len(models) != 0 {
			s.Cleaned.Store(false)
			for _, model := range models {
				if err := startDownload(model.Item.ChannelID, model); err != nil {
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
			var isFileCompleted bool
			file := update.(*client.UpdateFile).File
			lock.Lock()
			for id, task := range downloading {
				if state := task.state.Load(); state != nil && state.File.Id == file.Id {
					task.state.Store(&TaskFileState{
						File:      file,
						UpdatedAt: time.Now(),
					})
					if file.Local.IsDownloadingCompleted {
						isFileCompleted = true
						err := task.CompleteDownload()
						if err != nil {
							s.logger.Errorln("failed to complete download task:", err)
						} else {
							delete(downloading, id)
						}
					}
				}
			}
			lock.Unlock()
			if isFileCompleted {
				select {
				case ActivateTaskControl <- struct{}{}:
				default:
				}
			}

		case client.TypeUpdateNewMessage:
			// match new downloads
			created, err := ScanAndCreateNewDownloadTasks()
			if err != nil {
				s.logger.Errorln("failed to create downloads with message:", err)
				continue
			} else if created > 0 {
				select {
				case ActivateTaskControl <- struct{}{}:
				default:
				}
			}

		}
	}
}
