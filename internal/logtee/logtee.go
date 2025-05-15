package logtee

import (
	"bufio"
	"bytes"
	"container/list"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/tools"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
)

var reader io.Reader
var scan = make(chan struct{}, 1)

var bufLock sync.Mutex
var buf bytes.Buffer

var ringLock sync.Mutex
var ringBuffer = NewRingBuffer(config.Server.Get().LogRingSize)

var subLock sync.Mutex
var subscribers = list.New() // => chan [][]byte

func init() {
	pipeReader, pipeWriter := io.Pipe()
	log.SetOutput(io.MultiWriter(os.Stdout, pipeWriter))
	reader = pipeReader

	go _CopyWorker()
	go _SubscribeWorker()
}

func _CopyWorker() {
	for {
		buffer := tools.BufferCopy.Get().([]byte)
		n, err := reader.Read(buffer)
		if err != nil {
			logfield.New(logfield.ComLogTee).Fatalln("read log failed:", err)
		}
		bufLock.Lock()
		buf.Write(buffer[:n])
		bufLock.Unlock()
		tools.BufferCopy.Put(buffer)

		select {
		case scan <- struct{}{}:
		default:
		}
	}
}

func _SubscribeWorker() {
	for {
		<-scan

		var lines = list.New() // => []byte

		bufLock.Lock()

		ringLock.Lock()
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			line := scanner.Bytes()
			ringBuffer.Add(line)
			lines.PushBack(line)
		}
		ringLock.Unlock()

		remaining := buf.Bytes()
		buf.Reset()
		buf.Write(remaining)

		bufLock.Unlock()

		linesBytes := make([][]byte, lines.Len())
		for i, el := 0, lines.Front(); el != nil; i, el = i+1, el.Next() {
			linesBytes[i] = el.Value.([]byte)
		}

		subLock.Lock()
		for el := subscribers.Front(); el != nil; el = el.Next() {
			el.Value.(chan [][]byte) <- linesBytes
		}
		subLock.Unlock()
	}
}
