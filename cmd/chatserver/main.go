// It implements the main entry point for the chat server application.
package main

import (
	"fmt"
	"log"
	"net/http"
	"z-chat/internal/config"
	"z-chat/internal/handlers"
	"z-chat/internal/hub"
	"z-chat/internal/services"
	"z-chat/internal/storage/postgres"
	route "z-chat/internal/transport/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/net/context"
)

// main initializes and starts the chat server, setting up HTTP endpoints and launching the chat hub.
func main() {
	// Initialize configuration
	cfg := config.New()

	// Initialize database connection
	connConfig, err := pgxpool.ParseConfig(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Failed to parse database configuration: %v", err)
	}

	db, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := postgres.NewUserRepo(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.TokenDuration)

	messageRepo := postgres.NewMessageRepository(db)

	// Initialize chat hub
	chatHub := hub.NewHub(messageRepo)
	go chatHub.Run()

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(chatHub)

	// Set up router with all required dependencies
	router := route.NewRouter(wsHandler, authService)

	fmt.Printf("Chat server is running on port %s\n", cfg.Port)

	// Use a regular error check instead of log.Fatal to allow deferred functions to run
	err = http.ListenAndServe(cfg.Port, router)
	if err != nil {
		log.Printf("Server error: %v", err)
		log.Printf("HTTP server failed: %v", err)
		// Remove os.Exit(1) to allow deferred functions to run
	}
}
