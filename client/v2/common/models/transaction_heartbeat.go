package models

// TransactionHeartbeat fields for a heartbeat transaction.
// Definition:
// data/transactions/heartbeat.go : HeartbeatTxnFields
type TransactionHeartbeat struct {
	// HbAddress (hbad) HbAddress is the account this txn is proving onlineness for.
	HbAddress string `json:"hb-address"`

	// HbKeyDilution (hbkd) HbKeyDilution must match HbAddress account's current
	// KeyDilution.
	HbKeyDilution uint64 `json:"hb-key-dilution"`

	// HbProof (hbprf) HbProof is a signature using HeartbeatAddress's partkey, thereby
	// showing it is online.
	HbProof HbProofFields `json:"hb-proof"`

	// HbSeed (hbsd) HbSeed must be the block seed for the this transaction's
	// firstValid block.
	HbSeed []byte `json:"hb-seed"`

	// HbVoteId (hbvid) HbVoteID must match the HbAddress account's current VoteID.
	HbVoteId []byte `json:"hb-vote-id"`
}
