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
func NewRouter(wsHandler *handlers.WebSocketHandler, messageHandler *handlers.MessageHandler, roomHandler *handlers.RoomHandler, authService *services.AuthService) *chi.Mux {
	router := chi.NewRouter()
	userHandler := NewUserHandler(authService)

	authMiddleware := middlewares.NewAuthMiddleware(authService)

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Public routes
	router.Get("/api/health", handlers.HealthHandler)
	router.Post("/api/register", userHandler.Register)
	router.Post("/api/login", userHandler.Login)

	// WebSocket routes (authenticated)
	router.Get("/ws/{roomID}", authMiddleware.Authenticate(http.HandlerFunc(wsHandler.ServeWS)).ServeHTTP)

	// API routes (authenticated)
	router.Route("/api", func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		r.Route("/rooms", func(r chi.Router) {
			r.Get("/", roomHandler.GetRooms)
			r.Post("/", roomHandler.CreateRoom)
			r.Get("/{roomID}", roomHandler.GetRoomByID)
			r.Delete("/{roomID}", roomHandler.DeleteRoom)
			r.Get("/{roomID}/messages", messageHandler.GetMessagesByRoom)
			r.Get("/mine", roomHandler.GetUserRooms)

			// Member management
			r.Post("/{roomID}/members", roomHandler.AddMember)
			r.Delete("/{roomID}/members/{userID}", roomHandler.RemoveMember)
			r.Get("/{roomID}/members", roomHandler.GetMembers)

			// Admin management
			r.Post("/{roomID}/admins", roomHandler.AddAdmin)
			r.Delete("/{roomID}/admins/{userID}", roomHandler.RemoveAdmin)
			r.Get("/{roomID}/admins", roomHandler.GetAdmins)
		})
	})

	return router
}
