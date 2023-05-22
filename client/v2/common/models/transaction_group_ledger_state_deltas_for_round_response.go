package models

// TransactionGroupLedgerStateDeltasForRoundResponse response containing all ledger
// state deltas for transaction groups, with their associated Ids, in a single
// round.
type TransactionGroupLedgerStateDeltasForRoundResponse struct {
	// Deltas
	Deltas []LedgerStateDeltaForTransactionGroup `json:"Deltas"`
}
