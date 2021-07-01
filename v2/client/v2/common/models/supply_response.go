package models

// SupplyResponse supply represents the current supply of MicroAlgos in the system.
type SupplyResponse struct {
	// Current_round round
	Current_round uint64 `json:"current_round"`

	// OnlineMoney onlineMoney
	OnlineMoney uint64 `json:"online-money"`

	// TotalMoney totalMoney
	TotalMoney uint64 `json:"total-money"`
}
