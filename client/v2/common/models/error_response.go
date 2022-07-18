package models

// ErrorResponse an error response with optional data field.
type ErrorResponse struct {
	// Data
	Data *map[string]interface{} `json:"data,omitempty"`

	// Message
	Message string `json:"message"`
}
