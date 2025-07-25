// Package hub provides a client for interacting with the hub service.
package hub

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	"z-chat/internal/domain/models"
)

// MockMessageRepository is a mock implementation for testing.
type MockMessageRepository struct{}

// CreateMessage implements models.MessageRepository.
func (m *MockMessageRepository) CreateMessage(_ context.Context, _ *models.Message) error {
	return nil
}

// GetMessageByID implements models.MessageRepository.
func (m *MockMessageRepository) GetMessageByID(_ context.Context, _ string) (*models.Message, error) {
	return nil, nil
}

func (m *MockMessageRepository) SaveMessage(_ models.Message) error { return nil }

// Removed duplicate waitForClients function to resolve redeclaration error.

func TestNewHub(t *testing.T) {
	repo := &MockMessageRepository{} // Use the local mock implementation
	hub := NewHub(repo)
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
	repo := &MockMessageRepository{}
	hub := NewHub(repo)
	go hub.Run()

	client1 := &Client{hub: hub, send: make(chan []byte, 1)}
	client2 := &Client{hub: hub, send: make(chan []byte, 1)}

	hub.Register <- client1
	hub.Register <- client2
	if err := waitForClients(hub, 2, 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	hub.Unregister <- client1

	if err := waitForClients(hub, 1, 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	hub.broadcast <- models.Message{Content: "test message"}
	select {
	case msg := <-client2.send:
		var receivedMessage models.Message
		if err := json.Unmarshal(msg, &receivedMessage); err != nil {
			t.Fatalf("failed to unmarshal message: %v", err)
		}
		if receivedMessage.Content != "test message" {
			t.Errorf("expected content 'test message', got %q", receivedMessage.Content)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timeout waiting for broadcast message")
	}
}

func TestHubBroadcastToMultipleClients(t *testing.T) {
	repo := &MockMessageRepository{}
	hub := NewHub(repo)
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
	testMessage := models.Message{Content: "broadcast test"}
	hub.broadcast <- testMessage

	// Verify all clients receive the message
	clients := []*Client{client1, client2, client3}
	for i, client := range clients {
		select {
		case msg := <-client.send:
			// Parse the JSON to verify the message content
			var receivedMessage models.Message
			if err := json.Unmarshal(msg, &receivedMessage); err != nil {
				t.Fatalf("client %d: failed to unmarshal message: %v", i+1, err)
			}
			if receivedMessage.Content != "broadcast test" {
				t.Errorf("client %d: expected content 'broadcast test', got %q", i+1, receivedMessage.Content)
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatalf("client %d: timeout waiting for broadcast message", i+1)
		}
	}
}

func TestHubUnregisterNonExistentClient(t *testing.T) {
	repo := &MockMessageRepository{}
	hub := NewHub(repo)
	go hub.Run()

	client := &Client{hub: hub, send: make(chan []byte, 1)}

	// Try to unregister a client that was never registered
	hub.Unregister <- client

	// Should not cause any issues
	time.Sleep(10 * time.Millisecond)

	if len(hub.clients) != 0 {
		t.Errorf("expected no clients, got %d", len(hub.clients))
	}
}

func TestHubChannelInitialization(t *testing.T) {
	repo := &MockMessageRepository{}
	hub := NewHub(repo)
	go hub.Run()

	// Give the hub time to start
	time.Sleep(10 * time.Millisecond)

	// Test that all channels are properly initialized and can be used
	// Use goroutines to avoid blocking since these are unbuffered channels
	done := make(chan bool, 3)

	// Test broadcast channel
	go func() {
		hub.broadcast <- models.Message{Content: "test"}
		done <- true
	}()

	// Test Register channel
	client := &Client{hub: hub, send: make(chan []byte, 1)}
	go func() {
		hub.Register <- client
		done <- true
	}()

	// Wait for client to be registered
	if err := waitForClients(hub, 1, 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}

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

	// Verify client was unregistered
	if err := waitForClients(hub, 0, 50*time.Millisecond); err != nil {
		t.Fatal(err)
	}
}
