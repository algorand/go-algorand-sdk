package models

// Error an error response with optional data field.
type Error struct {
	// Data
	Data *map[string]interface{} `json:"data,omitempty"`

	// Message
	Message string `json:"message"`
}
