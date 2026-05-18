package forwarder

import (
	"errors"
	"net"
	"sync"

	"github.com/charmbracelet/log"
)

type Forwarder struct {
	Logger  *log.Logger
	Address string
	conn    *net.UDPConn
	started bool
	ch      chan []byte
	mu      sync.Mutex
	wg      sync.WaitGroup
}

func (f *Forwarder) Start() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.started {
		return errors.New("forwarder already started")
	}

	addr, err := net.ResolveUDPAddr("udp", f.Address)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	f.conn = conn
	f.ch = make(chan []byte, 10000)
	f.started = true

	ch := f.ch
	logger := f.Logger

	f.wg.Go(func() {
		for data := range ch {
			_, err := conn.Write(data)
			if err != nil && logger != nil {
				logger.Error("forward error", "addr", f.Address, "err", err)
			}
		}
	})

	return nil
}

func (f *Forwarder) Forward(data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.ch == nil {
		return
	}

	select {
	case f.ch <- data:
	default:
		if f.Logger != nil {
			f.Logger.Error("forward buffer full, dropping datagram", "addr", f.Address)
		}
	}
}

func (f *Forwarder) Stop() error {
	f.mu.Lock()
	if f.ch != nil {
		close(f.ch)
		f.ch = nil
	}
	f.mu.Unlock()

	f.wg.Wait()

	f.mu.Lock()
	defer f.mu.Unlock()

	if f.conn != nil {
		err := f.conn.Close()
		f.conn = nil
		f.started = false
		return err
	}
	return nil
}
