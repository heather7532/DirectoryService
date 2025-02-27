package handlers_test

import (
	"DirectoryService/cfg"
	"DirectoryService/db"
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"

	"DirectoryService/handlers"
	"DirectoryService/models"

	"github.com/stretchr/testify/assert"
)

func setupTestServer(t *testing.T) *handlers.Server {
	config, err := cfg.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load cfg: %v", err)
	}

	pool, err := db.ConnectDB(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	return handlers.NewServer(pool)
}

func TestRegisterServiceHandler(t *testing.T) {
	server := setupTestServer(t)

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
	server := setupTestServer(t)

	// Create Parent Service
	service := models.Service{
		Name:             "Parent Service",
		Description:      "Test service 1",
		OwnerInfo:        "Owner",
		IndustryCategory: "Category",
		ClientRating:     4.2,
	}

	svcBody, _ := json.Marshal(service)
	svcReq, err := http.NewRequest("POST", "/services", bytes.NewBuffer(svcBody))
	if err != nil {
		t.Fatal(err)
	}

	svcrr := httptest.NewRecorder()
	svcHandler := http.HandlerFunc(server.RegisterServiceHandler)
	svcHandler.ServeHTTP(svcrr, svcReq)

	assert.Equal(t, http.StatusOK, svcrr.Code)
	var responseService models.Service
	err = json.NewDecoder(svcrr.Body).Decode(&responseService)

	// Create Service Instance (child of Parent Service)
	if err != nil {
		t.Fatal(err)
	}

	serviceInstance := models.ServiceInstance{
		ServiceID:    responseService.ServiceID,
		Version:      "1.0.0",
		Host:         "localhost",
		Port:         8080,
		Url:          "http://localhost:8080",
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

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RegisterServiceInstanceHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseInstance models.ServiceInstance
	err = json.NewDecoder(rr.Body).Decode(&responseInstance)
	assert.NoError(t, err)
	assert.Equal(t, serviceInstance.ServiceID, responseInstance.ServiceID)
}

func TestRemoveServiceInstanceHandler(t *testing.T) {
	server := setupTestServer(t)

	// Create Parent Service
	service := models.Service{
		Name:             "Parent Service 2",
		Description:      "Test service 2",
		OwnerInfo:        "Owner",
		IndustryCategory: "Category",
		ClientRating:     4.4,
	}

	svcBody, _ := json.Marshal(service)
	svcReq, err := http.NewRequest("POST", "/services", bytes.NewBuffer(svcBody))
	if err != nil {
		t.Fatal(err)
	}

	svcrr := httptest.NewRecorder()
	svcHandler := http.HandlerFunc(server.RegisterServiceHandler)
	svcHandler.ServeHTTP(svcrr, svcReq)

	assert.Equal(t, http.StatusOK, svcrr.Code)
	var responseService models.Service
	err = json.NewDecoder(svcrr.Body).Decode(&responseService)

	// Create Service Instance (child of Parent Service)
	if err != nil {
		t.Fatal(err)
	}

	serviceInstance := models.ServiceInstance{
		ServiceID:    responseService.ServiceID,
		Version:      "1.1.0",
		Host:         "localhost",
		Port:         8080,
		Url:          "http://localhost:8080",
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

	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.RegisterServiceInstanceHandler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var responseInstance models.ServiceInstance
	err = json.NewDecoder(rr.Body).Decode(&responseInstance)
	assert.NoError(t, err)
	assert.Equal(t, serviceInstance.ServiceID, responseInstance.ServiceID)

	// Delete Service Instance
	deleteReq, err := http.NewRequest(
		"DELETE", "/service-instances/"+responseInstance.InstanceID.String(), nil,
	)
	if err != nil {
		t.Fatal(err)
	}

	router := mux.NewRouter()
	router.HandleFunc(
		"/service-instances/{id}", server.RemoveServiceInstanceHandler,
	).Methods("DELETE")

	deleteRR := httptest.NewRecorder()
	router.ServeHTTP(deleteRR, deleteReq) // Use router instead of calling handler directly

	assert.Equal(t, http.StatusNoContent, deleteRR.Code)
}
