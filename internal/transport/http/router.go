// Package http provides HTTP transport layer functionalities.
package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"z-chat/internal/handlers"
)

// NewRouter creates a new HTTP router with the necessary routes and middleware.
func NewRouter(wsHandler *handlers.WebSocketHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	router.Get("/ws", wsHandler.ServeWS)
	return router
}
