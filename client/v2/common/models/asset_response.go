package models

// AssetResponse
type AssetResponse struct {
	// Asset specifies both the unique identifier and the parameters for an asset
	Asset Asset `json:"asset,omitempty"`

	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round,omitempty"`
}
