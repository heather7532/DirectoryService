package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"DirectoryService/cfg"
	"DirectoryService/db"
)

func main() {
	config := &cfg.Config{
		DB: cfg.DBConfig{
			User:     "user",
			Password: "password",
			Host:     "localhost",
			Port:     5432,
			DBName:   "directoryservice",
			SSLMode:  "disable",
		},
	}

	// Initialize the database connection
	pool, err := db.ConnectDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Create RSData instance
	rsData := db.NewRSData(pool)

	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Example: Insert a service
	service := db.Service{
		ServiceID:   "service123",
		ServiceName: "Example Service",
		Metadata:    map[string]string{"version": "1.0"},
	}

	if err := rsData.InsertService(ctx, service); err != nil {
		log.Fatalf("Failed to insert service: %v", err)
	}

	// Example: Retrieve all services
	services, err := rsData.ListServices(ctx)
	if err != nil {
		log.Fatalf("Failed to list services: %v", err)
	}

	for _, svc := range services {
		fmt.Printf("Service: %+v\n", svc)
	}
}
