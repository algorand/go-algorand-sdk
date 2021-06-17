package models

import "github.com/algorand/go-algorand-sdk/types"

// PendingTransactionsResponse a potentially truncated list of transactions
// currently in the node's transaction pool. You can compute whether or not the
// list is truncated if the number of elements in the **top-transactions** array is
// fewer than **total-transactions**.
type PendingTransactionsResponse struct {
	// TopTransactions an array of signed transaction objects.
	TopTransactions []types.SignedTxn `json:"top-transactions"`

	// TotalTransactions total number of transactions in the pool.
	TotalTransactions uint64 `json:"total-transactions"`
}
