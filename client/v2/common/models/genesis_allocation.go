package models

// GenesisAllocation defines a model for GenesisAllocation.
type GenesisAllocation struct {
	// Addr
	Addr string `json:"addr"`

	// Comment
	Comment string `json:"comment"`

	// State
	State *map[string]interface{} `json:"state"`
}
