package main

import (
	"log"
	"os"

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

	server := NewServer(cfg)
	if err := server.Run(); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}

	// os.Exit(0) is redundant here since main() returning naturally exits with code 0
	// TODO: look into graceful shutdown via signal handling (SIGINT/SIGTERM) for cleaner Ctrl+C behaviour
}
