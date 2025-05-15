package logtee

import "sync"

type RingBuffer struct {
	data  [][]byte
	size  int
	start int
	count int
	lock  sync.RWMutex
}

func NewRingBuffer(capacity int) *RingBuffer {
	return &RingBuffer{
		data: make([][]byte, capacity),
		size: capacity,
	}
}

func (rb *RingBuffer) Add(line []byte) {
	rb.lock.Lock()
	defer rb.lock.Unlock()

	index := (rb.start + rb.count) % rb.size
	if rb.count < rb.size {
		rb.data[index] = line
		rb.count++
	} else {
		// buffer full, overwrite oldest
		rb.data[rb.start] = line
		rb.start = (rb.start + 1) % rb.size
	}
}

func (rb *RingBuffer) GetAll() [][]byte {
	rb.lock.RLock()
	defer rb.lock.RUnlock()

	result := make([][]byte, rb.count)
	for i := 0; i < rb.count; i++ {
		result[i] = rb.data[(rb.start+i)%rb.size]
	}
	return result
}
