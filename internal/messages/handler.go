package messages

import (
	"errors"
	"github.com/charmbracelet/log"
	"sync"
)

type OutputHandler func([]byte) error

type Handler struct {
	Logger     *log.Logger
	Out        OutputHandler
	PoolSize   int
	BufferSize int
	msgCh      chan []byte
	wg         sync.WaitGroup
}

func (h *Handler) New() *Handler {
	h.msgCh = make(chan []byte, h.BufferSize)
	h.wg.Add(h.PoolSize)

	// build a pool of goroutines to listen for messages to process
	for i := 0; i < h.PoolSize; i++ {
		go func() {
			defer h.wg.Done()
			for msg := range h.msgCh {
				err := h.Out(msg)
				if err != nil {
					return
				}
			}
		}()
	}

	return h
}

// Handle submits to the pool
func (h *Handler) Handle(msg []byte) error {
	select {
	case h.msgCh <- msg:
		return nil
	default:
	}

	return errors.New("pool capacity exceeded")
}

func (h *Handler) Stop() {
	close(h.msgCh)
	h.wg.Wait()
}
