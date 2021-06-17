package models

import "github.com/algorand/go-algorand-sdk/types"

// BlockResponse encoded block object.
type BlockResponse struct {
	// Block block header data.
	Block types.Block `json:"block"`

	// Cert optional certificate object. This is only included when the format is set
	// to message pack.
	Cert *map[string]interface{} `json:"cert,omitempty"`
}
