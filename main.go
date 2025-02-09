package main

import (
	"DirectoryService/cfg"
	"DirectoryService/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"DirectoryService/db"
)

func main() {
	// Set up database connection
	config, err := cfg.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load cfg: %v", err)
	}

	pool, err := db.ConnectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close() // Ensure connection pool is closed when the service shuts down

	fmt.Println("Database connection established")

	// Create DbCtx instance and inject into DbCtx
	server := handlers.NewServer(pool) // Inject DbCtx into DbCtx struct

	// Set up HTTP routes
	router := server.NewRouter()

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("DbCtx running on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("DbCtx error: %v", err)
	}
}
