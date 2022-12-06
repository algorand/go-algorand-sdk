package models

// NodeStatusResponse
type NodeStatusResponse struct {
	// Catchpoint the current catchpoint that is being caught up to
	Catchpoint string `json:"catchpoint,omitempty"`

	// CatchpointAcquiredBlocks the number of blocks that have already been obtained by
	// the node as part of the catchup
	CatchpointAcquiredBlocks uint64 `json:"catchpoint-acquired-blocks,omitempty"`

	// CatchpointProcessedAccounts the number of accounts from the current catchpoint
	// that have been processed so far as part of the catchup
	CatchpointProcessedAccounts uint64 `json:"catchpoint-processed-accounts,omitempty"`

	// CatchpointProcessedKvs the number of key-values (KVs) from the current
	// catchpoint that have been processed so far as part of the catchup
	CatchpointProcessedKvs uint64 `json:"catchpoint-processed-kvs,omitempty"`

	// CatchpointTotalAccounts the total number of accounts included in the current
	// catchpoint
	CatchpointTotalAccounts uint64 `json:"catchpoint-total-accounts,omitempty"`

	// CatchpointTotalBlocks the total number of blocks that are required to complete
	// the current catchpoint catchup
	CatchpointTotalBlocks uint64 `json:"catchpoint-total-blocks,omitempty"`

	// CatchpointTotalKvs the total number of key-values (KVs) included in the current
	// catchpoint
	CatchpointTotalKvs uint64 `json:"catchpoint-total-kvs,omitempty"`

	// CatchpointVerifiedAccounts the number of accounts from the current catchpoint
	// that have been verified so far as part of the catchup
	CatchpointVerifiedAccounts uint64 `json:"catchpoint-verified-accounts,omitempty"`

	// CatchpointVerifiedKvs the number of key-values (KVs) from the current catchpoint
	// that have been verified so far as part of the catchup
	CatchpointVerifiedKvs uint64 `json:"catchpoint-verified-kvs,omitempty"`

	// CatchupTime catchupTime in nanoseconds
	CatchupTime uint64 `json:"catchup-time"`

	// LastCatchpoint the last catchpoint seen by the node
	LastCatchpoint string `json:"last-catchpoint,omitempty"`

	// LastRound lastRound indicates the last round seen
	LastRound uint64 `json:"last-round"`

	// LastVersion lastVersion indicates the last consensus version supported
	LastVersion string `json:"last-version"`

	// NextVersion nextVersion of consensus protocol to use
	NextVersion string `json:"next-version"`

	// NextVersionRound nextVersionRound is the round at which the next consensus
	// version will apply
	NextVersionRound uint64 `json:"next-version-round"`

	// NextVersionSupported nextVersionSupported indicates whether the next consensus
	// version is supported by this node
	NextVersionSupported bool `json:"next-version-supported"`

	// StoppedAtUnsupportedRound stoppedAtUnsupportedRound indicates that the node does
	// not support the new rounds and has stopped making progress
	StoppedAtUnsupportedRound bool `json:"stopped-at-unsupported-round"`

	// TimeSinceLastRound timeSinceLastRound in nanoseconds
	TimeSinceLastRound uint64 `json:"time-since-last-round"`
}
