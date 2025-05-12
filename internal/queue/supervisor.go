package queue

import (
	"github.com/acgn-org/onest/internal/logfield"
	log "github.com/sirupsen/logrus"
	"time"
)

func supervisor() {
	// todo proactively update stats
	// todo clean downloads

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
		}
	}
}
