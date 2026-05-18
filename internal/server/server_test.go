package server

import (
	"fmt"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

func TestServerReceivesMessages(t *testing.T) {
	var received atomic.Int32

	handler := func(msg []byte) error {
		received.Add(1)
		return nil
	}

	srv := NewServer("127.0.0.1:0", handler, log.Default())

	done := make(chan error, 1)
	go func() {
		done <- srv.Listen()
	}()

	// Give the server a moment to start listening
	time.Sleep(100 * time.Millisecond)

	// We need to get the actual allocated port — resolve it from the server
	// Since our server doesn't expose the addr, we'll use a fixed port for testing
	// Stop the server and restart with a known port
	srv.Stop()

	// Use a known available port
	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv = NewServer(addr, handler, log.Default())
	go func() {
		done <- srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("page.views:1|c"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	// Give the server time to process
	time.Sleep(200 * time.Millisecond)

	srv.Stop()

	if got := received.Load(); got != 1 {
		t.Errorf("received %d messages, want 1", got)
	}
}

func TestServerSplitsMultiMessageDatagrams(t *testing.T) {
	var received atomic.Int32
	var messages [][]byte

	handler := func(msg []byte) error {
		received.Add(1)
		messages = append(messages, msg)
		return nil
	}

	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv := NewServer(addr, handler, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Send multiple messages in a single datagram separated by newlines
	_, err = conn.Write([]byte("page.views:1|c\nfuel.level:0.5|g"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	srv.Stop()

	if got := received.Load(); got != 2 {
		t.Errorf("received %d messages, want 2", got)
	}

	if len(messages) >= 2 {
		if string(messages[0]) != "page.views:1|c" {
			t.Errorf("messages[0] = %q, want %q", messages[0], "page.views:1|c")
		}
		if string(messages[1]) != "fuel.level:0.5|g" {
			t.Errorf("messages[1] = %q, want %q", messages[1], "fuel.level:0.5|g")
		}
	}
}

func TestServerStripsCarriageReturns(t *testing.T) {
	var received atomic.Int32
	var lastMsg []byte

	handler := func(msg []byte) error {
		received.Add(1)
		lastMsg = msg
		return nil
	}

	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv := NewServer(addr, handler, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	// Send message with trailing \r
	_, err = conn.Write([]byte("page.views:1|c\r\n"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	srv.Stop()

	if got := received.Load(); got != 1 {
		t.Errorf("received %d messages, want 1", got)
	}

	if string(lastMsg) != "page.views:1|c" {
		t.Errorf("message = %q, want %q", lastMsg, "page.views:1|c")
	}
}

func TestServerStops(t *testing.T) {
	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv := NewServer(addr, func(msg []byte) error { return nil }, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	err := srv.Stop()
	if err != nil {
		t.Errorf("Stop() returned error: %v", err)
	}
}

func findAvailablePort(t *testing.T) string {
	t.Helper()
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()
	return fmt.Sprintf("%d", conn.LocalAddr().(*net.UDPAddr).Port)
}
