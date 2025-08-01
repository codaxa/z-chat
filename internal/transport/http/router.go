// Package http provides HTTP transport layer functionalities.
package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"z-chat/internal/handlers"
	"z-chat/internal/services"
	middlewares "z-chat/internal/transport/http/middleware"
)

// NewRouter creates a new HTTP router with the necessary routes and middleware.
func NewRouter(wsHandler *handlers.WebSocketHandler, authService *services.AuthService) *chi.Mux {
	router := chi.NewRouter()
	userHandler := NewUserHandler(authService)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/api/health", handlers.HealthHandler)
	router.Get("/ws", middlewares.Authenticate(http.HandlerFunc(wsHandler.ServeWS)).ServeHTTP)
	router.Post("/api/register", userHandler.Register)
	router.Post("/api/login", userHandler.Login)

	return router
}
