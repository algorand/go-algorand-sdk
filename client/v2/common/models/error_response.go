package models

// ErrorResponse response for errors
type ErrorResponse struct {
	// Data
	Data *map[string]interface{} `json:"data,omitempty"`

	// Message
	Message string `json:"message"`
}
