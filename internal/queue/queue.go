package queue

import (
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"sync"
	"sync/atomic"
)

func init() {
	logger := logfield.New(logfield.ComQueue).WithAction("init")

	downloadRepo := database.NewRepository[repository.DownloadRepository]()

	// resume downloads

	downloadingSlice, err := downloadRepo.GetDownloadingPreloadItem()
	if err != nil {
		logger.Fatalln("load downloading failed:", err)
	}

	for _, repo := range downloadingSlice {
		err := startDownload(repo.Item.ChannelID, repo)
		if err != nil {
			logger.Errorln("resume download failed:", err)
		}
	}

	// run supervisor

	supervisor()
}

type _Queue struct {
	addLock sync.Mutex
	// repo Download.ID => *DownloadTask
	queue       sync.Map
	queueLength atomic.Int32
}

var queue = &_Queue{}

func (q *_Queue) AllKeys() []uint {
	maxParallel := int32(config.Telegram.Get().MaxParallelDownload)
	mapLen := q.queueLength.Load()
	if mapLen < maxParallel {
		mapLen = maxParallel
	}

	ids := make([]uint, 0, mapLen)
	q.queue.Range(func(key, value interface{}) bool {
		ids = append(ids, key.(uint))
		return true
	})
	return ids
}

func (q *_Queue) Len() int32 {
	return q.queueLength.Load()
}

func (q *_Queue) Load(key uint) (*DownloadTask, bool) {
	val, ok := q.queue.Load(key)
	if !ok {
		return nil, false
	}
	return val.(*DownloadTask), true
}

func (q *_Queue) Store(key uint, value *DownloadTask) {
	q.addLock.Lock()
	defer q.addLock.Unlock()
	_, loaded := q.queue.Swap(key, value)
	if !loaded {
		q.queueLength.Add(int32(1))
	}
}

func (q *_Queue) Delete(key uint) {
	_, ok := q.queue.LoadAndDelete(key)
	if ok {
		q.queueLength.Add(int32(-1))
	}
}

func (q *_Queue) LoadAndDelete(key uint) (*DownloadTask, bool) {
	val, ok := q.queue.LoadAndDelete(key)
	if ok {
		q.queueLength.Add(int32(-1))
	}
	return val.(*DownloadTask), ok
}

func (q *_Queue) Range(f func(key uint, value *DownloadTask) bool) {
	q.queue.Range(func(key, value interface{}) bool {
		return f(key.(uint), value.(*DownloadTask))
	})
}

// should call after ensured addLock is locked and queue is empty
func (q *_Queue) clean() error {
	q.queue.Clear()
	q.queueLength.Store(int32(0))
	if err := source.Telegram.RemoveAllDownloads(); err != nil {
		return err
	}
	if err := source.Telegram.CleanDownloadDirectory(); err != nil {
		return err
	}
	return nil
}
