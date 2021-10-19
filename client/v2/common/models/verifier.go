package models

type Verifier struct {
	// Root
	Root [32]byte `json:"r,omitempty"`

	// HasValidRoot
	HasValidRoot bool `json:"vr,omitempty"`
}
