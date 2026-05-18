package forwarder

import (
	"net"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

type recvResult struct {
	data []byte
	err  error
}

func TestForwarderForwardsDatagrams(t *testing.T) {
	// Start a UDP server to receive forwarded datagrams
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}
	upstream, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = upstream.Close() }()

	upstreamAddr := upstream.LocalAddr().String()
	done := make(chan recvResult, 1)
	buf := make([]byte, 1024)

	// Read forwarded datagrams in background
	go func() {
		_ = upstream.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, _, err := upstream.ReadFromUDP(buf)
		if err != nil {
			done <- recvResult{err: err}
			return
		}
		cp := make([]byte, n)
		copy(cp, buf[:n])
		done <- recvResult{data: cp}
	}()

	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: upstreamAddr,
	}

	if err := fwd.Start(); err != nil {
		t.Fatalf("failed to start forwarder: %v", err)
	}

	fwd.Forward([]byte("page.views:1|c"))

	select {
	case result := <-done:
		if result.err != nil {
			t.Fatalf("read error: %v", result.err)
		}
		if string(result.data) != "page.views:1|c" {
			t.Errorf("forwarded %q, want %q", string(result.data), "page.views:1|c")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for forwarded datagram")
	}

	if err := fwd.Stop(); err != nil {
		t.Fatalf("failed to stop forwarder: %v", err)
	}
}

func TestForwarderStops(t *testing.T) {
	// Start a real upstream to dial to
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}
	upstream, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = upstream.Close() }()

	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: upstream.LocalAddr().String(),
	}

	if err := fwd.Start(); err != nil {
		t.Fatalf("failed to start forwarder: %v", err)
	}

	if err := fwd.Stop(); err != nil {
		t.Errorf("Stop() returned error: %v", err)
	}
}

func TestForwarderNoopWithoutStart(t *testing.T) {
	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: "127.0.0.1:0",
	}

	// Forward without starting should not panic
	fwd.Forward([]byte("page.views:1|c"))
}

func TestForwarderNilLogger(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}
	upstream, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = upstream.Close() }()

	fwd := &Forwarder{
		Address: upstream.LocalAddr().String(),
	}

	if err := fwd.Start(); err != nil {
		t.Fatalf("failed to start forwarder: %v", err)
	}

	// Should not panic with nil logger
	fwd.Forward([]byte("page.views:1|c"))

	if err := fwd.Stop(); err != nil {
		t.Errorf("Stop() returned error: %v", err)
	}
}

func TestForwarderInvalidAddress(t *testing.T) {
	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: "invalid:address",
	}

	if err := fwd.Start(); err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}

func TestForwarderDoubleStart(t *testing.T) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to resolve addr: %v", err)
	}
	upstream, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer func() { _ = upstream.Close() }()

	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: upstream.LocalAddr().String(),
	}

	if err := fwd.Start(); err != nil {
		t.Fatalf("failed to start forwarder: %v", err)
	}
	defer func() { _ = fwd.Stop() }()

	if err := fwd.Start(); err == nil {
		t.Error("expected error on double start, got nil")
	}
}
