package queue

import (
	"context"
	"errors"
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/acgn-org/onest/tools"
	log "github.com/sirupsen/logrus"
	"github.com/zelenin/go-tdlib/client"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"sync"
	"sync/atomic"
	"time"
)

func UpdatePriority(id uint, priority int32) {
	download, ok := queue.Load(id)
	if !ok {
		return
	}
	download.priority.Store(priority)
	err := download.UpdateOrDownload(true)
	if err != nil {
		logfield.New(logfield.ComQueue).WithAction("update").Warnf("update priority of task %d to %d failed: %v", id, priority, err)
	}
}

func ForceAddDownloadQueue(channelId int64, model repository.Download) error {
	task, ok := queue.Load(model.ID)
	if ok {
		return task.UpdateOrDownload(false)
	}

	return startDownload(channelId, model)
}

func RemoveTasks(ids ...uint) {
	for _, id := range ids {
		task, ok := queue.LoadAndDelete(id)
		if !ok {
			continue
		}

		if err := task.Terminate(); err != nil {
			logfield.New(logfield.ComQueue).WithAction("remove").Errorf("terminate task %d with error: %v", id, err)
		}
	}
}

type TaskErrorState struct {
	Err string
	At  time.Time
}

type TaskFileState struct {
	File      *client.File
	UpdatedAt time.Time
}

func NewTaskLogger(id uint, logger log.FieldLogger) TaskLogger {
	taskLogger := TaskLogger{
		logger:     logger,
		id:         id,
		isFatal:    &atomic.Bool{},
		error:      &atomic.Pointer[TaskErrorState]{},
		errorCount: &atomic.Uint32{},
	}
	taskLogger.error.Store(&TaskErrorState{})
	return taskLogger
}

type TaskLogger struct {
	logger log.FieldLogger
	id     uint

	isFatal    *atomic.Bool
	error      *atomic.Pointer[TaskErrorState]
	errorCount *atomic.Uint32
}

func (tl TaskLogger) _SaveErrorState(state TaskErrorState) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	downloadRepo := database.BeginRepositoryWithContext[repository.DownloadRepository](ctx)
	defer downloadRepo.Rollback()

	if err := downloadRepo.UpdateDownloadError(tl.id, state.Err, state.At.Unix()); err != nil {
		return err
	}
	return downloadRepo.Commit().Error
}

func (tl TaskLogger) WithField(key string, value interface{}) TaskLogger {
	tl.logger = tl.logger.WithField(key, value)
	return tl
}

func (tl TaskLogger) Errorln(args ...interface{}) {
	errorState := TaskErrorState{
		Err: fmt.Sprintln(args...),
		At:  time.Now(),
	}
	tl.error.Store(&errorState)
	newErrorCount := tl.errorCount.Add(uint32(1))

	if newErrorCount >= config.Telegram.Get().MaxDownloadError {
		tl.FatalNow()
	}

	go func() {
		if err := tl._SaveErrorState(errorState); err != nil {
			tl.logger.Warnln("save error message to database failed:", err)
		}
	}()

	tl.logger.Errorln(args...)
}

func (tl TaskLogger) FatalNow() {
	tl.logger.Debugln("fatal now")
	tl.isFatal.Store(true)
}

func NewTask(channelId int64, model repository.Download) (*DownloadTask, error) {
	task := &DownloadTask{
		ID:        model.ID,
		ChannelID: channelId,
		MsgID:     model.MsgID,
		log:       NewTaskLogger(model.ID, logfield.New(logfield.ComTask).WithField("id", model.ID)),
	}
	task.priority.Store(model.Priority)
	return task, task.UpdateOrDownload(false)
}

type DownloadTask struct {
	ID uint

	ChannelID int64
	MsgID     int64

	log TaskLogger

	completed    atomic.Bool
	lockComplete sync.Mutex

	priority atomic.Int32
	state    atomic.Pointer[TaskFileState] // maybe nil
}

func (task *DownloadTask) _WriteFatalStateToDatabase() error {
	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	errorState := task.log.error.Load()
	err := downloadRepo.UpdateDownloadFatal(task.ID, errorState.Err, errorState.At.Unix())
	if err != nil {
		return err
	}
	return downloadRepo.Commit().Error
}

func (task *DownloadTask) CompleteDownload() (ok bool, err error) {
	if task.completed.Load() {
		return true, nil
	}
	if !task.lockComplete.TryLock() {
		return false, nil
	}
	ok = true
	defer task.lockComplete.Unlock()

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	download, err := downloadRepo.FirstByIDForUpdate(task.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			task.log.FatalNow()
		}
		task.log.Errorln("lookup task from database failed:", err)
		return
	}

	err = downloadRepo.UpdateDownloadComplete(task.ID)
	if err != nil {
		task.log.Errorln("mark download complete failed:", err)
		return
	}

	itemRepo := repository.ItemRepository{Repository: downloadRepo.Repository}
	item, err := itemRepo.FirstItemByID(download.ItemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			task.log.FatalNow()
		}
		task.log.Errorln("lookup item from database failed:", err)
		return
	}

	targetPath := item.TargetPath
	targetName, err := tools.ConvertPatternRegexpString(download.Text, item.Regexp, item.Pattern)
	if err != nil {
		task.log.Errorln("convert target path failed:", err)
		return
	}

	state := task.state.Load()
	if state == nil {
		panic("complete download called without file state")
	}

	info, err := os.Stat(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(targetPath, config.FilePerm)
			if err != nil {
				task.log.Errorln("create target directory failed:", err)
				return
			}
		} else {
			task.log.Errorln("stat target directory failed:", err)
			return
		}
	} else if !info.IsDir() {
		msg := fmt.Sprintf("target path '%s' is not a directory", targetPath)
		task.log.Errorln(msg)
		return ok, errors.New(msg)
	}

	fullPath := path.Join(targetPath, targetName) + path.Ext(state.File.Local.Path)
	err = os.Rename(state.File.Local.Path, fullPath)
	if err != nil {
		task.log.logger.Debugln("rename file failed:", err)

		var fileSource *os.File
		fileSource, err = os.OpenFile(state.File.Local.Path, os.O_RDONLY, 0600)
		if err != nil {
			task.log.Errorln("open source file failed:", err)
			return
		}
		defer fileSource.Close()

		var fileTarget *os.File
		fileTarget, err = os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, config.FilePerm)
		if err != nil {
			task.log.Errorln("create file failed:", err)
			return
		}
		defer fileTarget.Close()

		buffer := tools.BufferCopy.Get().([]byte)
		defer tools.BufferCopy.Put(buffer)
		_, err = io.CopyBuffer(fileTarget, fileSource, buffer)
		if err != nil {
			task.log.Errorln("copy file failed:", err)
			return
		}

		err = os.Remove(state.File.Local.Path)
		if err != nil {
			task.log.Errorln("remove source file failed:", err)
		}
	}

	if err = downloadRepo.Commit().Error; err != nil {
		task.log.Errorln("save changes into database failed:", err)
		return
	}

	task.completed.Store(true)
	return
}

func (task *DownloadTask) GetVideoFile() (bool, error) {
	msg, err := source.Telegram.GetMessage(task.ChannelID, task.MsgID)
	if err != nil {
		task.log.Errorln("get message failed:", err)
		return false, err
	}
	messageVideo, ok := source.Telegram.GetMessageVideo(msg)
	if !ok {
		return false, nil
	}
	task.state.CompareAndSwap(nil, &TaskFileState{
		File:      messageVideo.Video.Video,
		UpdatedAt: time.Now(),
	})
	return true, nil
}

func (task *DownloadTask) UpdateOrDownload(forceStart bool) error {
	if task.log.isFatal.Load() {
		return nil
	}

	state := task.state.Load()

	if state == nil {
		ok, err := task.GetVideoFile()
		if err != nil {
			return err
		} else if !ok {
			msg := "no video file found"
			task.log.Errorln(msg)
			task.log.FatalNow()
			return errors.New(msg)
		}
		state = task.state.Load()
	} else if !forceStart {
		file, err := source.Telegram.GetFile(state.File.Id)
		if err != nil {
			task.log.Errorln("get download file state failed:", err)
			return err
		}
		task.state.Store(&TaskFileState{
			File:      file,
			UpdatedAt: time.Now(),
		})
		if file.Local.IsDownloadingCompleted || file.Local.IsDownloadingActive {
			return nil
		}
	}

	file, err := source.Telegram.DownloadFile(state.File.Id, task.priority.Load(), false)
	if err != nil {
		task.log.Errorln("request download failed:", err)
	} else {
		task.state.Store(&TaskFileState{
			File:      file,
			UpdatedAt: time.Now(),
		})
	}
	return err
}

func (task *DownloadTask) Terminate() error {
	if task.log.isFatal.Load() {
		return nil
	}

	if state := task.state.Load(); state != nil {
		if !state.File.Local.IsDownloadingCompleted {
			if err := source.Telegram.CancelDownloadFile(state.File.Id); err != nil {
				return err
			}
		}
		return source.Telegram.RemoveFileFromDownloads(state.File.Id)
	}
	return nil
}
