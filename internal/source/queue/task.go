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

	lock     sync.RWMutex
	priority int32
	state    *client.File
	errorAt  time.Time
}

func retrieveDownloadInfo(repo repository.Download) (*DownloadTask, error) {
	msg, err := source.Telegram.GetMessage(repo.Item.ChannelID, repo.MsgID)
	if err != nil {
		return nil, err
	}
	video, ok := source.Telegram.GetMessageVideo(msg)
	if !ok {
		return nil, fmt.Errorf("download %d is not a video message", repo.ID)
	}
	file, err := source.Telegram.DownloadFile(video.Video.Id, repo.Priority)
	if err != nil {
		return nil, err
	}
	return &DownloadTask{
		RepoID:    repo.ID,
		MsgID:     repo.MsgID,
		VideoFile: video,
		priority:  repo.Priority,
		state:     file,
	}, nil
}
