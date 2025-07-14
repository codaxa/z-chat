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
		if h.ClientsCount() == expected {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for %d clients; got %d", expected, h.ClientsCount())
}

func TestNewHub(t *testing.T) {
	hub := NewHub()
	if hub == nil {
		t.Fatal("expected new Hub instance, got nil")
	}
	if len(hub.clients) != 0 {
		t.Errorf("expected empty clients map, got %d clients", len(hub.clients))
	}
	if hub.broadcast == nil || hub.Register == nil || hub.Unregister == nil {
		t.Error("expected non-nil channels for broadcast, register, and unregister")
	}
}

func TestHubRun(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{hub: hub, send: make(chan []byte, 1)}
	client2 := &Client{hub: hub, send: make(chan []byte, 1)}

	hub.Register <- client1
	hub.Register <- client2
	if err := waitForClients(hub, 2, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	hub.Unregister <- client1
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

func TestHubBroadcastToMultipleClients(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{hub: hub, send: make(chan []byte, 1)}
	client2 := &Client{hub: hub, send: make(chan []byte, 1)}
	client3 := &Client{hub: hub, send: make(chan []byte, 1)}

	// Register all clients
	hub.Register <- client1
	hub.Register <- client2
	hub.Register <- client3
	if err := waitForClients(hub, 3, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// Broadcast message
	testMessage := []byte("broadcast test")
	hub.broadcast <- testMessage

	// Verify all clients receive the message
	clients := []*Client{client1, client2, client3}
	for i, client := range clients {
		select {
		case msg := <-client.send:
			if got := string(msg); got != "broadcast test" {
				t.Errorf("client %d: expected 'broadcast test', got %q", i+1, got)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("client %d: timeout waiting for broadcast message", i+1)
		}
	}
}

func TestHubUnregisterNonExistentClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{hub: hub, send: make(chan []byte, 1)}

	// Try to unregister a client that was never registered
	hub.Unregister <- client

	// Should not cause any issues
	time.Sleep(10 * time.Millisecond)

	if len(hub.clients) != 0 {
		t.Errorf("expected 0 clients, got %d", len(hub.clients))
	}
}

func TestHubChannelInitialization(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// Give the hub time to start
	time.Sleep(10 * time.Millisecond)

	// Test that all channels are properly initialized and can be used
	// Use goroutines to avoid blocking since these are unbuffered channels
	done := make(chan bool, 3)

	// Test broadcast channel
	go func() {
		hub.broadcast <- []byte("test")
		done <- true
	}()

	// Test Register channel
	client := &Client{hub: hub, send: make(chan []byte, 1)}
	go func() {
		hub.Register <- client
		done <- true
	}()

	// Test Unregister channel
	go func() {
		hub.Unregister <- client
		done <- true
	}()

	// Wait for all operations to complete
	timeout := time.After(100 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case <-done:
			// Operation completed successfully
		case <-timeout:
			t.Fatalf("timeout waiting for channel operation %d to complete", i+1)
		}
	}
}
