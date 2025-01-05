package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "github.com/opd-ai/srvmce"
)

func main() {
	// Initialize logger
	logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags)

	// Create server configuration
	cfg := server.Config{
		RootDir:       "./content", // Directory for content files
		AdminPath:     "/admin",    // Admin panel URL prefix
		AdminUser:     os.Getenv("ADMIN_USER"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
	}

	// Create new server instance
	srv, err := server.NewServer(cfg)
	if err != nil {
		logger.Fatalf("Failed to create server: %v", err)
	}

	// Create middleware instances
	logging := server.NewLoggingMiddleware(logger)
	recovery := server.NewRecoveryMiddleware(logger)

	// Create public routes
	srv.HandleFunc("/", handleHome, logging)
	srv.HandleFunc("/about", handleAbout, logging)

	// Create API routes with additional middleware
	api := srv.Group("/api", logging, recovery)
	api.HandleFunc("/health", handleHealth)
	api.HandleFunc("/status", handleStatus)

	// Start server
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatalf("Failed to create listener: %v", err)
	}

	// Graceful shutdown setup
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Serve(listener); err != nil {
			logger.Fatalf("Server failed: %v", err)
		}
	}()

	logger.Printf("Server started on %s", listener.Addr())

	// Wait for interrupt signal
	<-done
	logger.Print("Server is shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server shutdown failed: %v", err)
	}

	logger.Print("Server stopped gracefully")
}

// Handler functions
func handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Welcome to the Home Page</h1>"))
}

func handleAbout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>About Us</h1>"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"healthy"}`))
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"running","version":"1.0.0"}`))
}
