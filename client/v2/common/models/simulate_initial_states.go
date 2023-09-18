package models

// SimulateInitialStates initial states of resources that were accessed during
// simulation.
type SimulateInitialStates struct {
	// AppInitialStates the initial states of accessed application before simulation.
	// The order of this array is arbitrary.
	AppInitialStates []ApplicationInitialStates `json:"app-initial-states,omitempty"`
}
