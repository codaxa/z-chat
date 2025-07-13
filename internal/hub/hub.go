// Package hub provides a client for interacting with the hub service.
package hub

import (
	"log"

	"github.com/gorilla/websocket"
)

// Hub represents a hub for managing WebSocket connections.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// ClientsCount returns the number of connected clients.
func (h *Hub) ClientsCount() int {
	return len(h.clients)
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
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
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// NewClient creates a new Client instance with the provided hub and WebSocket connection.
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
}

// ReadPump handles reading messages from the WebSocket connection and broadcasting them to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		if err := c.conn.Close(); err != nil {
			log.Print("can't close the connection %v", err)
		}

	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		c.hub.broadcast <- message
	}
}

// WritePump handles writing messages from the client's send channel to the WebSocket connection.
func (c *Client) WritePump() {
	defer func() {
		if err := c.conn.Close(); err != nil {
			log.Print("can't close the connection %v", err)
		}
	}()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
