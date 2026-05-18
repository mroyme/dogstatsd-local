package messages

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/charmbracelet/log"
)

func TestHandlerDispatchesMessages(t *testing.T) {
	var received atomic.Int32
	handler := &Handler{
		Logger:     log.Default(),
		PoolSize:   4,
		BufferSize: 100,
		Out: func(msg DogStatsDMessage) error {
			received.Add(1)
			return nil
		},
	}
	handler.New()

	for range 100 {
		err := handler.Handle([]byte("page.views:1|c"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	handler.Stop()

	if got := received.Load(); got != 100 {
		t.Errorf("received %d messages, want 100", got)
	}
}

func TestHandlerParsesAllMessageTypes(t *testing.T) {
	var mu sync.Mutex
	got := make(map[DogStatsDMessageType]int)

	handler := &Handler{
		Logger:     log.Default(),
		PoolSize:   4,
		BufferSize: 100,
		Out: func(msg DogStatsDMessage) error {
			mu.Lock()
			got[msg.Type()]++
			mu.Unlock()
			return nil
		},
	}
	handler.New()

	inputs := []struct {
		raw string
		typ DogStatsDMessageType
	}{
		{"page.views:1|c", MetricMessageType},
		{"fuel.level:0.5|g", MetricMessageType},
		{"_sc|db_check|0|#env:prod", ServiceCheckMessageType},
		{"_sc|cache|2|m:down", ServiceCheckMessageType},
		{"_e{5,4}:title|text", EventMessageType},
		{"_e{6,14}:title1|text with pipes", EventMessageType},
	}

	for _, input := range inputs {
		err := handler.Handle([]byte(input.raw))
		if err != nil {
			t.Fatalf("unexpected error handling %q: %v", input.raw, err)
		}
	}

	handler.Stop()

	want := map[DogStatsDMessageType]int{
		MetricMessageType:       2,
		ServiceCheckMessageType: 2,
		EventMessageType:        2,
	}

	for typ, wantCount := range want {
		if gotCount := got[typ]; gotCount != wantCount {
			t.Errorf("message type %v: got %d, want %d", typ, gotCount, wantCount)
		}
	}
}

func TestHandlerLogsParseErrors(t *testing.T) {
	logger := log.Default()
	logger.SetLevel(log.DebugLevel)

	handler := &Handler{
		Logger:     logger,
		PoolSize:   1,
		BufferSize: 10,
		Out: func(msg DogStatsDMessage) error {
			return nil
		},
	}
	handler.New()

	// Send an invalid message — should be logged, not crash the pool
	err := handler.Handle([]byte("invalid"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Send a valid message after — pool should still work
	err = handler.Handle([]byte("page.views:1|c"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	handler.Stop()

	// The key assertion is that the handler didn't panic
	// and the pool kept processing after the parse error.
}

func TestHandlerCapacityExceeded(t *testing.T) {
	// Use a handler that blocks to prevent workers from draining the channel.
	blockCh := make(chan struct{})
	var processed atomic.Int32

	handler := &Handler{
		Logger:     log.Default(),
		PoolSize:   1,
		BufferSize: 2,
		Out: func(msg DogStatsDMessage) error {
			<-blockCh
			processed.Add(1)
			return nil
		},
	}
	handler.New()

	// The channel capacity is BufferSize (2). The worker picks up 1 message
	// immediately, so we can fill 2 more into the channel. After that,
	// the non-blocking send should return an error.
	var overflow bool
	for range 10 {
		err := handler.Handle([]byte("page.views:1|c"))
		if err != nil {
			overflow = true
			break
		}
	}

	if !overflow {
		t.Error("expected pool capacity error, all messages accepted")
	}

	// Unblock workers and stop
	close(blockCh)
	handler.Stop()
}

func TestHandlerStop(t *testing.T) {
	handler := &Handler{
		Logger:     log.Default(),
		PoolSize:   4,
		BufferSize: 100,
		Out: func(msg DogStatsDMessage) error {
			return nil
		},
	}
	handler.New()

	// Stop should not block or panic
	handler.Stop()
}
