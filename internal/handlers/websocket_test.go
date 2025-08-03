package handlers

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	appContext "z-chat/internal/context"
	"z-chat/internal/domain/models"
	"z-chat/internal/hub"
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := map[string]interface{}{
			"username": "testuser",
		}
		ctx := context.WithValue(r.Context(), appContext.UserClaimsKey, claims)
		r = r.WithContext(ctx)

		handler.ServeWS(w, r)
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("WebSocket connection failed: %v", err)
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("expected status %d, got %d", http.StatusSwitchingProtocols, resp.StatusCode)
	}

	if err := conn.Close(); err != nil {
		t.Logf("Error closing connection: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		t.Logf("Error closing response body: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
}

func TestWebSocketHandler_ServeWS_InvalidUpgrade(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	handler := NewWebSocketHandler(h)

	req := httptest.NewRequest("GET", "/ws", nil)

	claims := map[string]interface{}{
		"username": "testuser",
	}
	ctx := context.WithValue(req.Context(), appContext.UserClaimsKey, claims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeWS(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestWebSocketHandler_ServeWS_MissingUsername(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	handler := NewWebSocketHandler(h)

	// Create a request with WebSocket headers but no username in JWT claims
	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	// Add empty JWT context (no username)
	claims := map[string]interface{}{}
	ctx := context.WithValue(req.Context(), appContext.UserClaimsKey, claims)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeWS(w, req)

	// Should return bad request for missing username
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "username is required"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("expected body to contain %q, got %q", expectedBody, w.Body.String())
	}
}

func TestWebSocketHandler_ServeWS_NoJWTContext(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	handler := NewWebSocketHandler(h)

	// Create a request with WebSocket headers but no JWT context at all
	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	// No JWT context added

	w := httptest.NewRecorder()

	handler.ServeWS(w, req)

	// Should return bad request for missing username
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "username is required"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("expected body to contain %q, got %q", expectedBody, w.Body.String())
	}
}

func TestWebSocketHandler_ServeWS_InvalidJWTClaims(t *testing.T) {
	repo := &mockMessageRepository{}
	h := hub.NewHub(repo)
	handler := NewWebSocketHandler(h)

	// Create a request with WebSocket headers
	req := httptest.NewRequest("GET", "/ws", nil)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	// Add invalid JWT context (not a map)
	ctx := context.WithValue(req.Context(), appContext.UserClaimsKey, "invalid")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	handler.ServeWS(w, req)

	// Should return bad request for invalid claims
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	expectedBody := "username is required"
	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf("expected body to contain %q, got %q", expectedBody, w.Body.String())
	}
}
