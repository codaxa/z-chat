// Package hub provides a client for interacting with the hub service.
package hub

import (
	"fmt"
	"testing"
	"time"
)

func waitForClients(h *Hub, expected int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if len(h.clients) == expected {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %d clients; got %d", expected, len(h.clients))
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("expected new Hub instance, got nil")
	}
	if len(hub.clients) != 0 {
		t.Errorf("expected empty clients map, got %d clients", len(hub.clients))
	}
	if hub.broadcast == nil || hub.register == nil || hub.unregister == nil {
		t.Error("expected non-nil channels for broadcast, register, and unregister")
	}
}

func TestHubRun(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{hub: hub, send: make(chan []byte, 1)}
	client2 := &Client{hub: hub, send: make(chan []byte, 1)}

	hub.register <- client1
	hub.register <- client2
	if err := waitForClients(hub, 2, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	hub.unregister <- client1
	if err := waitForClients(hub, 1, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	hub.broadcast <- []byte("test message")
	select {
	case msg := <-client2.send:
		if got := string(msg); got != "test message" {
			t.Errorf("expected 'test message', got %q", got)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for broadcast message")
	}
}
