package db

import (
	"DirectoryService/cfg"
	"DirectoryService/models"
	"context"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *DbCtx {
	config, err := cfg.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	pool, err := ConnectDB(config)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Clear the table before each test
	_, err = pool.Exec(
		context.Background(),
		`DELETE FROM r1.services WHERE name LIKE 'test_%'`,
	)
	if err != nil {
		t.Fatalf("Failed to clear table: %v", err)
	}

	return NewDbCtx(pool)
}

func TestInsertService(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	service := models.Service{
		ServiceID:        uuid.New(),
		Name:             "Test Service",
		Description:      "A test service",
		OwnerInfo:        "Test Owner",
		IndustryCategory: "Test Category",
		ClientRating:     4.5,
	}

	insertedService, err := rs.RegisterService(context.Background(), service)
	assert.NoError(t, err, "RegisterService should not return an error")
	assert.NotNil(t, insertedService, "RegisterService should return the inserted service")
	assert.Equal(t, service.ServiceID, insertedService.ServiceID, "ServiceID should match")
}

func TestUpdateService(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	service := models.Service{
		ServiceID:        uuid.New(),
		Name:             "Test Service",
		Description:      "A test service",
		OwnerInfo:        "Test Owner",
		IndustryCategory: "Test Category",
		ClientRating:     4.5,
	}

	insertedService, err := rs.RegisterService(context.Background(), service)
	assert.NoError(t, err, "RegisterService should not return an error")

	insertedService.Name = "Updated Service"
	insertedService.Description = "An updated test service"
	insertedService.OwnerInfo = "Updated Owner"
	insertedService.IndustryCategory = "Updated Category"
	insertedService.ClientRating = 4.8

	updatedService, err := rs.UpdateService(context.Background(), *insertedService)
	assert.NoError(t, err, "UpdateService should not return an error")
	assert.NotNil(t, updatedService, "UpdateService should return the updated service")
	assert.Equal(t, insertedService.ServiceID, updatedService.ServiceID, "ServiceID should match")
	assert.Equal(t, "Updated Service", updatedService.Name, "Name should match")
	assert.Equal(
		t, "An updated test service", updatedService.Description, "Description should match",
	)
	assert.Equal(t, "Updated Owner", updatedService.OwnerInfo, "OwnerInfo should match")
	assert.Equal(
		t, "Updated Category", updatedService.IndustryCategory, "IndustryCategory should match",
	)
	assert.Equal(t, 4.8, updatedService.ClientRating, "ClientRating should match")
}

// test to GetAllServices
func TestGetAllServices(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	services, err := rs.GetAllServices(context.Background())
	assert.NoError(t, err, "GetAllServices should not return an error")
	assert.NotEmpty(t, services, "GetAllServices should return a list of services")
}

func TestGetService(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	services, err := rs.GetAllServices(context.Background())
	assert.NoError(t, err, "GetAllServices should not return an error")
	assert.NotEmpty(t, services, "GetAllServices should return a list of services")

	serviceID := services[0].ServiceID

	service, err := rs.GetService(context.Background(), serviceID)
	assert.NoError(t, err, "GetService should not return an error")
	assert.NotNil(t, service, "GetService should return a service")
	assert.Equal(t, serviceID, service.ServiceID, "ServiceID should match")
}

func TestDeleteService(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	serviceID := uuid.New()

	err := rs.DeleteService(context.Background(), serviceID)
	assert.NoError(t, err, "DeleteService should not return an error")
}

func TestListServices(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	services, err := rs.ListServices(context.Background())
	assert.NoError(t, err, "ListServices should not return an error")
	assert.NotEmpty(t, services, "ListServices should return a list of services")
}

func TestCreateServiceInstance(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	service := models.Service{
		ServiceID:        uuid.New(),
		Name:             "Test Service",
		Description:      "A test service",
		OwnerInfo:        "Test Owner",
		IndustryCategory: "Test Category",
		ClientRating:     4.5,
	}

	insertedService, err := rs.RegisterService(context.Background(), service)
	assert.NoError(t, err, "RegisterService should not return an error")
	assert.NotNil(t, insertedService, "RegisterService should return the inserted service")

	// Populate a ServiceInstance struct
	serviceInstance := models.ServiceInstance{
		ServiceID: insertedService.ServiceID,
		Host:      "localhost",
		Port:      8080,
		Version:   "1.0.0",
		Url: url.URL{
			Scheme: "http",
			Host:   "localhost:8080",
		},
		Latitude:  0.0,
		Longitude: 0.0,
		ApiSpec:   "test",
	}

	newInstance, err := rs.CreateServiceInstance(context.Background(), serviceInstance)
	assert.NoError(t, err, "CreateServiceInstance should not return an error")
	assert.NotNil(t, newInstance, "CreateServiceInstance should return an instance")
	assert.Equal(t, insertedService.ServiceID, newInstance.ServiceID, "ServiceID should match")
}
