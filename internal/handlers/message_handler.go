package handlers

import (
	"encoding/json"
	"net/http"
	"z-chat/internal/domain/repository"

	"github.com/go-chi/chi/v5"
)

// MessageHandler manages HTTP requests related to messages in the chat application.
// It uses a MessageRepository to interact with the message data storage layer.
type MessageHandler struct {
	repo repository.MessageRepository
}

// NewMessageHandler creates a new MessageHandler with the provided message repository.
// It returns a pointer to the new MessageHandler instance.
func NewMessageHandler(repo repository.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

// GetMessagesByRoom handles HTTP requests to retrieve all messages from a specific chat room.
// It extracts the roomID from the URL parameters, fetches messages using the repository,
// and returns them as a JSON response.
func (h *MessageHandler) GetMessagesByRoom(w http.ResponseWriter, r *http.Request) {
	// Extract roomID from URL parameters
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	// Retrieve messages from repository
	messages, err := h.repo.GetMessagesByRoom(r.Context(), roomID)
	if err != nil {
		http.Error(w, "failed to get messages: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send messages as JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "failed to encode messages: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
