package models

// AccountTotals total Algos in the system grouped by account status
type AccountTotals struct {
	// NotParticipating amount of stake in non-participating accounts
	NotParticipating uint64 `json:"not-participating"`

	// Offline amount of stake in offline accounts
	Offline uint64 `json:"offline"`

	// Online amount of stake in online accounts
	Online uint64 `json:"online"`

	// RewardsLevel total number of algos received per reward unit since genesis
	RewardsLevel uint64 `json:"rewards-level"`
}
