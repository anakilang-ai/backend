// main.go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anakilang-ai/backend/routes"
)

func main() {
	// Define the server port
	port := ":8080"

	// Register route handlers
	http.HandleFunc("/", routes.URL)

	// Start the server and log the URL where it's running
	fmt.Printf("Server started at: http://localhost%s\n", port)

	// Start the HTTP server and handle any potential errors
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
