package models

// TransactionParametersResponse transactionParams contains the parameters that
// help a client construct a new transaction.
type TransactionParametersResponse struct {
	// ConsensusVersion consensusVersion indicates the consensus protocol version
	// as of LastRound.
	ConsensusVersion string `json:"consensus-version,omitempty"`

	// Fee fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but a group of N atomic transactions must
	// still have a fee of at least N*MinTxnFee for the current network protocol.
	Fee uint64 `json:"fee,omitempty"`

	// GenesisHash genesisHash is the hash of the genesis block.
	GenesisHash []byte `json:"genesis-hash,omitempty"`

	// GenesisId genesisID is an ID listed in the genesis block.
	GenesisId string `json:"genesis-id,omitempty"`

	// LastRound lastRound indicates the last round seen
	LastRound uint64 `json:"last-round,omitempty"`

	// MinFee the minimum transaction fee (not per byte) required for the
	// txn to validate for the current network protocol.
	MinFee uint64 `json:"min-fee,omitempty"`
}
