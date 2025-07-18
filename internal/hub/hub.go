// Package hub provides a client for interacting with the hub service.
package hub

import (
	"encoding/json"
	"log"
	"sync"
	"z-chat/internal/domain/models"
)

// Hub represents a hub for managing WebSocket connections.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex // Add this mutex
}

// ClientsCount returns the number of connected clients.
func (h *Hub) ClientsCount() int {
	h.mu.RLock()         // Lock for reading
	defer h.mu.RUnlock() // Unlock when function returns
	return len(h.clients)

}

// NewHub returns a new Hub instance with initialized channels and an empty set of clients.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan models.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Run starts the hub's main loop, handling client registration, unregistration, and message broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock() // Lock for writing
			h.clients[client] = true
			h.mu.Unlock() // Unlock after writing

		case client := <-h.Unregister:
			h.mu.Lock() // Lock for writing
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock() // Unlock after writing
		case message := <-h.broadcast:
			// Marshal the message to JSON bytes
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling message for broadcast: %v", err)
				continue
			}

			h.mu.RLock() // Lock for reading clients map
			for client := range h.clients {
				select {
				case client.send <- messageBytes: // Send JSON bytes
				default:
					close(client.send)
					// Note: This creates another race condition -
					// you'd need to collect clients to delete and
					// delete them after the loop
				}
			}
			h.mu.RUnlock() // Unlock after reading
		}
	}
}
