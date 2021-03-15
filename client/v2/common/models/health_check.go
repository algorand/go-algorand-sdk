package models;

// HealthCheck a health check response.
type HealthCheck struct {
   // Data
  Data *map[string]interface{} `json:"data,omitempty"`

   // DbAvailable
  DbAvailable bool `json:"db-available,omitempty"`

   // Errors
  Errors []string `json:"errors,omitempty"`

   // IsMigrating
  IsMigrating bool `json:"is-migrating,omitempty"`

   // Message
  Message string `json:"message,omitempty"`

   // Round
  Round uint64 `json:"round,omitempty"`
}
