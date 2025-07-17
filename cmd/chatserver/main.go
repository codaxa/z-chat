// It implements the main entry point for the chat server application.
package main

import (
	"fmt"
	"log"
	"net/http"
	"z-chat/internal/handlers"
	"z-chat/internal/hub"
)

func main() {
	// Initialize the chat server
	fmt.Println("Starting chat server...")
	http.HandleFunc("/health", handlers.HealthHandler)

	chatHub := hub.NewHub()

	go chatHub.Run()
	wsHandler := handlers.NewWebSocketHandler(chatHub)
	http.HandleFunc("/ws", wsHandler.ServeWS)

	fmt.Println("Running on port 8000")

	log.Fatal(http.ListenAndServe(":8000", nil))

}
