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

	srv := NewServer("127.0.0.1:0", handler, nil, log.Default())

	go func() {
		_ = srv.Listen()
	}()

	// Give the server a moment to start listening
	time.Sleep(100 * time.Millisecond)

	// Since the server doesn't expose the addr, restart with a known port
	if err := srv.Stop(); err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}

	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv = NewServer(addr, handler, nil, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	_, err = conn.Write([]byte("page.views:1|c"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	// Give the server time to process
	time.Sleep(200 * time.Millisecond)

	if err := srv.Stop(); err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}

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

	srv := NewServer(addr, handler, nil, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Send multiple messages in a single datagram separated by newlines
	_, err = conn.Write([]byte("page.views:1|c\nfuel.level:0.5|g"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if err := srv.Stop(); err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}

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

	srv := NewServer(addr, handler, nil, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Send message with trailing \r
	_, err = conn.Write([]byte("page.views:1|c\r\n"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if err := srv.Stop(); err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}

	if got := received.Load(); got != 1 {
		t.Errorf("received %d messages, want 1", got)
	}

	if string(lastMsg) != "page.views:1|c" {
		t.Errorf("message = %q, want %q", lastMsg, "page.views:1|c")
	}
}

func TestServerForwardsDatagrams(t *testing.T) {
	var received atomic.Int32
	done := make(chan []byte, 2)

	handler := func(msg []byte) error {
		return nil
	}

	forward := func(datagram []byte) {
		received.Add(1)
		cp := make([]byte, len(datagram))
		copy(cp, datagram)
		done <- cp
	}

	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv := NewServer(addr, handler, forward, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("udp", addr)
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer func() { _ = conn.Close() }()

	// Single message
	_, err = conn.Write([]byte("page.views:1|c"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	// Multi-message datagram with \r\n
	_, err = conn.Write([]byte("fuel.level:0.5|g\r\nsong.length:240|h"))
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	time.Sleep(200 * time.Millisecond)

	if err := srv.Stop(); err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}

	// Collect forwarded datagrams via channel to avoid data race
	var messages [][]byte
	for range int(received.Load()) {
		select {
		case msg := <-done:
			messages = append(messages, msg)
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for forwarded datagram")
		}
	}

	if got := received.Load(); got != 2 {
		t.Errorf("forwarded %d datagrams, want 2", got)
	}

	if len(messages) >= 1 {
		if string(messages[0]) != "page.views:1|c" {
			t.Errorf("messages[0] = %q, want %q", messages[0], "page.views:1|c")
		}
	}
	if len(messages) >= 2 {
		// Raw datagrams are forwarded as-is, including \r\n
		if string(messages[1]) != "fuel.level:0.5|g\r\nsong.length:240|h" {
			t.Errorf("messages[1] = %q, want %q", messages[1], "fuel.level:0.5|g\r\nsong.length:240|h")
		}
	}
}

func TestServerStops(t *testing.T) {
	port := findAvailablePort(t)
	addr := "127.0.0.1:" + port

	srv := NewServer(addr, func(msg []byte) error { return nil }, nil, log.Default())
	go func() {
		_ = srv.Listen()
	}()
	time.Sleep(100 * time.Millisecond)

	if err := srv.Stop(); err != nil {
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
	defer func() { _ = conn.Close() }()
	return fmt.Sprintf("%d", conn.LocalAddr().(*net.UDPAddr).Port)
}
