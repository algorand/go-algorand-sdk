package models

// TransactionParametersResponse transactionParams contains the parameters that
// help a client construct a new transaction.
type TransactionParametersResponse struct {
	// ConsensusVersion consensusVersion indicates the consensus protocol version
	// as of LastRound.
	ConsensusVersion string `json:"consensus-version"`

	// Fee fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	Fee uint64 `json:"fee"`

	// GenesisHash genesisHash is the hash of the genesis block.
	GenesisHash []byte `json:"genesis-hash"`

	// GenesisId genesisID is an ID listed in the genesis block.
	GenesisId string `json:"genesis-id"`

	// LastRound lastRound indicates the last round seen
	LastRound uint64 `json:"last-round"`

	// MinFee the minimum transaction fee (not per byte) required for the
	// txn to validate for the current network protocol.
	MinFee uint64 `json:"min-fee"`
}
