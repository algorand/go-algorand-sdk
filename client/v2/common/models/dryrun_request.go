package models

import "github.com/algorand/go-algorand-sdk/types"

// DryrunRequest request data type for dryrun endpoint. Given the Transactions and
// simulated ledger state upload, run TEAL scripts and return debugging
// information.
type DryrunRequest struct {
	// Accounts
	Accounts []Account `json:"accounts,omitempty"`

	// Apps
	Apps []Application `json:"apps,omitempty"`

	// LatestTimestamp latestTimestamp is available to some TEAL scripts. Defaults to
	// the latest confirmed timestamp this algod is attached to.
	LatestTimestamp uint64 `json:"latest-timestamp,omitempty"`

	// ProtocolVersion protocolVersion specifies a specific version string to operate
	// under, otherwise whatever the current protocol of the network this algod is
	// running in.
	ProtocolVersion string `json:"protocol-version,omitempty"`

	// Round round is available to some TEAL scripts. Defaults to the current round on
	// the network this algod is attached to.
	Round uint64 `json:"round,omitempty"`

	// Sources
	Sources []DryrunSource `json:"sources,omitempty"`

	// Txns
	Txns []types.SignedTxn `json:"txns,omitempty"`
}
