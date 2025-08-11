// It implements the main entry point for the chat server application.
package main

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"z-chat/internal/config"
	"z-chat/internal/handlers"
	"z-chat/internal/hub"
	"z-chat/internal/services"
	"z-chat/internal/storage/postgres"
	route "z-chat/internal/transport/http"
)

// main initializes and starts the chat server, setting up HTTP endpoints and launching the chat hub.
func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	cfg := config.New()

	// Initialize database connection
	connConfig, err := pgxpool.ParseConfig(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to parse database configuration: %v", err)
	}

	dbConn, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(dbConn)
	messageRepo := postgres.NewMessageRepository(dbConn)
	roomRepo := postgres.NewRoomRepository(dbConn)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenDuration)

	// Initialize chat hub
	messageHandler := handlers.NewMessageHandler(messageRepo)
	roomHandler := handlers.NewRoomHandler(roomRepo)

	chatHub := hub.NewManager(messageRepo, roomRepo)
	wsHandler := handlers.NewWebSocketHandler(chatHub, messageRepo, roomRepo, userRepo)
	router := route.NewRouter(wsHandler, messageHandler, roomHandler, authService)

	fmt.Printf("Chat server is running on port %s\n", cfg.Port)

	return http.ListenAndServe(cfg.Port, router)

}
