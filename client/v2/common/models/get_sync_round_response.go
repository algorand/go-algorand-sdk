package models

// GetSyncRoundResponse response containing the ledger's minimum sync round
type GetSyncRoundResponse struct {
	// Round the minimum sync round for the ledger.
	Round uint64 `json:"round"`
}
