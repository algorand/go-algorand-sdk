package models

// AccountDeltas exposes deltas for account based resources in a single round
type AccountDeltas struct {
	// Accounts array of Account updates for the round
	Accounts []AccountBalanceRecord `json:"accounts,omitempty"`

	// Apps array of App updates for the round.
	Apps []AppResourceRecord `json:"apps,omitempty"`

	// Assets array of Asset updates for the round.
	Assets []AssetResourceRecord `json:"assets,omitempty"`
}
