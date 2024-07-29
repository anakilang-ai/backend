package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/anakilang-ai/backend/routes"
)

func main() {
	// Set up a new server mux and attach the URL route
	mux := http.NewServeMux()
	mux.HandleFunc("/", routes.URL)

	// Define the server with a specific address and handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Run the server in a separate goroutine to allow for graceful shutdown
	go func() {
		fmt.Printf("Server started at: http://localhost%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Set up channel to listen for OS interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Block until an interrupt signal is received
	<-stop

	// Create a context with a timeout to allow for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
