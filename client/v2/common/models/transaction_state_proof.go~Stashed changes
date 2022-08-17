package models

// TransactionStateProof fields for a state proof transaction.
// Definition:
// data/transactions/stateproof.go : StateProofTxnFields
type TransactionStateProof struct {
	// Message (spmsg)
	Message IndexerStateProofMessage `json:"message,omitempty"`

	// StateProof (sp) represents a state proof.
	// Definition:
	// crypto/stateproof/structs.go : StateProof
	StateProof StateProofFields `json:"state-proof,omitempty"`

	// StateProofType (sptype) Type of the state proof. Integer representing an entry
	// defined in protocol/stateproof.go
	StateProofType uint64 `json:"state-proof-type,omitempty"`
}
