package db

import (
	"DirectoryService/cfg"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *RSData {
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

	return NewRSData(pool)
}

func TestInsertService(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	service := Service{
		ServiceID:        uuid.New(),
		Name:             "Test Service",
		Description:      "A test service",
		OwnerInfo:        "Test Owner",
		IndustryCategory: "Test Category",
		ClientRating:     4.5,
	}

	err := rs.InsertService(context.Background(), service)
	assert.NoError(t, err, "InsertService should not return an error")
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

func TestUpdateService(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	service := Service{
		ServiceID:        uuid.New(),
		Name:             "Updated Service",
		Description:      "An updated test service",
		OwnerInfo:        "Updated Owner",
		IndustryCategory: "Updated Category",
		ClientRating:     4.8,
	}

	err := rs.UpdateService(context.Background(), service)
	assert.NoError(t, err, "UpdateService should not return an error")
}

func TestListServices(t *testing.T) {
	rs := setupTestDB(t)
	defer rs.Pool.Close()

	services, err := rs.ListServices(context.Background())
	assert.NoError(t, err, "ListServices should not return an error")
	assert.NotEmpty(t, services, "ListServices should return a list of services")
}
