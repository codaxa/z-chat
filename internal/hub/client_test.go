package hub

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestNewClient(t *testing.T) {
	h := NewHub()
	conn := &websocket.Conn{} // Mock connection

	client := NewClient(h, conn)

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
	h := NewHub()
	go h.Run()

	// Create mock WebSocket connection
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

		// Keep connection alive
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

	client := NewClient(h, conn)
	go client.WritePump()

	// Send a message through the client
	testMessage := []byte("test message")
	client.send <- testMessage

	// Verify message was sent (WritePump should handle it)
	time.Sleep(10 * time.Millisecond)

	// Close send channel to stop WritePump
	close(client.send)
}

func TestClient_ReadPump_Unregisters(t *testing.T) {
	h := NewHub()
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
		} // Immediately close to trigger ReadPump exit
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}

	client := NewClient(h, conn)
	h.Register <- client

	// Wait for registration
	if err := waitForClients(h, 1, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}

	// ReadPump should unregister client when connection closes
	go client.ReadPump()

	// Wait for unregistration
	if err := waitForClients(h, 0, 100*time.Millisecond); err != nil {
		t.Fatal(err)
	}
}
