package models

// Genesis defines a model for Genesis.
type Genesis struct {
	// Alloc
	Alloc []GenesisAllocation `json:"alloc"`

	// Comment
	Comment string `json:"comment,omitempty"`

	// Devmode
	Devmode bool `json:"devmode,omitempty"`

	// Fees
	Fees string `json:"fees"`

	// Id
	Id string `json:"id"`

	// Network
	Network string `json:"network"`

	// Proto
	Proto string `json:"proto"`

	// Rwd
	Rwd string `json:"rwd"`

	// Timestamp
	Timestamp uint64 `json:"timestamp"`
}
