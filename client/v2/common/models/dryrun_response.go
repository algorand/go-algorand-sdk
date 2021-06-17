package models

// DryrunResponse dryrunResponse contains per-txn debug information from a dryrun.
type DryrunResponse struct {
	// Error
	Error string `json:"error"`

	// ProtocolVersion protocol version is the protocol version Dryrun was operated
	// under.
	ProtocolVersion string `json:"protocol-version"`

	// Txns
	Txns []DryrunTxnResult `json:"txns"`
}
