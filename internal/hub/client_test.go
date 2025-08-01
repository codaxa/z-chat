package hub

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"z-chat/internal/domain/models"
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

func TestNewClient(t *testing.T) {
	repo := &mockMessageRepositoryClientTest{}
	h := NewHub(repo)
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
	h := NewHub(repo)
	go h.Run()

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

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
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

	time.Sleep(10 * time.Millisecond)

	close(client.send)
}

func TestClient_ReadPump_Unregisters(t *testing.T) {
	repo := &mockMessageRepositoryClientTest{}
	h := NewHub(repo)
	go h.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upgrader := websocket.Upgrader{

			CheckOrigin: func(_ *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		if err := conn.Close(); err != nil {
			t.Logf("Error closing connection: %v", err)
		}

	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")

	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Logf("Error closing body: %v", err)
		}
	}()

	client := NewClient(h, conn, "testuser")
	h.Register <- client

	if err := waitForClients(h, 1, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

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
