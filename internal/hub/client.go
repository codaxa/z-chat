// Package hub provides a client for interacting with the hub service.
package hub

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"z-chat/internal/domain/models"
)

// Client represents a WebSocket client connected to the hub.
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	username string
}

// NewClient returns a new Client associated with the given hub, WebSocket connection, and username.
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

		message.Sender = c.username
		message.RoomID = c.hub.roomID

		// Populate Room field - you need access to room repository
		if c.hub.roomRepo != nil {
			room, err := c.hub.roomRepo.GetRoomByID(context.Background(), c.hub.roomID)
			if err == nil && room != nil {
				message.Room = *room
			}
		}

		if err := message.Validate(); err != nil {
			log.Printf("invalid message: %v", err)
			continue
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

// SendMessage method to send messages
func (c *Client) SendMessage(message []byte) bool {
	select {
	case c.send <- message:
		return true
	default:
		return false
	}
}

// Username getter methods
func (c *Client) Username() string {
	return c.username
}
