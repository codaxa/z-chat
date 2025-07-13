// It implements the main entry point for the chat server application.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize the chat server
	fmt.Println("Starting chat server...")
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println("Running on port 8000")

	log.Fatal(http.ListenAndServe(":8000", nil))
}
