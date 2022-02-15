package models

// HealthCheck a health check response.
type HealthCheck struct {
	// Data
	Data *map[string]interface{} `json:"data,omitempty"`

	// DbAvailable
	DbAvailable bool `json:"db-available"`

	// Errors
	Errors []string `json:"errors,omitempty"`

	// IsMigrating
	IsMigrating bool `json:"is-migrating"`

	// Message
	Message string `json:"message"`

	// Round
	Round uint64 `json:"round"`

	// Version current version.
	Version string `json:"version"`
}
