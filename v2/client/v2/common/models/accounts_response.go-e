package models

// AccountsResponse
type AccountsResponse struct {
	// Accounts
	Accounts []Account `json:"accounts"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// NextToken used for pagination, when making another request provide this token
	// with the next parameter.
	NextToken string `json:"next-token,omitempty"`
}
