package types

// BlockHash represents the hash of a block
type BlockHash Digest

// A BlockHeader represents the metadata and commitments to the state of a Block.
// The Algorand Ledger may be defined minimally as a cryptographically authenticated series of BlockHeader objects.
type BlockHeader struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Round Round `codec:"rnd"`

	// The hash of the previous block
	Branch BlockHash `codec:"prev"`

	// The hash of the previous block, using SHA-512
	Branch512 Sha512Digest `codec:"prev512"`

	// Sortition seed
	Seed Seed `codec:"seed"`

	// TxnCommitments authenticates the set of transactions appearing in the block.
	TxnCommitments

	// TimeStamp in seconds since epoch
	TimeStamp int64 `codec:"ts"`

	// Genesis ID to which this block belongs.
	GenesisID string `codec:"gen"`

	// Genesis hash to which this block belongs.
	GenesisHash Digest `codec:"gh"`

	// Proposer is the proposer of this block. Like the Seed, agreement adds
	// this after the block is assembled by the transaction pool, so that the same block can be prepared
	// for multiple participating accounts in the same node. Therefore, it can not be used
	// to influence block evaluation. Populated if proto.Payouts.Enabled
	Proposer Address `codec:"prp"`

	// FeesCollected is the sum of all fees paid by transactions in this
	// block. Populated if proto.Payouts.Enabled
	FeesCollected MicroAlgos `codec:"fc"`

	// Bonus is the bonus incentive to be paid for proposing this block.  It
	// begins as a consensus parameter value, and decays periodically.
	Bonus MicroAlgos `codec:"bi"`

	// ProposerPayout is the amount that is moved from the FeeSink to
	// the Proposer in this block.  It is basically the
	// bonus + the payouts percent of FeesCollected, but may be zero'd by
	// proposer ineligibility.
	ProposerPayout MicroAlgos `codec:"pp"`

	// Rewards.
	//
	// When a block is applied, some amount of rewards are accrued to
	// every account with AccountData.Status=/=NotParticipating.  The
	// amount is (thisBlock.RewardsLevel-prevBlock.RewardsLevel) of
	// MicroAlgos for every whole config.Protocol.RewardUnit of MicroAlgos in
	// that account's AccountData.MicroAlgos.
	//
	// Rewards are not compounded (i.e., not added to AccountData.MicroAlgos)
	// until some other transaction is executed on that account.
	//
	// Not compounding rewards allows us to precisely know how many algos
	// of rewards will be distributed without having to examine every
	// account to determine if it should get one more algo of rewards
	// because compounding formed another whole config.Protocol.RewardUnit
	// of algos.
	RewardsState

	// Consensus protocol versioning.
	//
	// Each block is associated with a version of the consensus protocol,
	// stored under UpgradeState.CurrentProtocol.  The protocol version
	// for a block can be determined without having to first decode the
	// block and its CurrentProtocol field, and this field is present for
	// convenience and explicitness.  Block.Valid() checks that this field
	// correctly matches the expected protocol version.
	//
	// Each block is associated with at most one active upgrade proposal
	// (a new version of the protocol).  An upgrade proposal can be made
	// by a block proposer, as long as no other upgrade proposal is active.
	// The upgrade proposal lasts for many rounds (UpgradeVoteRounds), and
	// in each round, that round's block proposer votes to support (or not)
	// the proposed upgrade.
	//
	// If enough votes are collected, the proposal is approved, and will
	// definitely take effect.  The proposal lingers for some number of
	// rounds to give clients a chance to notify users about an approved
	// upgrade, if the client doesn't support it, so the user has a chance
	// to download updated client software.
	//
	// Block proposers influence this upgrade machinery through two fields
	// in UpgradeVote: UpgradePropose, which proposes an upgrade to a new
	// protocol, and UpgradeApprove, which signals approval of the current
	// proposal.
	//
	// Once a block proposer determines its UpgradeVote, then UpdateState
	// is updated deterministically based on the previous UpdateState and
	// the new block's UpgradeVote.
	UpgradeState
	UpgradeVote

	// TxnCounter is the number of the next transaction that will be
	// committed after this block.  Genesis blocks can start at either
	// 0 or 1000, depending on a consensus parameter (AppForbidLowResources).
	TxnCounter uint64 `codec:"tc"`

	// StateProofTracking tracks the status of the state proofs, potentially
	// for multiple types of ASPs (Algorand's State Proofs).
	StateProofTracking map[StateProofType]StateProofTrackingData `codec:"spt"`

	// ParticipationUpdates contains the information needed to mark
	// certain accounts offline because their participation keys expired
	ParticipationUpdates
}

// TxnCommitments represents the commitments computed from the transactions in the block.
// It contains multiple commitments based on different algorithms and hash functions, to support different use cases.
type TxnCommitments struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	// Root of transaction Merkle tree using the SHA-512/256 hash function.
	// This commitment is computed based on the PaysetCommit type specified in the block's consensus protocol.
	NativeSha512_256Commitment Digest `codec:"txn"`

	// Root of transaction vector commitment Merkle tree using the SHA-256 hash function.
	Sha256Commitment Digest `codec:"txn256"`

	// Root of transaction vector commitment Merkle tree using the SHA-512 hash function.
	Sha512Commitment Sha512Digest `codec:"txn512"`
}

// ParticipationUpdates represents participation account data that
// needs to be checked/acted on by the network
type ParticipationUpdates struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// ExpiredParticipationAccounts contains a list of online accounts
	// that needs to be converted to offline since their
	// participation key expired.
	ExpiredParticipationAccounts []Address `codec:"partupdrmv"`

	// AbsentParticipationAccounts contains a list of online accounts that
	// needs to be converted to offline since they are not proposing.
	AbsentParticipationAccounts []Address `codec:"partupdabs"`
}

// RewardsState represents the global parameters controlling the rate
// at which accounts accrue rewards.
type RewardsState struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// The FeeSink accepts transaction fees. It can only spend to
	// the incentive pool.
	FeeSink Address `codec:"fees"`

	// The RewardsPool accepts periodic injections from the
	// FeeSink and continually redistributes them to addresses as
	// rewards.
	RewardsPool Address `codec:"rwd"`

	// RewardsLevel specifies how many rewards, in MicroAlgos,
	// have been distributed to each config.Protocol.RewardUnit
	// of MicroAlgos since genesis.
	RewardsLevel uint64 `codec:"earn"`

	// The number of new MicroAlgos added to the participation stake from rewards at the next round.
	RewardsRate uint64 `codec:"rate"`

	// The number of leftover MicroAlgos after the distribution of RewardsRate/rewardUnits
	// MicroAlgos for every reward unit in the next round.
	RewardsResidue uint64 `codec:"frac"`

	// The round at which the RewardsRate will be recalculated.
	RewardsRecalculationRound Round `codec:"rwcalr"`
}

// UpgradeVote represents the vote of the block proposer with
// respect to protocol upgrades.
type UpgradeVote struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// UpgradePropose indicates a proposed upgrade
	UpgradePropose string `codec:"upgradeprop"`

	// UpgradeDelay indicates the time between acceptance and execution
	UpgradeDelay Round `codec:"upgradedelay"`

	// UpgradeApprove indicates a yes vote for the current proposal
	UpgradeApprove bool `codec:"upgradeyes"`
}

// UpgradeState tracks the protocol upgrade state machine.  It is,
// strictly speaking, computable from the history of all UpgradeVotes
// but we keep it in the block for explicitness and convenience
// (instead of materializing it separately, like balances).
//
//msgp:ignore UpgradeState
type UpgradeState struct {
	CurrentProtocol string `codec:"proto"`
	NextProtocol    string `codec:"nextproto"`
	// NextProtocolApprovals is the number of approvals for the next protocol proposal. It is expressed in Round because it is a count of rounds.
	NextProtocolApprovals Round `codec:"nextyes"`
	// NextProtocolVoteBefore specify the last voting round for the next protocol proposal. If there is no voting for
	// an upgrade taking place, this would be zero.
	NextProtocolVoteBefore Round `codec:"nextbefore"`
	// NextProtocolSwitchOn specify the round number at which the next protocol would be adopted. If there is no upgrade taking place,
	// nor a wait for the next protocol, this would be zero.
	NextProtocolSwitchOn Round `codec:"nextswitch"`
}

// StateProofTrackingData tracks the status of state proofs.
type StateProofTrackingData struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// StateProofVotersCommitment is the root of a vector commitment containing the
	// online accounts that will help sign a state proof.  The
	// VC root, and the state proof, happen on blocks that
	// are a multiple of ConsensusParams.StateProofRounds.  For blocks
	// that are not a multiple of ConsensusParams.StateProofRounds,
	// this value is zero.
	StateProofVotersCommitment GenericDigest `codec:"v"`

	// StateProofOnlineTotalWeight is the total number of microalgos held by the online accounts
	// during the StateProof round (or zero, if the merkle root is zero - no commitment for StateProof voters).
	// This is intended for computing the threshold of votes to expect from StateProofVotersCommitment.
	StateProofOnlineTotalWeight MicroAlgos `codec:"t"`

	// StateProofNextRound is the next round for which we will accept
	// a StateProof transaction.
	StateProofNextRound Round `codec:"n"`
}

// A Block contains the Payset and metadata corresponding to a given Round.
type Block struct {
	BlockHeader
	Payset Payset `codec:"txns"`
}

// A Payset represents a common, unforgeable, consistent, ordered set of SignedTxn objects.
//
//msgp:allocbound Payset 100000
type Payset []SignedTxnInBlock

// SignedTxnInBlock is how a signed transaction is encoded in a block.
type SignedTxnInBlock struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	SignedTxnWithAD

	HasGenesisID   bool `codec:"hgi"`
	HasGenesisHash bool `codec:"hgh"`
}

// SignedTxnWithAD is a (decoded) SignedTxn with associated ApplyData
type SignedTxnWithAD struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	SignedTxn
	ApplyData
}

// ApplyData contains information about the transaction's execution.
type ApplyData struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Closing amount for transaction.
	ClosingAmount MicroAlgos `codec:"ca"`

	// Closing amount for asset transaction.
	AssetClosingAmount uint64 `codec:"aca"`

	// Rewards applied to the Sender, Receiver, and CloseRemainderTo accounts.
	SenderRewards   MicroAlgos `codec:"rs"`
	ReceiverRewards MicroAlgos `codec:"rr"`
	CloseRewards    MicroAlgos `codec:"rc"`
	EvalDelta       EvalDelta  `codec:"dt"`

	// If asa or app is being created, the id used. Else 0.
	// Names chosen to match naming the corresponding txn.
	// These are populated only when MaxInnerTransactions > 0 (TEAL 5)
	ConfigAsset   AssetIndex `codec:"caid"`
	ApplicationID AppIndex   `codec:"apid"`
}

// EvalDelta stores StateDeltas for an application's global key/value store, as
// well as StateDeltas for some number of accounts holding local state for that
// application
type EvalDelta struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	GlobalDelta StateDelta `codec:"gd"`

	// When decoding EvalDeltas, the integer key represents an offset into
	// [txn.Sender, txn.Accounts[0], txn.Accounts[1], ..., SharedAccts[0], SharedAccts[1], ...]
	LocalDeltas map[uint64]StateDelta `codec:"ld"`

	// If a program modifies the local of an account that is not the Sender, or
	// in txn.Accounts, it must be recorded here, so that the key in LocalDeltas
	// can refer to it.
	SharedAccts []Address `codec:"sa"`

	// The total allocbound calculation here accounts for the worse possible case of having bounds.MaxLogCalls individual log entries
	// with the length of all of them summing up to bounds.MaxEvalDeltaTotalLogSize which is the limit for the sum of individual log lengths
	Logs []string `codec:"lg"`

	InnerTxns []SignedTxnWithAD `codec:"itx"`
}

// StateDelta is a map from key/value store keys to ValueDeltas, indicating
// what should happen for that key
//
//msgp:allocbound StateDelta config.MaxStateDeltaKeys
type StateDelta map[string]ValueDelta

// ValueDelta links a DeltaAction with a value to be set
type ValueDelta struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Action DeltaAction `codec:"at"`
	Bytes  string      `codec:"bs"`
	Uint   uint64      `codec:"ui"`
}

// DeltaAction is an enum of actions that may be performed when applying a
// delta to a TEAL key/value store
type DeltaAction uint64

const (
	// SetBytesAction indicates that a TEAL byte slice should be stored at a key
	SetBytesAction DeltaAction = 1

	// SetUintAction indicates that a Uint should be stored at a key
	SetUintAction DeltaAction = 2

	// DeleteAction indicates that the value for a particular key should be deleted
	DeleteAction DeltaAction = 3
)
