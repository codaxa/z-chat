package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// AuthServicer defines the authentication service interface
type AuthServicer interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) (string, error)
}

// UserHandler handles HTTP requests related to user operations.
// It processes user-related endpoints and interacts with the authentication
// service to perform operations such as user registration, login, and
// profile management.
type UserHandler struct {
	authService AuthServicer
}

// NewUserHandler creates and returns a new UserHandler instance with the provided
// authentication service. It's responsible for handling HTTP requests related to user
// operations such as authentication, registration, and user management.
func NewUserHandler(authService AuthServicer) *UserHandler {
	return &UserHandler{authService: authService}
}

type tokenResponse struct {
	Token string `json:"token"`
}

// Structs for input

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Register handles HTTP requests for user registration.
// It accepts a JSON payload containing username and password,
// validates that both fields are not empty, and attempts to register
// the user via the authentication service.
//
// On success, it returns HTTP 201 Created.
// On error, it returns either HTTP 400 Bad Request for invalid input
// or HTTP 500 Internal Server Error for service failures.
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("ERROR parsing request body: %v", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	log.Printf("Parsed registration request for user: %s", req.Username)

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	log.Printf("Registration request received for username: %s", req.Username)
	err := h.authService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		log.Printf("ERROR in registration: %v", err)
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Login handles user authentication requests. It accepts a JSON payload with
// username and password credentials, validates them through the auth service,
// and returns a JWT token upon successful authentication. Returns 400 for invalid
// request format, 401 for invalid credentials, and 500 for internal errors.
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	resp := tokenResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("ERROR encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
