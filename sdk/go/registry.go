package _go

// ServiceRegistry defines the interface for managing service lifecycle, health checks, and statistics.
type ServiceRegistry interface {
	// RegisterService registers the service in the registry.
	RegisterService(
		serviceID, name, description, ownerInfo, industryCategory string, clientRating float64,
	) error

	// DeregisterService removes the service from the registry.
	DeregisterService(serviceID string) error

	// UpdateServiceHealth updates the health status of the service in the registry.
	UpdateServiceHealth(serviceID string, status HealthStatus) error

	// PerformHealthCheck performs a health check for the given service.
	PerformHealthCheck(serviceID string) (HealthStatus, error)

	// RetrieveStatistics retrieves the statistics for the given service in JSON format.
	RetrieveStatistics(serviceID string) (ServiceStatistics, error)
}

// HealthStatus represents the health status of a service.
type HealthStatus string

const (
	Up      HealthStatus = "up"
	Down    HealthStatus = "down"
	Unknown HealthStatus = "unknown"
)

// ServiceStatistics represents performance and usage statistics for a service.
type ServiceStatistics struct {
	ServiceID        string                 `json:"service_id"`
	TransactionCount int64                  `json:"transaction_count"`
	AvgResponseTime  float64                `json:"average_response_time"` // in milliseconds
	Details          map[string]interface{} `json:"details,omitempty"`     // Additional service-specific metrics
}

// RegistryClient is an implementation of ServiceRegistry for interacting with the directory service.
type RegistryClient struct {
	BaseURL string // Base URL of the directory service
}

// NewRegistryClient creates a new RegistryClient instance.
func NewRegistryClient(baseURL string) *RegistryClient {
	return &RegistryClient{BaseURL: baseURL}
}

// Implement the methods for RegistryClient...
