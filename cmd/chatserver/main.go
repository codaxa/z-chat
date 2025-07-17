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

func main() {
	// Initialize the chat server
	fmt.Println("Starting chat server...")

	chatHub := hub.NewHub()

	go chatHub.Run()
	wsHandler := handlers.NewWebSocketHandler(chatHub)
	router := route.NewRouter(wsHandler)

	cfg := config.New()

	log.Fatal(http.ListenAndServe(cfg.Port, router))
}
