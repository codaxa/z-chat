// Package handlers provides HTTP handlers for the chat application.
package handlers

import (
	"log"
	"net/http"
	"z-chat/internal/hub"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for the chat application.
type WebSocketHandler struct {
	hub *hub.Hub
}

// NewWebSocketHandler creates a new WebSocketHandler instance.
func NewWebSocketHandler(hub *hub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// ServeWS handles WebSocket upgrade requests and manages client connections.
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}()

	client := hub.NewClient(h.hub, conn)
	if client == nil {
		log.Printf("Failed to create client")
		return
	}

	h.hub.Register <- client

	go client.ReadPump()
	go client.WritePump()
}
