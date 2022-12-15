package models

// LedgerStateDelta contains ledger updates.
type LedgerStateDelta struct {
	// Accts accountDeltas object
	Accts AccountDeltas `json:"accts,omitempty"`

	// KvMods array of KV Deltas
	KvMods []KvDelta `json:"kv-mods,omitempty"`

	// ModifiedApps list of modified Apps
	ModifiedApps []ModifiedApp `json:"modified-apps,omitempty"`

	// ModifiedAssets list of modified Assets
	ModifiedAssets []ModifiedAsset `json:"modified-assets,omitempty"`

	// PrevTimestamp previous block timestamp
	PrevTimestamp uint64 `json:"prev-timestamp,omitempty"`

	// StateProofNext next round for which we expect a state proof
	StateProofNext uint64 `json:"state-proof-next,omitempty"`

	// Totals account Totals
	Totals AccountTotals `json:"totals,omitempty"`

	// TxLeases list of transaction leases
	TxLeases []TxLease `json:"tx-leases,omitempty"`
}
