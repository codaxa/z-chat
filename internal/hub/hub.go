// Package hub provides a client for interacting with the hub service.
package hub

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"z-chat/internal/domain/models"
	"z-chat/internal/domain/repository"
)

// Hub represents a hub for managing WebSocket connections.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
	repo       repository.MessageRepository
}

// ClientCount returns the current number of clients connected to the hub.p.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// Clients provides access to the connected clients in the Hub.
func (h *Hub) Clients() {
	panic("unimplemented")
}

// NewHub returns a new Hub instance with initialized channels and an empty set of clients.
func NewHub(repo repository.MessageRepository) *Hub {
	return &Hub{
		broadcast:  make(chan models.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		repo:       repo,
	}
}

// Run starts the hub's main loop, handling client registration, unregistration, and message broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
		case message := <-h.broadcast:
			messageBytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshaling message for broadcast: %v", err)
				continue
			}
			if err := h.repo.CreateMessage(context.Background(), &message); err != nil {
				log.Printf("error saving message: %v", err)
			}
			h.mu.RLock()
			var failedClients []*Client
			for client := range h.clients {
				select {
				case client.send <- messageBytes:
				default:
					failedClients = append(failedClients, client)
				}
			}
			h.mu.RUnlock()

			if len(failedClients) > 0 {
				h.mu.Lock()
				for _, client := range failedClients {
					if _, ok := h.clients[client]; ok {
						delete(h.clients, client)
						close(client.send)
					}
				}
				h.mu.Unlock()
			}
		}
	}
}
