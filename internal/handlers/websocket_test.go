package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"z-chat/internal/domain/models"
	"z-chat/internal/hub"

	"github.com/gorilla/websocket"
)

// Mock message repository for testing
type mockMessageRepository struct{}

func (m *mockMessageRepository) CreateMessage(_ context.Context, _ *models.Message) error {
	return nil
}

func (m *mockMessageRepository) GetMessageByID(_ context.Context, _ string) (*models.Message, error) {
	return nil, nil
}

func (m *mockMessageRepository) SaveMessage(_ models.Message) error {
	return nil
}

func (m *mockMessageRepository) GetMessages() ([]models.Message, error) {
	return nil, nil
}

func TestNewWebSocketHandler(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	handler := NewWebSocketHandler(h)

	if handler == nil {
		t.Fatal("expected WebSocketHandler instance, got nil")
	}
	if handler.hub != h {
		t.Error("expected handler to have correct hub reference")
	}
}

func TestWebSocketUpgrade(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	go h.Run()

	handler := NewWebSocketHandler(h)
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	// Add ?username=test to the URL
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?username=test"

	// Test successful WebSocket connection
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}

	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Fatal(err)
	}

	// Verify client is registered
	if h.ClientCount() != 1 {
		t.Errorf("expected 1 client registered, got %d", h.ClientCount())
	}
}

func TestWebSocketHandler_ServeWS_InvalidUpgrade(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	handler := NewWebSocketHandler(h)

	// Create a regular HTTP request (not WebSocket upgrade)
	req := httptest.NewRequest("GET", "/ws", nil)
	w := httptest.NewRecorder()

	handler.ServeWS(w, req)

	// Should return bad request since it's not a valid WebSocket upgrade
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
