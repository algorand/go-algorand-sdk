package models

// DryrunResponse dryrunResponse contains per-txn debug information from a dryrun.
type DryrunResponse struct {
	// Error
	Error string `json:"error,omitempty"`

	// ProtocolVersion protocol version is the protocol version Dryrun was operated
	// under.
	ProtocolVersion string `json:"protocol-version,omitempty"`

	// Txns
	Txns []DryrunTxnResult `json:"txns,omitempty"`
}
