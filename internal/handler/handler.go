package handler

import (
	"errors"
	"sync"
)

type MessageHandler func([]byte) error

type AsyncMessageHandler interface {
	Handler([]byte) error
	Stop()
}

type handler struct {
	fn    MessageHandler
	msgCh chan []byte
	wg    sync.WaitGroup
}

func NewAsyncMessageHandler(fn MessageHandler, poolSize int, bufferSize int) AsyncMessageHandler {
	h := &handler{
		fn:    fn,
		msgCh: make(chan []byte, bufferSize),
	}

	h.wg.Add(poolSize)

	// build a pool of goroutines to listen for messages to process
	for i := 0; i < poolSize; i++ {
		go func() {
			defer h.wg.Done()
			for msg := range h.msgCh {
				err := h.fn(msg)
				if err != nil {
					return
				}
			}
		}()
	}

	return h
}

// Handler submits to the pool
func (a *handler) Handler(msg []byte) error {
	select {
	case a.msgCh <- msg:
		return nil
	default:
	}

	return errors.New("POOL_CAPACITY_EXCEEDED")
}

func (a *handler) Stop() {
	close(a.msgCh)
	a.wg.Wait()
}
