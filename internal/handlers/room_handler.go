// Package handlers provides HTTP handlers for managing chat rooms.
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
	appContext "z-chat/internal/context"
	"z-chat/internal/domain/models"
	"z-chat/internal/domain/repository"
)

// RoomHandler manages HTTP requests related to chat rooms.
type RoomHandler struct {
	repo repository.RoomRepository
}

// NewRoomHandler creates a new RoomHandler with the provided room repository.
func NewRoomHandler(repo repository.RoomRepository) *RoomHandler {
	return &RoomHandler{repo: repo}
}

// CreateRoom handles HTTP requests to create a new chat room.
func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract user ID from JWT claims
	userID := h.extractUserID(r)
	fmt.Println("Extracted userID:", userID)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	room := &models.Room{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatedBy: userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.CreateRoom(r.Context(), room); err != nil {
		http.Error(w, fmt.Sprintf("Failed to create room, error: %v", err), http.StatusInternalServerError)
		return
	}

	// Add creator as admin
	if err := h.repo.AddRoomAdmin(r.Context(), room.ID, userID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to add room admin, error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(room); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode room response, error: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetRooms handles HTTP requests to retrieve a paginated list of chat rooms.
func (h *RoomHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
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

	rooms, err := h.repo.GetRooms(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, "failed to get rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode rooms response, error: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetRoomByID handles HTTP requests to retrieve a chat room by its ID.
func (h *RoomHandler) GetRoomByID(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	room, err := h.repo.GetRoomByID(r.Context(), roomID)
	if err != nil {
		http.Error(w, "failed to get room", http.StatusInternalServerError)
		return
	}

	if room == nil {
		http.Error(w, "room not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(room); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode room response, error: %v", err), http.StatusInternalServerError)
		return
	}
}

// DeleteRoom handles HTTP requests to delete a chat room by its ID.
func (h *RoomHandler) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteRoom(r.Context(), roomID); err != nil {
		http.Error(w, "failed to delete room", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddMember handles HTTP requests to add a member to a chat room.
func (h *RoomHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	currentUserID := h.extractUserID(r)

	isAdmin, err := h.repo.IsRoomAdmin(r.Context(), roomID, currentUserID)
	if err != nil || !isAdmin {
		http.Error(w, "Only room admins can add members", http.StatusForbidden)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.AddRoomMember(r.Context(), roomID, req.UserID); err != nil {
		http.Error(w, "Failed to add member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveMember handles HTTP requests to remove a member from a chat room.
func (h *RoomHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.RemoveRoomMember(r.Context(), roomID, userID); err != nil {
		http.Error(w, "Failed to remove member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddAdmin handles HTTP requests to add an admin to a chat room.
func (h *RoomHandler) AddAdmin(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.AddRoomAdmin(r.Context(), roomID, req.UserID); err != nil {
		http.Error(w, "Failed to add admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetMembers handles HTTP requests to retrieve all members of a chat room.
func (h *RoomHandler) GetMembers(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	members, err := h.repo.GetRoomMembers(r.Context(), roomID)
	if err != nil {
		http.Error(w, "failed to get members", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(members); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode members response, error: %v", err), http.StatusInternalServerError)
		return
	}
}

// GetAdmins handles HTTP requests to retrieve all admins of a chat room.
func (h *RoomHandler) GetAdmins(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	admins, err := h.repo.GetRoomAdmins(r.Context(), roomID)
	if err != nil {
		http.Error(w, "failed to get admins", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(admins); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode admins response, error: %v", err), http.StatusInternalServerError)
		return
	}
}

// RemoveAdmin handles HTTP requests to remove an admin from a chat room.
func (h *RoomHandler) RemoveAdmin(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if roomID == "" {
		http.Error(w, "roomID is required", http.StatusBadRequest)
		return
	}

	userID := chi.URLParam(r, "userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	if err := h.repo.RemoveRoomAdmin(r.Context(), roomID, userID); err != nil {
		http.Error(w, "Failed to remove admin", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserRooms handles HTTP requests to retrieve all rooms a user is a member of.
func (h *RoomHandler) GetUserRooms(w http.ResponseWriter, r *http.Request) {
	userID := h.extractUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	rooms, err := h.repo.GetUserRooms(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get user rooms", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(rooms); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode rooms response, error: %v", err), http.StatusInternalServerError)
		return
	}
}

// extractUserID extracts the user ID from the JWT claims in the request context.
func (h *RoomHandler) extractUserID(r *http.Request) (userID string) {
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
			if userIDVal, exists := claimsMap["user_id"]; exists {
				if userIDStr, ok := userIDVal.(string); ok {
					userID = userIDStr
				}
			}
		}
	}
	return userID
}
