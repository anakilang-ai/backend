package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anakilang-ai/backend/routes"
)

func main() {
	// Setup route handlers
	http.HandleFunc("/", routes.URL)

	// Get the port from environment variables or use a default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)

	// Start the server
	fmt.Printf("Server started at: http://localhost%s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
