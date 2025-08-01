// Package handlers provides HTTP handlers for the chat application.
package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	appContext "z-chat/internal/context"
	"z-chat/internal/hub"
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
	var username string

	if claims := r.Context().Value(appContext.UserClaimsKey); claims != nil {

		var claimsMap map[string]interface{}
		var ok bool

		switch v := claims.(type) {
		case map[string]interface{}:
			claimsMap = v
			ok = true
		case *map[string]interface{}:
			claimsMap = *v
			ok = true
		case jwt.MapClaims:
			claimsMap = map[string]interface{}(v)
			ok = true
		case *jwt.MapClaims:
			claimsMap = map[string]interface{}(*v)
			ok = true
		default:
			log.Printf("Unexpected claims type: %T", v)
		}

		if ok && claimsMap != nil {
			if usernameVal, exists := claimsMap["username"]; exists {
				if usernameStr, strOk := usernameVal.(string); strOk {
					username = usernameStr
				} else {
					log.Printf("Username value is not a string: %v (type: %T)", usernameVal, usernameVal)
				}
			} else {
				log.Printf("Username key not found in claims: %+v", claimsMap)
			}
		} else {
			log.Printf("Failed to convert claims to map[string]interface{}")
		}
	} else {
		log.Printf("No claims found in context with key: %+v", appContext.UserClaimsKey)
	}

	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := hub.NewClient(h.hub, conn, username)
	if client == nil {
		log.Printf("Failed to create client")
		return
	}

	h.hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
