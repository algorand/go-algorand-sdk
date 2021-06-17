package models

// ApplicationResponse
type ApplicationResponse struct {
	// Application application index and its parameters
	Application Application `json:"application,omitempty"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`
}
