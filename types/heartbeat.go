package types

// HeartbeatTxnFields captures the fields used for an account to prove it is
// online (really, it proves that an entity with the account's part keys is able
// to submit transactions, so it should be able to propose/vote.)
type HeartbeatTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// HeartbeatAddress is the account this txn is proving onlineness for.
	HbAddress Address `codec:"hbad"`

	// HbProof is a signature using HeartbeatAddress's partkey, thereby showing it is online.
	HbProof HeartbeatProof `codec:"hbprf"`

	// The final three fields are included to allow early, concurrent check of
	// the HbProof.

	// HbSeed must be the block seed for this transaction's firstValid
	// block. It is the message that must be signed with HbAddress's part key.
	HbSeed Seed `codec:"hbsd"`

	// HbVoteID must match the HbAddress account's current VoteID.
	HbVoteID OneTimeSignatureVerifier `codec:"hbvid"`

	// HbKeyDilution must match HbAddress account's current KeyDilution.
	HbKeyDilution uint64 `codec:"hbkd"`
}
