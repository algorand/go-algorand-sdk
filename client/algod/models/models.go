// Package models defines models used by an algod rest client//
// IF YOU MODIFY THIS FILE: IMPORTANT
// In practice, this is straight up copied from /v1/models/model.go. It is duplicated
// from internal model code, to maintain the internal/external client encapsulation.
// It does flatten some embedded structs, however.
// No client should depend on any package in v1.
package models

import "github.com/algorand/go-algorand-sdk/types"

// NodeStatus contains the information about a node status
// swagger:model NodeStatus
type NodeStatus struct {
	// LastRound indicates the last round seen
	//
	// required: true
	LastRound uint64 `json:"lastRound"`

	// LastVersion indicates the last consensus version supported
	//
	// required: true
	LastVersion string `json:"lastConsensusVersion"`

	// NextVersion of consensus protocol to use
	//
	// required: true
	NextVersion string `json:"nextConsensusVersion"`

	// NextVersionRound is the round at which the next consensus version will apply
	//
	// required: true
	NextVersionRound uint64 `json:"nextConsensusVersionRound"`

	// NextVersionSupported indicates whether the next consensus version is supported by this node
	//
	// required: true
	NextVersionSupported bool `json:"nextConsensusVersionSupported"`

	// TimeSinceLastRound in nanoseconds
	//
	// required: true
	TimeSinceLastRound int64 `json:"timeSinceLastRound"`

	// CatchupTime in nanoseconds
	//
	// required: true
	CatchupTime int64 `json:"catchupTime"`
}

// TransactionID Description
// swagger:model transactionID
type TransactionID struct {
	// TxId is the string encoding of the transaction hash
	//
	// required: true
	TxID string `json:"txId"`
}

// Account Description
// swagger:model Account
type Account struct {
	// Round indicates the round for which this information is relevant
	//
	// required: true
	Round uint64 `json:"round"`

	// Address indicates the account public key
	//
	// required: true
	Address string `json:"address"`

	// Amount indicates the total number of MicroAlgos in the account
	//
	// required: true
	Amount uint64 `json:"amount"`

	// PendingRewards specifies the amount of MicroAlgos of pending
	// rewards in this account.
	//
	// required: true
	PendingRewards uint64 `json:"pendingrewards"`

	// AmountWithoutPendingRewards specifies the amount of MicroAlgos in
	// the account, without the pending rewards.
	//
	// required: true
	AmountWithoutPendingRewards uint64 `json:"amountwithoutpendingrewards"`

	// Rewards indicates the total rewards of MicroAlgos the account has recieved
	//
	// required: true
	Rewards uint64 `json:"rewards"`

	// Status indicates the delegation status of the account's MicroAlgos
	// Offline - indicates that the associated account is delegated.
	// Online  - indicates that the associated account used as part of the delegation pool.
	// NotParticipating - indicates that the associated account is neither a delegator nor a delegate.
	//
	// required: true
	Status string `json:"status"`

	// AssetParams specifies the parameters of assets created by this account.
	//
	// required: false
	AssetParams map[uint64]AssetParams `json:"thisassettotal,omitempty"`

	// Assets specifies the holdings of assets by this account,
	// indexed by the asset ID.
	//
	// required: false
	Assets map[uint64]AssetHolding `json:"assets,omitempty"`
}

// AssetParams specifies the parameters for an asset.
// swagger:model AssetParams
type AssetParams struct {
	// Creator specifies the address that created this asset.
	// This is the address where the parameters for this asset
	// can be found, and also the address where unwanted asset
	// units can be sent in the worst case.
	//
	// required: true
	Creator string `json:"creator"`

	// Total specifies the total number of units of this asset.
	//
	// required: true
	Total uint64 `json:"total"`

	// DefaultFrozen specifies whether holdings in this asset
	// are frozen by default.
	//
	// required: false
	DefaultFrozen bool `json:"defaultfrozen"`

	// UnitName specifies the name of a unit of this asset,
	// as supplied by the creator.
	//
	// required: false
	UnitName string `json:"unitname"`

	// AssetName specifies the name of this asset,
	// as supplied by the creator.
	//
	// required: false
	AssetName string `json:"assetname"`

	// URL specifies a URL where more information about the asset can be
	// retrieved
	//
	// required: false
	URL string `json:"url"`

	// MetadataHash specifies a commitment to some unspecified asset
	// metadata. The format of this metadata is up to the application.
	//
	// required: false
	// swagger:strfmt byte
	MetadataHash []byte `json:"metadatahash"`

	// ManagerAddr specifies the address used to manage the keys of this
	// asset and to destroy it.
	//
	// required: false
	ManagerAddr string `json:"managerkey"`

	// ReserveAddr specifies the address holding reserve (non-minted)
	// units of this asset.
	//
	// required: false
	ReserveAddr string `json:"reserveaddr"`

	// FreezeAddr specifies the address used to freeze holdings of
	// this asset.  If empty, freezing is not permitted.
	//
	// required: false
	FreezeAddr string `json:"freezeaddr"`

	// ClawbackAddr specifies the address used to clawback holdings of
	// this asset.  If empty, clawback is not permitted.
	//
	// required: false
	ClawbackAddr string `json:"clawbackaddr"`
}

// AssetHolding specifies the holdings of a particular asset.
// swagger:model AssetHolding
type AssetHolding struct {
	// Creator specifies the address that created this asset.
	// This is the address where the parameters for this asset
	// can be found, and also the address where unwanted asset
	// units can be sent in the worst case.
	//
	// required: true
	Creator string `json:"creator"`

	// Amount specifies the number of units held.
	//
	// required: true
	Amount uint64 `json:"amount"`

	// Frozen specifies whether this holding is frozen.
	//
	// required: false
	Frozen bool `json:"frozen"`
}

// Bytes is a byte array
// swagger:strfmt binary
type Bytes = []byte // note that we need to make this its own object to get the strfmt annotation to work properly. Otherwise swagger generates []uint8 instead of type binary
// ^ one day we should probably fork swagger, to avoid this workaround.

// Transaction contains all fields common to all transactions and serves as an envelope to all transactions
// type
// swagger:model Transaction
type Transaction struct {
	// Type is the transaction type
	//
	// required: true
	Type types.TxType `json:"type"`

	// TxID is the transaction ID
	//
	// required: true
	TxID string `json:"tx"`

	// From is the sender's address
	//
	// required: true
	From string `json:"from"`

	// Fee is the transaction fee
	//
	// required: true
	Fee uint64 `json:"fee"`

	// FirstRound indicates the first valid round for this transaction
	//
	// required: true
	FirstRound uint64 `json:"first-round"`

	// LastRound indicates the last valid round for this transaction
	//
	// required: true
	LastRound uint64 `json:"last-round"`

	// Note is a free form data
	//
	// required: false
	Note Bytes `json:"noteb64,omitempty"`

	// ConfirmedRound indicates the block number this transaction appeared in
	//
	// required: false
	ConfirmedRound uint64 `json:"round,omitempty"`

	// PoolError indicates the transaction was evicted from this node's transaction
	// pool (if non-empty).  A non-empty PoolError does not guarantee that the
	// transaction will never be committed; other nodes may not have evicted the
	// transaction and may attempt to commit it in the future.
	//
	// required: false
	PoolError string `json:"poolerror,omitempty"`

	// This is a list of all supported transactions.
	// To add another one, create a struct with XXXTransactionType and embed it here.
	// To prevent extraneous fields, all must have the "omitempty" tag.
	Payment *PaymentTransactionType `json:"payment,omitempty"`

	// FromRewards is the amount of pending rewards applied to the From
	// account as part of this transaction.
	//
	// required: false
	FromRewards uint64 `json:"fromrewards,omitempty"`

	// Genesis ID
	//
	// required: true
	GenesisID string `json:"genesisID"`

	// Genesis hash
	//
	// required: true
	GenesisHash Bytes `json:"genesishashb64"`
}

// PaymentTransactionType contains the additional fields for a payment Transaction
// swagger:model PaymentTransactionType
type PaymentTransactionType struct {
	// To is the receiver's address
	//
	// required: true
	To string `json:"to"`

	// CloseRemainderTo is the address the sender closed to
	//
	// required: false
	CloseRemainderTo string `json:"close,omitempty"`

	// CloseAmount is the amount sent to CloseRemainderTo, for committed transaction
	//
	// required: false
	CloseAmount uint64 `json:"closeamount,omitempty"`

	// Amount is the amount of MicroAlgos intended to be transferred
	//
	// required: true
	Amount uint64 `json:"amount"`

	// ToRewards is the amount of pending rewards applied to the To account
	// as part of this transaction.
	//
	// required: false
	ToRewards uint64 `json:"torewards,omitempty"`

	// CloseRewards is the amount of pending rewards applied to the CloseRemainderTo
	// account as part of this transaction.
	//
	// required: false
	CloseRewards uint64 `json:"closerewards,omitempty"`
}

// TransactionList contains a list of transactions
// swagger:model TransactionList
type TransactionList struct {
	// TransactionList is a list of transactions
	//
	// required: true
	Transactions []Transaction `json:"transactions,omitempty"`
}

// TransactionFee contains the suggested fee
// swagger:model TransactionFee
type TransactionFee struct {
	// Fee is transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	//
	// required: true
	Fee uint64 `json:"fee"`
}

// TransactionParams contains the parameters that help a client construct
// a new transaction.
// swagger:model TransactionParams
type TransactionParams struct {
	// Fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	//
	// required: true
	Fee uint64 `json:"fee"`

	// Genesis ID
	//
	// required: true
	GenesisID string `json:"genesisID"`

	// Genesis hash
	//
	// required: true
	GenesisHash Bytes `json:"genesishashb64"`

	// LastRound indicates the last round seen
	//
	// required: true
	LastRound uint64 `json:"lastRound"`

	// ConsensusVersion indicates the consensus protocol version
	// as of LastRound.
	//
	// required: true
	ConsensusVersion string `json:"consensusVersion"`
}

// Block contains a block information
// swagger:model Block
type Block struct {
	// Hash is the current block hash
	//
	// required: true
	Hash string `json:"hash"`

	// PreviousBlockHash is the previous block hash
	//
	// required: true
	PreviousBlockHash string `json:"previousBlockHash"`

	// Seed is the sortition seed
	//
	// required: true
	Seed string `json:"seed"`

	// Proposer is the address of this block proposer
	//
	// required: true
	Proposer string `json:"proposer"`

	// Round is the current round on which this block was appended to the chain
	//
	// required: true
	Round uint64 `json:"round"`

	// Period is the period on which the block was confirmed
	//
	// required: true
	Period uint64 `json:"period"`

	// TransactionsRoot authenticates the set of transactions appearing in the block.
	// More specifically, it's the root of a merkle tree whose leaves are the block's Txids, in lexicographic order.
	// For the empty block, it's 0.
	// Note that the TxnRoot does not authenticate the signatures on the transactions, only the transactions themselves.
	// Two blocks with the same transactions but in a different order and with different signatures will have the same TxnRoot.
	//
	// required: true
	TransactionsRoot string `json:"txnRoot"`

	// RewardsLevel specifies how many rewards, in MicroAlgos,
	// have been distributed to each config.Protocol.RewardUnit
	// of MicroAlgos since genesis.
	RewardsLevel uint64 `json:"reward"`

	// The number of new MicroAlgos added to the participation stake from rewards at the next round.
	RewardsRate uint64 `json:"rate"`

	// The number of leftover MicroAlgos after the distribution of RewardsRate/rewardUnits
	// MicroAlgos for every reward unit in the next round.
	RewardsResidue uint64 `json:"frac"`

	// Transactions is the list of transactions in this block
	Transactions TransactionList `json:"txns"`

	// TimeStamp in seconds since epoch
	//
	// required: true
	Timestamp int64 `json:"timestamp"`

	UpgradeState
	UpgradeVote
}

// UpgradeState contains the information about a current state of an upgrade
// swagger:model UpgradeState
type UpgradeState struct {
	// CurrentProtocol is a string that represents the current protocol
	//
	// required: true
	CurrentProtocol string `json:"currentProtocol"`

	// NextProtocol is a string that represents the next proposed protocol
	//
	// required: true
	NextProtocol string `json:"nextProtocol"`

	// NextProtocolApprovals is the number of blocks which approved the protocol upgrade
	//
	// required: true
	NextProtocolApprovals uint64 `json:"nextProtocolApprovals"`

	// NextProtocolVoteBefore is the deadline round for this protocol upgrade (No votes will be consider after this round)
	//
	// required: true
	NextProtocolVoteBefore uint64 `json:"nextProtocolVoteBefore"`

	// NextProtocolSwitchOn is the round on which the protocol upgrade will take effect
	//
	// required: true
	NextProtocolSwitchOn uint64 `json:"nextProtocolSwitchOn"`
}

// UpgradeVote represents the vote of the block proposer with respect to protocol upgrades.
// swagger:model UpgradeVote
type UpgradeVote struct {
	// UpgradePropose indicates a proposed upgrade
	//
	// required: true
	UpgradePropose string `json:"upgradePropose"`

	// UpgradeApprove indicates a yes vote for the current proposal
	//
	// required: true
	UpgradeApprove bool `json:"upgradeApprove"`
}

// Supply represents the current supply of MicroAlgos in the system
// swagger:model Supply
type Supply struct {
	// Round
	//
	// required: true
	Round uint64 `json:"round"`

	// TotalMoney
	//
	// required: true
	TotalMoney uint64 `json:"totalMoney"`

	// OnlineMoney
	//
	// required: true
	OnlineMoney uint64 `json:"onlineMoney"`
}

// PendingTransactions represents a potentially truncated list of transactions currently in the
// node's transaction pool.
// swagger:model PendingTransactions
type PendingTransactions struct {
	// TruncatedTxns
	// required: true
	TruncatedTxns TransactionList `json:"truncatedTxns"`
	// TotalTxns
	// required: true
	TotalTxns uint64 `json:"totalTxns"`
}

// Version contains the current algod version.
//
// Note that we annotate this as a model so that legacy clients
// can directly import a swagger generated Version model.
// swagger:model Version
type Version struct {
	// required: true
	Versions []string `json:"versions"`
	// required: true
	GenesisID string `json:"genesis_id"`
	// required: true
	GenesisHash Bytes `json:"genesis_hash_b64"`
}

// VersionsResponse is the response to 'GET /versions'
//
// swagger:response VersionsResponse
type VersionsResponse struct {
	// in: body
	Body Version
}
