package server

import (
	"errors"
	"github.com/charmbracelet/log"
	"github.com/mroyme/dogstatsd-local/internal/messages"
	"net"
	"sync"
	"time"
)

type Server interface {
	Listen() error
	Stop() error
}

func NewServer(addr string, fn messages.OutputHandler, logger *log.Logger) Server {
	return &udpServer{
		logger:        logger,
		msgHandler:    fn,
		rawAddr:       addr,
		readDeadline:  time.Second / 4,
		writeDeadline: time.Second / 4,
		errCh:         make(chan error, 1),
		stopCh:        make(chan struct{}),
		wg:            sync.WaitGroup{},
	}
}

type udpServer struct {
	logger        *log.Logger
	msgHandler    messages.OutputHandler
	rawAddr       string
	readDeadline  time.Duration
	writeDeadline time.Duration
	stopCh        chan struct{}
	errCh         chan error
	wg            sync.WaitGroup
}

func (u *udpServer) Listen() error {
	addr, err := net.ResolveUDPAddr("udp", u.rawAddr)
	if err != nil {
		return err
	}

	serverConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	u.wg.Add(1)
	go u.errHandler()

	buf := make([]byte, 1024)
	var respMsg []byte

	for {
		select {
		case <-u.stopCh:
			goto stop
		default:
		}

		err := serverConn.SetDeadline(time.Now().Add(u.readDeadline))
		if err != nil {
			return err
		}

		n, clientAddr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			// check if the error is a timeout error
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				continue
			}

			u.errCh <- err
			continue
		}

		// copy the message and pass it to the format function
		msg := make([]byte, n)
		copy(msg, buf[:n])
		err = u.msgHandler(msg)
		if err != nil {
			return err
		}

		// respond to the origin connection
		err = serverConn.SetWriteDeadline(time.Now().Add(u.writeDeadline))
		if err != nil {
			return err
		}
		_, err = serverConn.WriteToUDP(respMsg, clientAddr)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				continue
			}

			u.errCh <- err
		}
	}

stop:
	serverConn.Close()
	close(u.stopCh)
	return nil
}

func (u *udpServer) errHandler() {
	for err := range u.errCh {
		log.Error(err)
	}

	u.wg.Done()
}

func (u *udpServer) Stop() error {
	// stop the server and wait for it to finish
	u.stopCh <- struct{}{}
	<-u.stopCh

	// close the err channel and wait for any in progress errors to complete
	close(u.errCh)
	u.wg.Wait()
	return nil
}
