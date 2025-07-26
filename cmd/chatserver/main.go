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
	userRepository := postgres.NewUserRepository(dbConn)

	// Initialize services
	authService := services.NewAuthService(userRepository, cfg.JWTSecret, cfg.TokenDuration)

	messageRepo := postgres.NewMessageRepository(dbConn)

	// Initialize chat hub
	chatHub := hub.NewHub(messageRepo)
	go chatHub.Run()

	wsHandler := handlers.NewWebSocketHandler(chatHub)
	router := route.NewRouter(wsHandler, authService)

	fmt.Printf("Chat server is running on port %s\n", cfg.Port)

	return http.ListenAndServe(cfg.Port, router)
}
