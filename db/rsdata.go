package db

import (
	"DirectoryService/cfg"
	"DirectoryService/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

// DbCtx represents the database layer for the registry service.
type DbCtx struct {
	Pool *pgxpool.Pool
}

// NewDbCtx creates a new DbCtx instance.
func NewDbCtx(pool *pgxpool.Pool) *DbCtx {
	return &DbCtx{pool}
}

// ConnectDB creates a connection pool to the PostgreSQL database
func ConnectDB(config *cfg.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DB.User, config.DB.Password, config.DB.Host, config.DB.Port, config.DB.DBName,
		config.DB.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return pool, nil
}

// RegisterService inserts a new service into the database and returns the inserted service.
func (s *DbCtx) RegisterService(ctx context.Context, service models.Service) (
	*models.Service, error,
) {
	query := `
		INSERT INTO r1.services (service_id, name, description, owner_info, industry_category, client_rating)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING service_id, name, description, owner_info, industry_category, client_rating, created_at, updated_at
	`

	var newService models.Service
	err := s.Pool.QueryRow(
		ctx, query, service.ServiceID, service.Name, service.Description, service.OwnerInfo,
		service.IndustryCategory, service.ClientRating,
	).Scan(
		&newService.ServiceID, &newService.Name, &newService.Description, &newService.OwnerInfo,
		&newService.IndustryCategory, &newService.ClientRating, &newService.CreatedAt,
		&newService.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert service: %w", err)
	}

	return &newService, nil
}

// UpdateService updates the service details and returns the updated service.
func (s *DbCtx) UpdateService(ctx context.Context, service models.Service) (
	*models.Service, error,
) {
	query := `
		UPDATE r1.services
		SET name = $1, description = $2, owner_info = $3, industry_category = $4, client_rating = $5, updated_at = CURRENT_TIMESTAMP
		WHERE service_id = $6
		RETURNING service_id, name, description, owner_info, industry_category, client_rating, created_at, updated_at
	`

	var updatedService models.Service
	err := s.Pool.QueryRow(
		ctx, query, service.Name, service.Description, service.OwnerInfo, service.IndustryCategory,
		service.ClientRating, service.ServiceID,
	).Scan(
		&updatedService.ServiceID, &updatedService.Name, &updatedService.Description,
		&updatedService.OwnerInfo,
		&updatedService.IndustryCategory, &updatedService.ClientRating, &updatedService.CreatedAt,
		&updatedService.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return &updatedService, nil
}

// GetService retrieves a service by ID.
func (s *DbCtx) GetService(ctx context.Context, serviceID uuid.UUID) (*models.Service, error) {
	query := `
		SELECT service_id, name, description, owner_info, industry_category, client_rating
		FROM r1.services
		WHERE service_id = $1
	`

	var service models.Service

	err := s.Pool.QueryRow(ctx, query, serviceID).Scan(
		&service.ServiceID, &service.Name, &service.Description, &service.OwnerInfo,
		&service.IndustryCategory, &service.ClientRating,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve service: %w", err)
	}

	return &service, nil
}

// GetAllServices retrieves all services from the database.
func (s *DbCtx) GetAllServices(ctx context.Context) ([]models.Service, error) {
	query := `
		SELECT service_id, name, description, owner_info, industry_category, client_rating
		FROM r1.services
	`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var service models.Service
		if err := rows.Scan(
			&service.ServiceID, &service.Name, &service.Description, &service.OwnerInfo,
			&service.IndustryCategory, &service.ClientRating,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		services = append(services, service)
	}

	return services, nil
}

// DeleteService deletes a service by ID.
func (s *DbCtx) DeleteService(ctx context.Context, serviceID uuid.UUID) error {
	query := `
		DELETE FROM r1.services
		WHERE service_id = $1
	`

	_, err := s.Pool.Exec(ctx, query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

// ListServices retrieves all services from the database.
func (s *DbCtx) ListServices(ctx context.Context) ([]models.Service, error) {
	query := `
		SELECT service_id, name, description, owner_info, industry_category, client_rating
		FROM r1.services
	`

	rows, err := s.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var service models.Service
		if err := rows.Scan(
			&service.ServiceID, &service.Name, &service.Description, &service.OwnerInfo,
			&service.IndustryCategory, &service.ClientRating,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		services = append(services, service)
	}

	return services, nil
}

// Create a new ServiceInstance in the database.
func (s *DbCtx) CreateServiceInstance(
	ctx context.Context, instance models.ServiceInstance,
) (*models.ServiceInstance, error) {
	instance.InstanceID = uuid.New()
	instance.HealthStatus = models.HealthStatus("starting")
	instance.CreatedAt = time.Now().UTC()
	instance.LastChecked = time.Now().UTC()

	query := `
  INSERT INTO r1.service_instances (
   service_id, instance_id, version, host, port, url, api_spec, latitude, longitude, health_status, created_at, last_checked
  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
  RETURNING instance_id
 `

	var instanceID uuid.UUID
	err := s.Pool.QueryRow(
		ctx, query, instance.ServiceID, instance.InstanceID, instance.Version, instance.Host,
		instance.Port,
		instance.Url.String(), instance.ApiSpec, instance.Latitude, instance.Longitude,
		instance.HealthStatus,
		instance.CreatedAt, instance.LastChecked,
	).Scan(&instanceID)
	if err != nil {
		return nil, fmt.Errorf("failed to create service instance: %w", err)
	}

	instance.InstanceID = instanceID
	return &instance, nil
}

// GetServiceInstance retrieves a instanceID.
func (s *DbCtx) GetServiceInstance(
	ctx context.Context, instanceID uuid.UUID,
) (*models.ServiceInstance, error) {
	// create query to get service instance using InstanceID
	query := `
		SELECT service_id, instance_id, version, host, port, url, api_spec, latitude, longitude, health_status, created_at, last_checked
		FROM r1.service_instances
		WHERE instance_id = $1
	`

	var serviceInstance models.ServiceInstance

	err := s.Pool.QueryRow(ctx, query, instanceID).Scan(
		&serviceInstance.ServiceID, &serviceInstance.InstanceID, &serviceInstance.Version,
		&serviceInstance.Host,
		&serviceInstance.Port, &serviceInstance.Url, &serviceInstance.ApiSpec,
		&serviceInstance.Latitude, &serviceInstance.Longitude,
		&serviceInstance.HealthStatus, &serviceInstance.CreatedAt, &serviceInstance.LastChecked,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve service: %w", err)
	}

	return &serviceInstance, nil
}

// RemoveServiceInstance - copies pertinent columns from ServiceInstance to service_instance_history and inserts new row before it deletes the service instance by ID.
func (s *DbCtx) RemoveServiceInstance(ctx context.Context, instanceID uuid.UUID) error {

	// call GetServiceInstance to retrieve the service instance by ID
	serviceInstance, err := s.GetServiceInstance(ctx, instanceID)
	if err != nil {
		return fmt.Errorf("failed to get service instance: %w", err)
	}

	// User serviceIntance to insert into service_instance_history
	query := `
	  INSERT INTO r1.service_instance_history (		
		service_id, instance_id, version, url, metrics, started_at, stopped_at
	  ) VALUES ($1, $2, $3, $4, $5, $6, $7)
	  RETURNING history_id
`
	var historyID uuid.UUID
	// create a map to hold the Metrics
	metrics := make(map[string]interface{})
	metrics["health_status"] = serviceInstance.HealthStatus

	err = s.Pool.QueryRow(
		ctx, query, serviceInstance.ServiceID, serviceInstance.InstanceID,
		serviceInstance.Version, serviceInstance.Url, metrics, serviceInstance.CreatedAt,
		time.Now(),
	).Scan(&historyID)

	if err != nil {
		return fmt.Errorf("failed to copy service instance to history: %w", err)
	}

	// now delete the entry in service_instances
	deleteQuery := `
  DELETE FROM r1.service_instances
  WHERE instance_id = $1
 `

	_, err = s.Pool.Exec(ctx, deleteQuery, instanceID)
	if err != nil {
		return fmt.Errorf("failed to delete service instance: %w", err)
	}

	return nil
}
