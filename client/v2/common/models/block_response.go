package models

// BlockResponse encoded block object.
type BlockResponse struct {
	// Block block header data.
	Block *map[string]interface{} `json:"block,omitempty"`

	// Cert optional certificate object. This is only included when the format is set
	// to message pack.
	Cert *map[string]interface{} `json:"cert,omitempty"`
}
