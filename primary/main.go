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

func init() {
	// Customize the logger if needed
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func main() {
	// Create a new router
	r := mux.NewRouter()
	r.HandleFunc("/", routes.URL).Methods(http.MethodGet)

	// Determine the port from environment variables or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serverAddr := fmt.Sprintf(":%s", port)

	// Create a server
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		fmt.Printf("Server started at: http://localhost%s\n", serverAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig

	// Create a context with a timeout for shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	log.Info("Server exited gracefully")
}
