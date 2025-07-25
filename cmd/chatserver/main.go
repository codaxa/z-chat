// It implements the main entry point for the chat server application.
package main

import (
	"fmt"
	"log"
	"net/http"
	"z-chat/internal/config"
	"z-chat/internal/handlers"
	"z-chat/internal/hub"
	"z-chat/internal/storage/postgres"
	route "z-chat/internal/transport/http"

	"github.com/joho/godotenv"
)

// main initializes and starts the chat server, setting up HTTP endpoints and launching the chat hub.
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize the chat server
	fmt.Println("Starting chat server...")

	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	// Create database connection
	dbConn, err := postgres.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Create repository instances
	messageRepo := postgres.NewMessageRepository(dbConn)

	// Initialize hub with repository
	chatHub := hub.NewHub(messageRepo)
	go chatHub.Run()

	// Initialize handlers
	wsHandler := handlers.NewWebSocketHandler(chatHub)
	router := route.NewRouter(wsHandler)

	cfg := config.New()
	fmt.Printf("Chat server is running on port %s\n", cfg.Port)

	return http.ListenAndServe(cfg.Port, router)
}
