package models

import "github.com/algorand/go-algorand-sdk/types"

// DryrunRequest request data type for dryrun endpoint. Given the Transactions and
// simulated ledger state upload, run TEAL scripts and return debugging
// information.
type DryrunRequest struct {
	// Accounts
	Accounts []Account `json:"accounts"`

	// Apps
	Apps []Application `json:"apps"`

	// LatestTimestamp latestTimestamp is available to some TEAL scripts. Defaults to
	// the latest confirmed timestamp this algod is attached to.
	LatestTimestamp uint64 `json:"latest-timestamp"`

	// ProtocolVersion protocolVersion specifies a specific version string to operate
	// under, otherwise whatever the current protocol of the network this algod is
	// running in.
	ProtocolVersion string `json:"protocol-version"`

	// Round round is available to some TEAL scripts. Defaults to the current round on
	// the network this algod is attached to.
	Round uint64 `json:"round"`

	// Sources
	Sources []DryrunSource `json:"sources"`

	// Txns
	Txns []types.SignedTxn `json:"txns"`
}
