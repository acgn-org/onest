package queue

import (
	"errors"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/zelenin/go-tdlib/client"
	"sync"
	"sync/atomic"
	"time"
)

func NewTask(model repository.Download) (*DownloadTask, error) {
	task := &DownloadTask{
		RepoID:    model.ID,
		ChannelID: model.Item.ChannelID,
		MsgID:     model.MsgID,
		priority:  model.Priority,
	}
	return task, task.UpdateOrDownload()
}

type DownloadTask struct {
	RepoID    uint
	ChannelID int64
	MsgID     int64

	lock     sync.RWMutex
	priority int32
	// maybe nil
	state          *client.File
	stateUpdatedAt time.Time
	errorAt        time.Time
	errorCount     uint8
	isFatal        atomic.Bool
}

func (task *DownloadTask) fatal() {
	task.isFatal.Store(true)
}

func (task *DownloadTask) getVideoFile() (bool, error) {
	msg, err := source.Telegram.GetMessage(task.ChannelID, task.MsgID)
	if err != nil {
		_ = task.setError(err.Error())
		return false, err
	}
	video, ok := source.Telegram.GetMessageVideo(msg)
	if !ok {
		return false, nil
	}
	task.state = video.Video
	return true, nil
}

func (task *DownloadTask) setError(msg string) error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	task.errorAt = time.Now()
	task.errorCount++

	if err := downloadRepo.UpdateDownloadError(task.RepoID, msg, task.errorAt.Unix()); err != nil {
		// todo warn with log
		return err
	}
	if err := downloadRepo.Commit().Error; err != nil {
		return err
	}

	if config.Telegram.Get().MaxDownloadError <= task.errorCount {
		task.fatal()
	}
	return nil
}

func (task *DownloadTask) UpdateOrDownload() error {
	if task.isFatal.Load() {
		return nil
	}

	task.lock.Lock()
	defer task.lock.Unlock()

	if task.state == nil {
		ok, err := task.getVideoFile()
		if err != nil {
			return err
		} else if !ok {
			task.fatal()
			return errors.New("no video file found")
		}
	}

	newState, err := source.Telegram.DownloadFile(task.state.Id, task.priority)
	if err != nil {
		_ = task.setError(err.Error())
	} else {
		task.state = newState
		task.stateUpdatedAt = time.Now()
	}
	return nil
}
