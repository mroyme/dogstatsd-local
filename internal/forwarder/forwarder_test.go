package forwarder

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/charmbracelet/log"
)

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

	var received atomic.Int32
	buf := make([]byte, 1024)

	// Read forwarded datagrams in background
	go func() {
		for {
			_ = upstream.SetReadDeadline(time.Now().Add(2 * time.Second))
			n, _, err := upstream.ReadFromUDP(buf)
			if err != nil {
				return
			}
			if string(buf[:n]) == "page.views:1|c" {
				received.Add(1)
			}
		}
	}()

	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: upstreamAddr,
	}

	if err := fwd.Start(); err != nil {
		t.Fatalf("failed to start forwarder: %v", err)
	}

	fwd.Forward([]byte("page.views:1|c"))

	time.Sleep(200 * time.Millisecond)

	if err := fwd.Stop(); err != nil {
		t.Fatalf("failed to stop forwarder: %v", err)
	}

	if got := received.Load(); got != 1 {
		t.Errorf("forwarded %d datagrams, want 1", got)
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

func TestForwarderInvalidAddress(t *testing.T) {
	fwd := &Forwarder{
		Logger:  log.Default(),
		Address: "invalid:address",
	}

	if err := fwd.Start(); err == nil {
		t.Error("expected error for invalid address, got nil")
	}
}