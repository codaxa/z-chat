package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthHandler responds to GET requests on the /health endpoint with a JSON status message.
// It returns a 405 error for non-GET methods and a 500 error if the response cannot be encoded.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "healthy"}); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
