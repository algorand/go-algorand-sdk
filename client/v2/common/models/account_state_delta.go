package models

// AccountStateDelta application state delta.
type AccountStateDelta struct {
	// Address
	Address string `json:"address,omitempty"`

	// Delta application state delta.
	Delta []EvalDeltaKeyValue `json:"delta,omitempty"`
}
