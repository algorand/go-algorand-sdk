package models

// AssetHoldingsResponse defines model for AssetHoldingsResponse.
type AssetHoldingsResponse struct {
	Assets []AssetHolding `json:"assets"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token,omitempty"`
}
