package models

// TransactionPayment fields for a payment transaction.
// Definition:
// data/transactions/payment.go : PaymentTxnFields
type TransactionPayment struct {
	// Amount (amt) number of MicroAlgos intended to be transferred.
	Amount uint64 `json:"amount"`

	// CloseAmount number of MicroAlgos that were sent to the close-remainder-to
	// address when closing the sender account.
	CloseAmount uint64 `json:"close-amount,omitempty"`

	// CloseRemainderTo (close) when set, indicates that the sending account should be
	// closed and all remaining funds be transferred to this address.
	CloseRemainderTo string `json:"close-remainder-to,omitempty"`

	// Receiver (rcv) receiver's address.
	Receiver string `json:"receiver"`
}
