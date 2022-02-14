package models

// ApplicationLocalStatesResponse defines model for ApplicationLocalStatesResponse.
type ApplicationLocalStatesResponse struct {
	AppsLocalStates []ApplicationLocalState `json:"apps-local-states"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token,omitempty"`
}
