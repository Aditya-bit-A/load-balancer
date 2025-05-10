package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Response represents the JSON structure to be returned
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	// Get server ID from environment variable, default to "unknown" if not set
	serverID := os.Getenv("SERVER_ID")
	if serverID == "" {
		serverID = "unknown"
	}

	// Define the handler for the /home endpoint
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		// Only handle GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Create the response
		response := Response{
			Message: "Hello from Server: " + serverID,
			Status:  "successful",
		}

		// Set content type to JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Encode the response as JSON and write to response writer
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding JSON: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	})

	// Define the handler for the /heartbeat endpoint
	http.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		// Return empty response with 200 status code
		w.WriteHeader(http.StatusOK)
	})

	// Start the server on port 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server %s starting on port %s...", serverID, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}