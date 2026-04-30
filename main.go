package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	// Load .env file if present (ignored in production where env vars are set directly)
	if err := godotenv.Load(); err != nil {
		// Suppress the noisy warning when running in environments without a .env file
		// Note: only log this in DEBUG mode to keep output clean during normal dev
		if os.Getenv("DEBUG") == "true" {
			log.Println("No .env file found, using environment variables")
		}
	}

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log the full startup info including host so it's easier to identify
	// which interface we're binding to during local dev (useful when running
	// multiple services simultaneously)
	log.Printf("Starting ds2api %s on %s:%s", Version, cfg.Host, cfg.Port)

	// Also log the current DEBUG mode status so it's obvious at startup
	// whether verbose logging is active — saves confusion during local testing
	if os.Getenv("DEBUG") == "true" {
		log.Println("DEBUG mode enabled")
	}

	// Basic signal handling for graceful shutdown on Ctrl+C / SIGTERM
	// Addresses the TODO above — at least we log a clean exit message now
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutdown signal received, exiting...")
		os.Exit(0)
	}()

	server := NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}

	// TODO: consider adding a /healthz endpoint for use with Docker HEALTHCHECK
}
