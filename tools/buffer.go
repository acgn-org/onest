package tools

import (
	"sync"
)

var BufferCopy = sync.Pool{
	New: func() any {
		return make([]byte, 32*1024)
	},
}

type BufferHttpUtil struct{}

func (b BufferHttpUtil) Get() []byte {
	return BufferCopy.Get().([]byte)
}

func (b BufferHttpUtil) Put(data []byte) {
	BufferCopy.Put(data)
}
