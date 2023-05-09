package models

// LedgerStateDeltaForTransactionGroup contains a ledger delta for a single
// transaction group
type LedgerStateDeltaForTransactionGroup struct {
	// Delta ledger StateDelta object
	Delta *map[string]interface{} `json:"delta"`

	// Ids
	Ids []string `json:"ids"`
}
