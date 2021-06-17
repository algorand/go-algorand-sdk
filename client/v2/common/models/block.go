package models

// Block block information.
// Definition:
// data/bookkeeping/block.go : Block
type Block struct {
	// GenesisHash (gh) hash to which this block belongs.
	GenesisHash []byte `json:"genesis-hash"`

	// GenesisId (gen) ID to which this block belongs.
	GenesisId string `json:"genesis-id"`

	// PreviousBlockHash (prev) Previous block hash.
	PreviousBlockHash []byte `json:"previous-block-hash"`

	// Rewards fields relating to rewards,
	Rewards BlockRewards `json:"rewards,omitempty"`

	// Round (rnd) Current round on which this block was appended to the chain.
	Round uint64 `json:"round"`

	// Seed (seed) Sortition seed.
	Seed []byte `json:"seed"`

	// Timestamp (ts) Block creation timestamp in seconds since eposh
	Timestamp uint64 `json:"timestamp"`

	// Transactions (txns) list of transactions corresponding to a given round.
	Transactions []Transaction `json:"transactions,omitempty"`

	// TransactionsRoot (txn) TransactionsRoot authenticates the set of transactions
	// appearing in the block. More specifically, it's the root of a merkle tree whose
	// leaves are the block's Txids, in lexicographic order. For the empty block, it's
	// 0. Note that the TxnRoot does not authenticate the signatures on the
	// transactions, only the transactions themselves. Two blocks with the same
	// transactions but in a different order and with different signatures will have
	// the same TxnRoot.
	TransactionsRoot []byte `json:"transactions-root"`

	// TxnCounter (tc) TxnCounter counts the number of transactions committed in the
	// ledger, from the time at which support for this feature was introduced.
	// Specifically, TxnCounter is the number of the next transaction that will be
	// committed after this block. It is 0 when no transactions have ever been
	// committed (since TxnCounter started being supported).
	TxnCounter uint64 `json:"txn-counter,omitempty"`

	// UpgradeState fields relating to a protocol upgrade.
	UpgradeState BlockUpgradeState `json:"upgrade-state,omitempty"`

	// UpgradeVote fields relating to voting for a protocol upgrade.
	UpgradeVote BlockUpgradeVote `json:"upgrade-vote,omitempty"`
}
