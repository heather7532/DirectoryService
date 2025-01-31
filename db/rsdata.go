package db

import (
	"DirectoryService/cfg"
	"context"
	"fmt"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RSData represents the data access layer for the service registry.
type RSData struct {
	Pool *pgxpool.Pool
}

// NewRSData creates a new instance of RSData with the given connection pool.
func NewRSData(pool *pgxpool.Pool) *RSData {
	return &RSData{Pool: pool}
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

// InsertService inserts a new service into the database.
func (rs *RSData) InsertService(ctx context.Context, service Service) error {
	query := `
		INSERT INTO r1.services (service_id, name, description, owner_info, industry_category, client_rating)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := rs.Pool.Exec(
		ctx, query, service.ServiceID, service.Name, service.Description, service.OwnerInfo,
		service.IndustryCategory, service.ClientRating,
	)
	if err != nil {
		return fmt.Errorf("failed to insert service: %w", err)
	}

	return nil
}

// GetService retrieves a service by ID.
func (rs *RSData) GetService(ctx context.Context, serviceID uuid.UUID) (*Service, error) {
	query := `
		SELECT service_id, name, description, owner_info, industry_category, client_rating
		FROM r1.services
		WHERE service_id = $1
	`

	var service Service

	err := rs.Pool.QueryRow(ctx, query, serviceID).Scan(
		&service.ServiceID, &service.Name, &service.Description, &service.OwnerInfo,
		&service.IndustryCategory, &service.ClientRating,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve service: %w", err)
	}

	return &service, nil
}

// GetAllServices retrieves all services from the database.
func (rs *RSData) GetAllServices(ctx context.Context) ([]Service, error) {
	query := `
		SELECT service_id, name, description, owner_info, industry_category, client_rating
		FROM r1.services
	`

	rows, err := rs.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
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
func (rs *RSData) DeleteService(ctx context.Context, serviceID uuid.UUID) error {
	query := `
		DELETE FROM r1.services
		WHERE service_id = $1
	`

	_, err := rs.Pool.Exec(ctx, query, serviceID)
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	return nil
}

// UpdateService updates the service details.
func (rs *RSData) UpdateService(ctx context.Context, service Service) error {
	query := `
		UPDATE r1.services
		SET name = $1, description = $2, owner_info = $3, industry_category = $4, client_rating = $5
		WHERE service_id = $6
	`

	_, err := rs.Pool.Exec(
		ctx, query, service.Name, service.Description, service.OwnerInfo, service.IndustryCategory,
		service.ClientRating, service.ServiceID,
	)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	return nil
}

// ListServices retrieves all services from the database.
func (rs *RSData) ListServices(ctx context.Context) ([]Service, error) {
	query := `
		SELECT service_id, name, description, owner_info, industry_category, client_rating
		FROM r1.services
	`

	rows, err := rs.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query services: %w", err)
	}
	defer rows.Close()

	var services []Service
	for rows.Next() {
		var service Service
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
