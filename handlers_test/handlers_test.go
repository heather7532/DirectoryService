package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"DirectoryService/handlers"
	"DirectoryService/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

func setupTestServer() *handlers.Server {
	// Mock the database connection pool
	dbPool, _ := pgxpool.New(context.Background(), "postgres://user:password@localhost:5432/testdb")
	return &handlers.Server{DB: dbPool}
}

func TestRegisterServiceHandler(t *testing.T) {
	server := setupTestServer()

	service := models.Service{
		Name:             "Test Service",
		Description:      "A test service",
		OwnerInfo:        "Owner",
		IndustryCategory: "Category",
		ClientRating:     4.5,
	}

	body, _ := json.Marshal(service)
	req, err := http.NewRequest("POST", "/services", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RegisterServiceHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseService models.Service
	err = json.NewDecoder(rr.Body).Decode(&responseService)
	assert.NoError(t, err)
	assert.Equal(t, service.Name, responseService.Name)
}

func TestRegisterServiceInstanceHandler(t *testing.T) {
	server := setupTestServer()

	parsedURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return
	}

	serviceInstance := models.ServiceInstance{
		ServiceID:    uuid.New(),
		Version:      "1.0.0",
		Host:         "localhost",
		Port:         8080,
		Url:          *parsedURL,
		ApiSpec:      "API Spec",
		Latitude:     37.7749,
		Longitude:    -122.4194,
		HealthStatus: "Healthy",
	}

	body, _ := json.Marshal(serviceInstance)
	req, err := http.NewRequest("POST", "/service-instances", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RegisterServiceInstanceHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseInstance models.ServiceInstance
	err = json.NewDecoder(rr.Body).Decode(&responseInstance)
	assert.NoError(t, err)
	assert.Equal(t, serviceInstance.ServiceID, responseInstance.ServiceID)
}
