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

	// HbSeed must be the block seed for the block before this transaction's
	// firstValid. It is supplied in the transaction so that Proof can be
	// checked at submit time without a ledger lookup, and must be checked at
	// evaluation time for equality with the actual blockseed.
	HbSeed Seed `codec:"hbsd"`
}
