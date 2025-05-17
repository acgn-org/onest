package queue

import (
	"errors"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/acgn-org/onest/tools"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

func NewTask(model repository.Download) (*DownloadTask, error) {
	task := &DownloadTask{
		ID:        model.ID,
		ChannelID: model.Item.ChannelID,
		MsgID:     model.MsgID,
		logger: logfield.New(logfield.ComTask).
			WithField("id", model.ID),
		priority: model.Priority,
	}
	return task, task.doUpdateOrDownload()
}

type DownloadTask struct {
	ID uint

	ChannelID int64
	MsgID     int64

	logger  log.FieldLogger
	isFatal atomic.Bool

	lock     sync.RWMutex
	priority int32
	// maybe nil
	state          *client.File
	stateUpdatedAt time.Time
	errorAt        time.Time
	errorCount     uint8
}

func (task *DownloadTask) fatal() {
	task.logger.Errorln("task failed with too many errors or an fatal error")
	task.isFatal.Store(true)
}

func (task *DownloadTask) writeFatalStateToDatabase() error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	err := downloadRepo.UpdateDownloadFatal(task.ID)
	if err != nil {
		return err
	}
	return downloadRepo.Commit().Error
}

func (task *DownloadTask) completeDownload() error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	download, err := downloadRepo.FirstByID(task.ID)
	if err != nil {
		return err
	}

	err = downloadRepo.UpdateDownloadComplete(task.ID)
	if err != nil {
		return err
	}

	itemRepo := repository.ItemRepository{Repository: downloadRepo.Repository}
	item, err := itemRepo.FirstItemByID(download.ItemID)
	if err != nil {
		return err
	}

	targetPath := item.TargetPath
	targetName, err := tools.ConvertWithPattern(download.Text, item.Regexp, item.Pattern)
	if err != nil {
		task.logger.Errorln("convert target path failed:", err)
		return err
	}

	err = os.Rename(task.state.Local.Path, path.Join(targetPath, targetName))
	if err != nil {
		task.logger.Errorln("move file failed:", err)
		return err
	}

	return downloadRepo.Commit().Error
}

func (task *DownloadTask) getVideoFile() (bool, error) {
	msg, err := source.Telegram.GetMessage(task.ChannelID, task.MsgID)
	if err != nil {
		_ = task.setError(err.Error(), false)
		return false, err
	}
	messageVideo, ok := source.Telegram.GetMessageVideo(msg)
	if !ok {
		return false, nil
	}
	task.state = messageVideo.Video.Video
	task.stateUpdatedAt = time.Now()
	return true, nil
}

func (task *DownloadTask) setError(msg string, fatalNow bool) error {
	logger := task.logger.WithField("msg", msg)

	task.errorAt = time.Now()
	task.errorCount++

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	if err := downloadRepo.UpdateDownloadError(task.ID, msg, task.errorAt.Unix()); err != nil {
		logger.Errorln("save error message to database failed:", err)
		return err
	}
	if err := downloadRepo.Commit().Error; err != nil {
		return err
	}

	if fatalNow || config.Telegram.Get().MaxDownloadError <= task.errorCount {
		task.fatal()
	}
	return nil
}

func (task *DownloadTask) doUpdateOrDownload() error {
	if task.state == nil {
		ok, err := task.getVideoFile()
		if err != nil {
			return err
		} else if !ok {
			msg := "no video file found"
			_ = task.setError(msg, true)
			return errors.New(msg)
		}
	}

	newState, err := source.Telegram.DownloadFile(task.state.Id, task.priority)
	if err != nil {
		_ = task.setError(err.Error(), false)
	} else {
		task.state = newState
		task.stateUpdatedAt = time.Now()
	}
	return nil
}

func (task *DownloadTask) UpdateOrDownload() error {
	if task.isFatal.Load() {
		return nil
	}

	task.lock.Lock()
	defer task.lock.Unlock()

	return task.doUpdateOrDownload()
}
