package queue

import "time"

func supervisor() {
	for {
		time.Sleep(time.Second * 10)

		// todo proactively update stats
	}
}
