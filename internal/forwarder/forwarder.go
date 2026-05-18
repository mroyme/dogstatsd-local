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
	mu      sync.Mutex
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
	f.started = true

	return nil
}

func (f *Forwarder) Forward(data []byte) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.conn == nil {
		return
	}

	_, err := f.conn.Write(data)
	if err != nil {
		f.Logger.Error("forward error", "addr", f.Address, "err", err)
	}
}

func (f *Forwarder) Stop() error {
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
