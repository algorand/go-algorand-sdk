package models

// PostTransactionsResponse transaction ID of the submission.
type PostTransactionsResponse struct {
	// Txid encoding of the transaction hash.
	Txid string `json:"txId"`
}
