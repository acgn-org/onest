package queue

import (
	"fmt"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/zelenin/go-tdlib/client"
	"sync"
	"time"
)

type DownloadTask struct {
	RepoID    uint
	MsgID     int64
	VideoFile *client.Video

	lock           sync.RWMutex
	priority       int32
	state          *client.File
	stateUpdatedAt time.Time
	errorAt        time.Time
}

func retrieveAndStartDownload(model repository.Download) (*DownloadTask, error) {
	msg, err := source.Telegram.GetMessage(model.Item.ChannelID, model.MsgID)
	if err != nil {
		return nil, err
	}
	video, ok := source.Telegram.GetMessageVideo(msg)
	if !ok {
		return nil, fmt.Errorf("download %d is not a video message", model.ID)
	}
	file, err := source.Telegram.DownloadFile(video.Video.Id, model.Priority)
	if err != nil {
		return nil, err
	}
	return &DownloadTask{
		RepoID:         model.ID,
		MsgID:          model.MsgID,
		VideoFile:      video,
		priority:       model.Priority,
		state:          file,
		stateUpdatedAt: time.Now(),
	}, nil
}
