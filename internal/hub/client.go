// Package hub provides a client for interacting with the hub service.
package hub

import "github.com/gorilla/websocket"

// Client represents a WebSocket client connected to the hub.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}
