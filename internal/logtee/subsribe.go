package logtee

import "container/list"

type Subscribe struct {
	element *list.Element
	receive chan [][]byte
	write   chan *list.List
	read    chan [][]byte
}

func NewSubscribe() ([][]byte, *Subscribe) {
	receive := make(chan [][]byte)

	subLock.Lock()
	ringLock.Lock()
	ringLogs := ringBuffer.GetAll()
	el := subscribers.PushFront(receive)
	ringLock.Unlock()
	subLock.Unlock()

	sub := Subscribe{
		element: el,
		receive: receive,
		write:   make(chan *list.List),
		read:    make(chan [][]byte),
	}
	go sub.receiveWorker()
	return ringLogs, &sub
}

func (sub *Subscribe) Listen() <-chan [][]byte {
	return sub.read
}

func (sub *Subscribe) Close() error {
	subLock.Lock()
	subscribers.Remove(sub.element)
	subLock.Unlock()
	close(sub.receive)
	return nil
}

func (sub *Subscribe) receiveWorker() {
	buf := list.New() // => [][]byte
	for {
		data, ok := <-sub.receive
		if !ok {
			close(sub.write)
			return
		}
		buf.PushBack(data)

		select {
		case sub.write <- buf:
			buf = list.New()
		default:
		}
	}
}

func (sub *Subscribe) sendWorker() {
	for {
		buf, ok := <-sub.write
		if !ok {
			close(sub.read)
			return
		}

		var length int
		for el := buf.Front(); el != nil; el = el.Next() {
			length += len(el.Value.([][]byte))
		}
		var data = make([][]byte, 0, length)
		for el := buf.Front(); el != nil; el = el.Next() {
			data = append(data, el.Value.([][]byte)...)
		}

		sub.read <- data
	}
}
