// Package route provides HTTP transport layer functionalities.
package route

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"z-chat/internal/handlers"
	"z-chat/internal/hub"
)

// NewRouter creates a new HTTP router with the necessary routes and middleware.
func NewRouter(_ *hub.Hub, wsHandler *handlers.WebSocketHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	router.Get("/ws", wsHandler.ServeWS)
	return router
}
