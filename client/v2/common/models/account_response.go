package models

// AccountResponse
type AccountResponse struct {
	// Account account information at a given round.
	// Definition:
	// data/basics/userBalance.go : AccountData
	Account Account `json:"account,omitempty"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round,omitempty"`
}
