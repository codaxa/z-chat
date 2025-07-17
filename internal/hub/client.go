// Package hub provides a client for interacting with the hub service.
package hub

import (
	"github.com/gorilla/websocket"
	"log"
)

// Client represents a WebSocket client connected to the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
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
			log.Printf("can't close the connection: %v", err)
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
