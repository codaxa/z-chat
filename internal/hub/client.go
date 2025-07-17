// Package hub provides a client for interacting with the hub service.
package hub

import (
	"encoding/json"
	"log"
	"z-chat/internal/domain/models"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connected to the hub.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

// NewClient creates a new Client instance with the provided hub and WebSocket connection.
func NewClient(hub *Hub, conn *websocket.Conn, username string) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
	}
}

// ReadPump handles reading messages from the WebSocket connection and broadcasting them to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.Unregister <- c
		if err := c.conn.Close(); err != nil {
			log.Printf("can't close the connection: %v", err)
		}
	}()
	for {
		_, messageBytes, err := c.conn.ReadMessage()

		if err != nil {
			break
		}
		var message models.Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("error parsing message: %v", err)
			continue
		}
		if err := message.Validate(); err != nil {
			log.Printf("invalid message: %v", err)
			continue
		}
		// Set the sender to the current client's username before broadcasting
		message.Sender = c.username
		c.hub.broadcast <- message
	}
}

// WritePump handles writing messages from the client's send channel to the WebSocket connection.
func (c *Client) WritePump() {
	defer func() {
		c.hub.Unregister <- c
		if err := c.conn.Close(); err != nil {
			log.Printf("can't close the connection: %v", err)
		}
	}()
	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
