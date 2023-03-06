package models

// SimulateResponse result of a transaction group simulation.
type SimulateResponse struct {
	// LastRound the round immediately preceding this simulation. State changes through
	// this round were used to run this simulation.
	LastRound uint64 `json:"last-round"`

	// TxnGroups a result object for each transaction group that was simulated.
	TxnGroups []SimulateTransactionGroupResult `json:"txn-groups"`

	// Version the version of this response object.
	Version uint64 `json:"version"`

	// WouldSucceed indicates whether the simulated transactions would have succeeded
	// during an actual submission. If any transaction fails or is missing a signature,
	// this will be false.
	WouldSucceed bool `json:"would-succeed"`
}
