package models

import (
	"github.com/google/uuid"
	"net/url"
	"time"
)

// ServiceInstance represents an instance of a service.
type ServiceInstance struct {
	ServiceID    uuid.UUID    `json:"service_id"`
	InstanceID   uuid.UUID    `json:"instance_id"`
	Version      string       `json:"version"`
	Host         string       `json:"host"`
	Port         int          `json:"port"`
	Url          url.URL      `json:"url"`
	Latitude     float64      `json:"latitude"`
	Longitude    float64      `json:"longitude"`
	HealthStatus HealthStatus `json:"health_status"`
	ApiSpec      string       `json:"api_spec"`
	CreatedAt    time.Time    `json:"created_at"`
	LastChecked  time.Time    `json:"last_checked"`
}
