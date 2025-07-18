// It implements the main entry point for the chat server application.
package main

import (
	"fmt"
	"log"
	"net/http"
	"z-chat/internal/config"
	"z-chat/internal/handlers"
	"z-chat/internal/hub"
	route "z-chat/internal/transport/http"
)

// main initializes and starts the chat server, setting up HTTP endpoints and launching the chat hub.
func main() {
	// Initialize the chat server
	fmt.Println("Starting chat server...")

	chatHub := hub.NewHub()

	go chatHub.Run()
	wsHandler := handlers.NewWebSocketHandler(chatHub)
	router := route.NewRouter(wsHandler)

	cfg := config.New()

	fmt.Printf("Chat server is running on port %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, router))
}
