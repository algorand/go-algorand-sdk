package models

// SupplyResponse supply represents the current supply of MicroAlgos in the system.
type SupplyResponse struct {
	// Current_round round
	Current_round uint64 `json:"current_round"`

	// OnlineMoney total stake held by accounts with status Online at current_round,
	// including those whose participation keys have expired but have not yet been
	// marked offline.
	OnlineMoney uint64 `json:"online-money"`

	// OnlineStake online stake used by agreement to vote for current_round, excluding
	// accounts whose participation keys have expired.
	OnlineStake uint64 `json:"online-stake"`

	// TotalMoney totalMoney
	TotalMoney uint64 `json:"total-money"`
}
