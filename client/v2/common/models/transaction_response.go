package models

// TransactionResponse
type TransactionResponse struct {
	// CurrentRound round at which the results were computed.
	CurrentRound uint64 `json:"current-round,omitempty"`

	// Transaction contains all fields common to all transactions and serves as an
	// envelope to all transactions type.
	// Definition:
	// data/transactions/signedtxn.go : SignedTxn
	// data/transactions/transaction.go : Transaction
	Transaction Transaction `json:"transaction,omitempty"`
}
