// Package handlers provides HTTP handlers for managing chat messages.
package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"z-chat/internal/domain/repository"
)

// MessageHandler manages HTTP requests related to messages in the chat application.
type MessageHandler struct {
	repo repository.MessageRepository
}

// NewMessageHandler creates a new MessageHandler with the provided message repository.
func NewMessageHandler(repo repository.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

// GetMessagesByRoom handles HTTP requests to retrieve all messages from a specific chat room.
func (h *MessageHandler) GetMessagesByRoom(w http.ResponseWriter, r *http.Request) {
	// Extract roomID from URL parameters
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)

	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Retrieve messages from repository
	messages, err := h.repo.GetMessagesByRoom(r.Context(), roomID, limit, offset)
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
