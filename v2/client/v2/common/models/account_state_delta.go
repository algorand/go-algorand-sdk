package models

// AccountStateDelta application state delta.
type AccountStateDelta struct {
	// Address
	Address string `json:"address"`

	// Delta application state delta.
	Delta []EvalDeltaKeyValue `json:"delta"`
}
