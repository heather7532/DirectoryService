package models

import (
	"github.com/google/uuid"
	"time"
)

// Service represents a service entity in the database.
type Service struct {
	ServiceID        uuid.UUID `json:"service_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	OwnerInfo        string    `json:"owner_info"`
	IndustryCategory string    `json:"industry_category"`
	ClientRating     float64   `json:"client_rating"`
	TransactionCount int64     `json:"transaction_count,omitempty"`
	AvgResponseTime  float64   `json:"average_response_time,omitempty"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
}
