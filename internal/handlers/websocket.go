// Package handlers provides HTTP handlers for the chat application.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	appContext "z-chat/internal/context"
	"z-chat/internal/domain/repository"
	"z-chat/internal/hub"
)

// WebSocketHandler handles WebSocket connections for the chat application.
type WebSocketHandler struct {
	HubManager        *hub.Manager
	MessageRepository repository.MessageRepository
	RoomRepository    repository.RoomRepository
	UserRepository    repository.UserRepository
}

// NewWebSocketHandler creates a new WebSocketHandler instance.
func NewWebSocketHandler(hubManager *hub.Manager, messageRepo repository.MessageRepository, roomRepo repository.RoomRepository, userRepo repository.UserRepository) *WebSocketHandler {
	return &WebSocketHandler{
		HubManager:        hubManager,
		MessageRepository: messageRepo,
		RoomRepository:    roomRepo,
		UserRepository:    userRepo,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
}

// ServeWS handles WebSocket upgrade requests and manages client connections.
func (h *WebSocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	username, userID := h.extractUserInfo(r)
	if username == "" || userID == "" {
		http.Error(w, "authentication required", http.StatusUnauthorized)
		return
	}
	fmt.Println("Extracted username:", username, "userID:", userID)

	// Get room based on URL pattern
	roomID, err := h.getRoom(r, userID)
	fmt.Print("Extracted roomID:", roomID)
	if err != nil {
		log.Printf("Failed to get room: %v", err)
		http.Error(w, "failed to access room", http.StatusInternalServerError)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	roomHub := h.HubManager.GetOrCreateHub(roomID)
	client := hub.NewClient(roomHub, conn, username)

	if client == nil {
		log.Printf("Failed to create client")
		return
	}

	h.sendRecentMessages(r.Context(), client, roomID)

	roomHub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (h *WebSocketHandler) extractUserInfo(r *http.Request) (username, userID string) {
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
		}

		if ok && claimsMap != nil {
			if usernameVal, exists := claimsMap["username"]; exists {
				if usernameStr, ok := usernameVal.(string); ok {
					username = usernameStr
				}
			}
			if userIDVal, exists := claimsMap["user_id"]; exists {
				if userIDStr, ok := userIDVal.(string); ok {
					userID = userIDStr
				}
			}
		}
	}
	return username, userID
}

func (h *WebSocketHandler) getRoom(r *http.Request, userID string) (roomID string, err error) {
	roomID = chi.URLParam(r, "roomID")

	if roomID == "" {
		return "", fmt.Errorf("roomID is required")
	}

	room, err := h.RoomRepository.GetRoomByID(r.Context(), roomID)
	if err != nil {
		return "", fmt.Errorf("failed to get room: %w", err)
	}

	if room == nil {
		return "", fmt.Errorf("room not found")
	}

	// Then check membership
	isMember, err := h.RoomRepository.IsRoomMember(r.Context(), roomID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to check room membership: %w", err)
	}

	if !isMember {
		return "", fmt.Errorf("user is not a member of this room")
	}

	return roomID, nil
}

func (h *WebSocketHandler) sendRecentMessages(ctx context.Context, client *hub.Client, roomID string) {
	messages, err := h.MessageRepository.GetMessagesByRoom(ctx, roomID, 50, 0)
	if err != nil {
		log.Printf("Failed to get recent messages: %v", err)
		return
	}

	for i := len(messages) - 1; i >= 0; i-- {
		msgBytes, err := json.Marshal(messages[i])
		if err != nil {
			log.Printf("Failed to marshal message: %v", err)
			continue
		}

		if !client.SendMessage(msgBytes) {
			log.Printf("Failed to send historical message to client")
			break
		}
	}
}
