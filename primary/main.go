package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/anakilang-ai/backend/routes"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Logger instance
var log = logrus.New()

// Server configuration constants
const (
	defaultPort     = "8080"
	shutdownTimeout = 5 * time.Second
)

func init() {
	// Customize the logger if needed
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func main() {
	// Initialize router
	router := mux.NewRouter()
	router.HandleFunc("/", routes.URL).Methods(http.MethodGet)

	// Load port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	serverAddr := fmt.Sprintf(":%s", port)

	// Initialize HTTP server
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Start server in a separate goroutine
	go func() {
		log.Infof("Server starting at: http://localhost%s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// Create a context with timeout for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Info("Server exited gracefully")
}
