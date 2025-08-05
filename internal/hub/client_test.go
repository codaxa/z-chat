package hub

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"z-chat/internal/domain/models"

	"github.com/gorilla/websocket"
)

// MockMessageRepository is a mock implementation for testing
type mockMessageRepositoryClientTest struct{}

func (m *mockMessageRepositoryClientTest) CreateMessage(_ context.Context, _ *models.Message) error {
	return nil
}

func (m *mockMessageRepositoryClientTest) GetMessageByID(_ context.Context, _ string) (*models.Message, error) {
	return nil, nil
}

func (m *mockMessageRepositoryClientTest) SaveMessage(_ models.Message) error {
	return nil
}

func (m *mockMessageRepositoryClientTest) GetMessages() ([]models.Message, error) {
	return nil, nil
}

func (m *mockMessageRepositoryClientTest) GetMessagesByRoom(_ context.Context, _ string, _ int, _ int) ([]*models.Message, error) {
	return []*models.Message{}, nil
}

func TestNewClient(t *testing.T) {
	repo := &mockMessageRepositoryClientTest{}
	h := NewHub(repo, "test")
	conn := &websocket.Conn{} // Mock connection

	client := NewClient(h, conn, "testuser")

	if client == nil {
		t.Fatal("expected Client instance, got nil")
	}
	if client.hub != h {
		t.Error("expected client to have correct hub reference")
	}
	if client.conn != conn {
		t.Error("expected client to have correct connection reference")
	}
	if client.send == nil {
		t.Error("expected client send channel to be initialized")
	}
}

func TestClient_WritePump(t *testing.T) {
	repo := &mockMessageRepositoryClientTest{}
	h := NewHub(repo, "test")
	go h.Run()

	messageReceived := make(chan []byte, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(_ *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				t.Logf("Error closing connection: %v", err)
			}
		}()

		// Read one message and signal
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}
		messageReceived <- msg
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not connect to websocket: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("Error closing connection: %v", err)
		}
	}()

	if err := resp.Body.Close(); err != nil {
		t.Logf("Error closing connection: %v", err)
	}

	client := NewClient(h, conn, "testuser")
	go client.WritePump()

	testMessage := []byte("test message")
	client.send <- testMessage

	select {
	case received := <-messageReceived:
		if string(received) != string(testMessage) {
			t.Errorf("handler received wrong message: got %s want %s", received, testMessage)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timed out waiting for message")
	}

	close(client.send)
}

func TestClient_ReadPump_Unregisters(t *testing.T) {
	repo := &mockMessageRepositoryClientTest{}
	h := NewHub(repo, "test")
	go h.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{

			CheckOrigin: func(_ *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		// Immediately close the connection to trigger an error on the client's ReadPump
		if err := conn.Close(); err != nil {
			t.Logf("Error closing connection: %v", err)
		}

	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not connect to websocket: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Error closing body: %v", err)
		}
		// The client will close the connection in its ReadPump defer.
		// Closing it here again might cause a "close of closed connection" panic.
		// conn.Close()
	}()

	client := NewClient(h, conn, "testuser")
	h.Register <- client

	// Wait for registration to be processed
	time.Sleep(50 * time.Millisecond)

	go client.ReadPump()

	if err := waitForClients(h, 0, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
}

// Helper function to wait for expected number of clients
func waitForClients(h *Hub, expected int, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if len(h.clients) == expected {
			return nil
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("timed out waiting for %d clients, got %d", expected, len(h.clients))
}
