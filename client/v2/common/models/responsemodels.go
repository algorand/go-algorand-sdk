package models

import "github.com/algorand/go-algorand-sdk/types"

// Account information at a given round.
// Definition:
// data/basics/userBalance.go : AccountData
type Account struct {
	// Address the account public key
	Address string `json:"address,omitempty"`

	// Amount (algo) total number of MicroAlgos in the account
	Amount uint64 `json:"amount,omitempty"`

	// AmountWithoutPendingRewards specifies the amount of MicroAlgos in the account,
	// without the pending rewards.
	AmountWithoutPendingRewards uint64 `json:"amount-without-pending-rewards,omitempty"`

	// AppsLocalState (appl) applications local data stored in this account.
	// Note the raw object uses `map[int] -> AppLocalState` for this type.
	AppsLocalState []ApplicationLocalStates `json:"apps-local-state,omitempty"`

	// AppsTotalSchema (tsch) stores the sum of all of the local schemas and global
	// schemas in this account.
	// Note: the raw account uses `StateSchema` for this type.
	AppsTotalSchema ApplicationStateSchema `json:"apps-total-schema,omitempty"`

	// Assets (asset) assets held by this account.
	// Note the raw object uses `map[int] -> AssetHolding` for this type.
	Assets []AssetHolding `json:"assets,omitempty"`

	// AuthAddr (spend) the address against which signing should be checked. If empty,
	// the address of the current account is used. This field can be updated in any
	// transaction by setting the RekeyTo field.
	AuthAddr string `json:"auth-addr,omitempty"`

	// CreatedApps (appp) parameters of applications created by this account including
	// app global data.
	// Note: the raw account uses `map[int] -> AppParams` for this type.
	CreatedApps []Application `json:"created-apps,omitempty"`

	// CreatedAssets (apar) parameters of assets created by this account.
	// Note: the raw account uses `map[int] -> Asset` for this type.
	CreatedAssets []Asset `json:"created-assets,omitempty"`

	// Participation accountParticipation describes the parameters used by this account
	// in consensus protocol.
	Participation AccountParticipation `json:"participation,omitempty"`

	// PendingRewards amount of MicroAlgos of pending rewards in this account.
	PendingRewards uint64 `json:"pending-rewards,omitempty"`

	// RewardBase (ebase) used as part of the rewards computation. Only applicable to
	// accounts which are participating.
	RewardBase uint64 `json:"reward-base,omitempty"`

	// Rewards (ern) total rewards of MicroAlgos the account has received, including
	// pending rewards.
	Rewards uint64 `json:"rewards,omitempty"`

	// Round the round for which this information is relevant.
	Round uint64 `json:"round,omitempty"`

	// SigType indicates what type of signature is used by this account, must be one
	// of:
	// * sig
	// * msig
	// * lsig
	SigType string `json:"sig-type,omitempty"`

	// Status (onl) delegation status of the account's MicroAlgos
	// * Offline - indicates that the associated account is delegated.
	// * Online - indicates that the associated account used as part of the delegation
	// pool.
	// * NotParticipating - indicates that the associated account is neither a
	// delegator nor a delegate.
	Status string `json:"status,omitempty"`
}

// AccountParticipation describes the parameters used by this account in consensus protocol.
type AccountParticipation struct {

	// \[sel\] Selection public key (if any) currently registered for this round.
	SelectionParticipationKey []byte `json:"selection-participation-key,omitempty"`

	// \[voteFst\] First round for which this participation is valid.
	VoteFirstValid uint64 `json:"vote-first-valid,omitempty"`

	// \[voteKD\] Number of subkeys in each batch of participation keys.
	VoteKeyDilution uint64 `json:"vote-key-dilution,omitempty"`

	// \[voteLst\] Last round for which this participation is valid.
	VoteLastValid uint64 `json:"vote-last-valid,omitempty"`

	// \[vote\] root participation public key (if any) currently registered for this round.
	VoteParticipationKey []byte `json:"vote-participation-key,omitempty"`
}

// Asset specifies both the unique identifier and the parameters for an asset
type Asset struct {

	// unique asset identifier
	Index uint64 `json:"index"`

	// AssetParams specifies the parameters for an asset.
	//
	// \[apar\] when part of an AssetConfig transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params AssetParams `json:"params"`
}

// AssetHolding describes an asset held by an account.
type AssetHolding struct {

	// \[a\] number of units held.
	Amount uint64 `json:"amount"`

	// Asset ID of the holding.
	AssetId uint64 `json:"asset-id"`

	// Address that created this asset. This is the address where the parameters for this asset can be found, and also the address where unwanted asset units can be sent in the worst case.
	Creator string `json:"creator"`

	// \[f\] whether or not the holding is frozen.
	IsFrozen bool `json:"is-frozen,omitempty"`
}

// AssetParams specifies the parameters for an asset.
type AssetParams struct {

	// \[c\] Address of account used to clawback holdings of this asset.  If empty, clawback is not permitted.
	Clawback string `json:"clawback,omitempty"`

	// The address that created this asset. This is the address where the parameters for this asset can be found, and also the address where unwanted asset units can be sent in the worst case.
	Creator string `json:"creator"`

	// \[dc\] The number of digits to use after the decimal point when displaying this asset. If 0, the asset is not divisible. If 1, the base unit of the asset is in tenths. If 2, the base unit of the asset is in hundredths, and so on. This value must be between 0 and 19 (inclusive).
	Decimals uint64 `json:"decimals"`

	// \[df\] Whether holdings of this asset are frozen by default.
	DefaultFrozen bool `json:"default-frozen,omitempty"`

	// \[f\] Address of account used to freeze holdings of this asset.  If empty, freezing is not permitted.
	Freeze string `json:"freeze,omitempty"`

	// \[m\] Address of account used to manage the keys of this asset and to destroy it.
	Manager string `json:"manager,omitempty"`

	// \[am\] A commitment to some unspecified asset metadata. The format of this metadata is up to the application.
	MetadataHash []byte `json:"metadata-hash,omitempty"`

	// \[an\] Name of this asset, as supplied by the creator.
	Name string `json:"name,omitempty"`

	// \[r\] Address of account holding reserve (non-minted) units of this asset.
	Reserve string `json:"reserve,omitempty"`

	// \[t\] The total number of units of this asset.
	Total uint64 `json:"total"`

	// \[un\] Name of a unit of this asset, as supplied by the creator.
	UnitName string `json:"unit-name,omitempty"`

	// \[au\] URL where more information about the asset can be retrieved.
	Url string `json:"url,omitempty"`
}

type RawBlockJson struct {
	Block string
}
type RawBlockMsgpack struct {
	Block string `json:"url,omitempty"`
}

// Block defines model for Block.
type Block struct {
	Cert string `json:"cert"`

	// \[gh\] hash to which this block belongs.
	GenesisHash []byte `json:"genesis-hash"`

	// \[gen\] ID to which this block belongs.
	GenesisId string `json:"genesis-id"`

	// Current block hash
	Hash []byte `json:"hash"`

	// Period on which the block was confirmed.
	Period uint64 `json:"period"`

	// \[prev\] Previous block hash.
	PreviousBlockHash []byte `json:"previous-block-hash"`

	// Address that proposed this block.
	Proposer string `json:"proposer"`

	// Fields relating to rewards,
	Rewards BlockRewards `json:"rewards,omitempty"`

	// \[rnd\] Current round on which this block was appended to the chain.
	Round uint64 `json:"round"`

	// \[seed\] Sortition seed.
	Seed []byte `json:"seed"`

	// \[ts\] Block creation timestamp in seconds since eposh
	Timestamp uint64 `json:"timestamp"`

	// \[txns\] list of transactions corresponding to a given round.
	Transactions []Transaction `json:"transactions,omitempty"`

	// \[txn\] TransactionsRoot authenticates the set of transactions appearing in the block. More specifically, it's the root of a merkle tree whose leaves are the block's Txids, in lexicographic order. For the empty block, it's 0. Note that the TxnRoot does not authenticate the signatures on the transactions, only the transactions themselves. Two blocks with the same transactions but in a different order and with different signatures will have the same TxnRoot.
	TransactionsRoot []byte `json:"transactions-root"`

	// \[tc\] TxnCounter counts the number of transactions committed in the ledger, from the time at which support for this feature was introduced.
	//
	// Specifically, TxnCounter is the number of the next transaction that will be committed after this block.  It is 0 when no transactions have ever been committed (since TxnCounter started being supported).
	TxnCounter uint64 `json:"txn-counter,omitempty"`

	// Fields relating to a protocol upgrade.
	UpgradeState BlockUpgradeState `json:"upgrade-state,omitempty"`

	// Fields relating to voting for a protocol upgrade.
	UpgradeVote BlockUpgradeVote `json:"upgrade-vote,omitempty"`
}

// BlockRewards defines model for BlockRewards.
type BlockRewards struct {

	// \[fees\] accepts transaction fees, it can only spend to the incentive pool.
	FeeSink string `json:"fee-sink"`

	// \[rwcalr\] number of leftover MicroAlgos after the distribution of rewards-rate MicroAlgos for every reward unit in the next round.
	RewardsCalculationRound uint64 `json:"rewards-calculation-round"`

	// \[earn\] How many rewards, in MicroAlgos, have been distributed to each RewardUnit of MicroAlgos since genesis.
	RewardsLevel uint64 `json:"rewards-level"`

	// \[rwd\] accepts periodic injections from the fee-sink and continually redistributes them as rewards.
	RewardsPool string `json:"rewards-pool"`

	// \[rate\] Number of new MicroAlgos added to the participation stake from rewards at the next round.
	RewardsRate uint64 `json:"rewards-rate"`

	// \[frac\] Number of leftover MicroAlgos after the distribution of RewardsRate/rewardUnits MicroAlgos for every reward unit in the next round.
	RewardsResidue uint64 `json:"rewards-residue"`
}

// BlockUpgradeState defines model for BlockUpgradeState.
type BlockUpgradeState struct {

	// \[proto\] The current protocol version.
	CurrentProtocol string `json:"current-protocol"`

	// \[nextproto\] The next proposed protocol version.
	NextProtocol string `json:"next-protocol,omitempty"`

	// \[nextyes\] Number of blocks which approved the protocol upgrade.
	NextProtocolApprovals uint64 `json:"next-protocol-approvals,omitempty"`

	// \[nextswitch\] Round on which the protocol upgrade will take effect.
	NextProtocolSwitchOn uint64 `json:"next-protocol-switch-on,omitempty"`

	// \[nextbefore\] Deadline round for this protocol upgrade (No votes will be consider after this round).
	NextProtocolVoteBefore uint64 `json:"next-protocol-vote-before,omitempty"`
}

// BlockUpgradeVote defines model for BlockUpgradeVote.
type BlockUpgradeVote struct {

	// \[upgradeyes\] Indicates a yes vote for the current proposal.
	UpgradeApprove bool `json:"upgrade-approve,omitempty"`

	// \[upgradedelay\] Indicates the time between acceptance and execution.
	UpgradeDelay uint64 `json:"upgrade-delay,omitempty"`

	// \[upgradeprop\] Indicates a proposed upgrade.
	UpgradePropose string `json:"upgrade-propose,omitempty"`
}

// MiniAssetHolding defines model for MiniAssetHolding.
type MiniAssetHolding struct {
	Address  string `json:"address"`
	Amount   uint64 `json:"amount"`
	IsFrozen bool   `json:"is-frozen"`
}

// NodeStatus contains the information about a node's 'status.
type NodeStatus struct {

	// CatchupTime in nanoseconds
	CatchupTime uint64 `json:"catchup-time"`

	// HasSyncedSinceStartup indicates whether a round has completed since startup
	HasSyncedSinceStartup bool `json:"has-synced-since-startup"`

	// LastRound indicates the last round seen
	LastRound uint64 `json:"last-round"`

	// LastVersion indicates the last consensus version supported
	LastVersion string `json:"last-version,omitempty"`

	// NextVersion of consensus protocol to use
	NextVersion string `json:"next-version,omitempty"`

	// NextVersionRound is the round at which the next consensus version will apply
	NextVersionRound uint64 `json:"next-version-round,omitempty"`

	// NextVersionSupported indicates whether the next consensus version is supported by this node
	NextVersionSupported bool `json:"next-version-supported,omitempty"`

	// StoppedAtUnsupportedRound indicates that the node does not support the new rounds and has stopped making progress
	StoppedAtUnsupportedRound bool `json:"stopped-at-unsupported-round"`

	// TimeSinceLastRound in nanoseconds
	TimeSinceLastRound uint64 `json:"time-since-last-round"`
}

// Supply defines model for Supply.
type Supply struct {

	// OnlineMoney
	OnlineMoney uint64 `json:"online-money"`

	// Round
	Round uint64 `json:"current_round"`

	// TotalMoney
	TotalMoney uint64 `json:"total-money"`
}

// Transaction defines model for Transaction.
type Transaction struct {

	// Fields for asset allocation, re-configuration, and destruction.
	//
	//
	// A zero value for asset-id indicates asset creation.
	// A zero value for the params indicates asset destruction.
	//
	// Definition:
	// data/transactions/asset.go : AssetConfigTxnFields
	AssetConfigTransaction TransactionAssetConfig `json:"asset-config-transaction,omitempty"`

	// Fields for an asset freeze transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetFreezeTxnFields
	AssetFreezeTransaction TransactionAssetFreeze `json:"asset-freeze-transaction,omitempty"`

	// Fields for an asset transfer transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetTransferTxnFields
	AssetTransferTransaction TransactionAssetTransfer `json:"asset-transfer-transaction,omitempty"`

	// \[rc\] rewards applied to close-remainder-to account.
	CloseRewards uint64 `json:"close-rewards,omitempty"`

	// \[ca\] closing amount for transaction.
	ClosingAmount uint64 `json:"closing-amount,omitempty"`

	// Round when the transaction was confirmed.
	ConfirmedRound uint64 `json:"confirmed-round,omitempty"`

	// Specifies an asset index (ID) if an asset was created with this transaction.
	CreatedAssetIndex uint64 `json:"created-asset-index,omitempty"`

	// \[fee\] Transaction fee.
	Fee uint64 `json:"fee"`

	// \[fv\] First valid round for this transaction.
	FirstValid uint64 `json:"first-valid"`

	// \[gh\] Hash of genesis block.
	GenesisHash []byte `json:"genesis-hash,omitempty"`

	// \[gen\] genesis block ID.
	GenesisId string `json:"genesis-id,omitempty"`

	// \[grp\] Base64 encoded byte array of a sha512/256 digest. When present indicates that this transaction is part of a transaction group and the value is the sha512/256 hash of the transactions in that group.
	Group []byte `json:"group,omitempty"`

	// Transaction ID
	Id string `json:"id"`

	// Fields for a keyreg transaction.
	//
	// Definition:
	// data/transactions/keyreg.go : KeyregTxnFields
	KeyregTransaction TransactionKeyreg `json:"keyreg-transaction,omitempty"`

	// \[lv\] Last valid round for this transaction.
	LastValid uint64 `json:"last-valid"`

	// \[lx\] Base64 encoded 32-byte array. Lease enforces mutual exclusion of transactions.  If this field is nonzero, then once the transaction is confirmed, it acquires the lease identified by the (Sender, Lease) pair of the transaction until the LastValid round passes.  While this transaction possesses the lease, no other transaction specifying this lease can be confirmed.
	Lease []byte `json:"lease,omitempty"`

	// \[note\] Free form data.
	Note []byte `json:"note,omitempty"`

	// Fields for a payment transaction.
	//
	// Definition:
	// data/transactions/payment.go : PaymentTxnFields
	PaymentTransaction TransactionPayment `json:"payment-transaction,omitempty"`

	// \[rr\] rewards applied to receiver account.
	ReceiverRewards uint64 `json:"receiver-rewards,omitempty"`

	// Time when the block this transaction is in was confirmed.
	RoundTime uint64 `json:"round-time,omitempty"`

	// \[snd\] Sender's address.
	Sender string `json:"sender"`

	// \[rs\] rewards applied to sender account.
	SenderRewards uint64 `json:"sender-rewards,omitempty"`

	// Validation signature associated with some data. Only one of the signatures should be provided.
	Signature TransactionSignature `json:"signature"`

	// \[type\] Indicates what type of transaction this is. Different types have different fields.
	//
	// Valid types, and where their fields are stored:
	// * \[pay\] payment-transaction
	// * \[keyreg\] keyreg-transaction
	// * \[acfg\] asset-config-transaction
	// * \[axfer\] asset-transfer-transaction
	// * \[afrz\] asset-freeze-transaction
	Type string `json:"tx-type"`

	// \[rekey\] when included in a valid transaction,
	// the accounts auth addr will be updated with this value
	// and future signatures must be signed with the key represented by this address.
	RekeyTo string `json:"rekey-to,omitempty"`
}

// TransactionAssetConfig defines model for TransactionAssetConfig.
type TransactionAssetConfig struct {

	// \[xaid\] ID of the asset being configured or empty if creating.
	AssetId uint64 `json:"asset-id,omitempty"`

	// AssetParams specifies the parameters for an asset.
	//
	// \[apar\] when part of an AssetConfig transaction.
	//
	// Definition:
	// data/transactions/asset.go : AssetParams
	Params AssetParams `json:"params,omitempty"`
}

// TransactionAssetFreeze defines model for TransactionAssetFreeze.
type TransactionAssetFreeze struct {

	// \[fadd\] Address of the account whose asset is being frozen or thawed.
	Address string `json:"address"`

	// \[faid\] ID of the asset being frozen or thawed.
	AssetId uint64 `json:"asset-id"`

	// \[afrz\] The new freeze status.
	NewFreezeStatus bool `json:"new-freeze-status"`
}

// TransactionAssetTransfer defines model for TransactionAssetTransfer.
type TransactionAssetTransfer struct {

	// \[aamt\] Amount of asset to transfer. A zero amount transferred to self allocates that asset in the account's Assets map.
	Amount uint64 `json:"amount"`

	// \[xaid\] ID of the asset being transferred.
	AssetId uint64 `json:"asset-id"`

	// \[aclose\] Indicates that the asset should be removed from the account's Assets map, and specifies where the remaining asset holdings should be transferred.  It's always valid to transfer remaining asset holdings to the creator account.
	CloseTo string `json:"close-to,omitempty"`

	// \[arcv\] Recipient address of the transfer.
	Receiver string `json:"receiver"`

	// \[asnd\] The effective sender during a clawback transactions. If this is not a zero value, the real transaction sender must be the Clawback address from the AssetParams.
	Sender string `json:"sender,omitempty"`

	CloseAmount uint64 `json:"close-amount,omitempty"`
}

// TransactionKeyreg defines model for TransactionKeyreg.
type TransactionKeyreg struct {

	// \[nonpart\] Mark the account as participating or non-participating.
	NonParticipation bool `json:"non-participation,omitempty"`

	// \[selkey\] Public key used with the Verified Random Function (VRF) result during committee selection.
	SelectionParticipationKey []byte `json:"selection-participation-key,omitempty"`

	// \[votefst\] First round this participation key is valid.
	VoteFirstValid uint64 `json:"vote-first-valid,omitempty"`

	// \[votekd\] Number of subkeys in each batch of participation keys.
	VoteKeyDilution uint64 `json:"vote-key-dilution,omitempty"`

	// \[votelst\] Last round this participation key is valid.
	VoteLastValid uint64 `json:"vote-last-valid,omitempty"`

	// \[votekey\] Participation public key used in key registration transactions.
	VoteParticipationKey []byte `json:"vote-participation-key,omitempty"`
}

// TransactionParams contains the parameters that help a client construct a new transaction.
type TransactionParams struct {

	// ConsensusVersion indicates the consensus protocol version
	// as of LastRound.
	ConsensusVersion string `json:"consensus-version"`

	// Fee is the suggested transaction fee
	// Fee is in units of micro-Algos per byte.
	// Fee may fall to zero but transactions must still have a fee of
	// at least MinTxnFee for the current network protocol.
	Fee uint64 `json:"fee"`

	// GenesisID is an ID listed in the genesis block.
	GenesisID string `json:"genesis-id"`

	// GenesisHash is the hash of the genesis block.
	Genesishash []byte `json:"genesis-hash"`

	// LastRound indicates the last round seen
	LastRound uint64 `json:"last-round"`

	// The minimum transaction fee (not per byte) required for the
	// txn to validate for the current network protocol.
	MinFee uint64 `json:"min-fee,omitempty"`
}

// TransactionPayment defines model for TransactionPayment.
type TransactionPayment struct {

	// \[amt\] number of MicroAlgos intended to be transferred.
	Amount uint64 `json:"amount"`

	// Number of MicroAlgos that were sent to the close-remainder-to address when closing the sender account.
	CloseAmount uint64 `json:"close-amount,omitempty"`

	// \[close\] when set, indicates that the sending account should be closed and all remaining funds be transferred to this address.
	CloseRemainderTo string `json:"close-remainder-to,omitempty"`

	// \[rcv\] receiver's address.
	Receiver string `json:"receiver"`
}

// TransactionSignature defines model for TransactionSignature.
type TransactionSignature struct {

	// \[lsig\] Programatic transaction signature.
	//
	// Definition:
	// data/transactions/logicsig.go
	Logicsig TransactionSignatureLogicsig `json:"logicsig,omitempty"`

	// \[msig\] structure holding multiple subsignatures.
	//
	// Definition:
	// crypto/multisig.go : MultisigSig
	Multisig TransactionSignatureMultisig `json:"multisig,omitempty"`

	// \[sig\] Standard ed25519 signature.
	Sig []byte `json:"sig,omitempty"`
}

// TransactionSignatureLogicsig defines model for TransactionSignatureLogicsig.
type TransactionSignatureLogicsig struct {

	// \[arg\] Logic arguments, base64 encoded.
	Args []string `json:"args,omitempty"`

	// \[l\] Program signed by a signature or multi signature, or hashed to be the address of ana ccount. Base64 encoded TEAL program.
	Logic []byte `json:"logic"`

	// \[msig\] structure holding multiple subsignatures.
	//
	// Definition:
	// crypto/multisig.go : MultisigSig
	MultisigSignature TransactionSignatureMultisig `json:"multisig-signature,omitempty"`

	// \[sig\] ed25519 signature.
	Signature []byte `json:"signature,omitempty"`
}

// TransactionSignatureMultisig defines model for TransactionSignatureMultisig.
type TransactionSignatureMultisig struct {

	// \[subsig\] holds pairs of public key and signatures.
	Subsignature []TransactionSignatureMultisigSubsignature `json:"subsignature,omitempty"`

	// \[thr\]
	Threshold uint64 `json:"threshold,omitempty"`

	// \[v\]
	Version uint64 `json:"version,omitempty"`
}

// TransactionSignatureMultisigSubsignature defines model for TransactionSignatureMultisigSubsignature.
type TransactionSignatureMultisigSubsignature struct {

	// \[pk\]
	PublicKey []byte `json:"public-key,omitempty"`

	// \[s\]
	Signature []byte `json:"signature,omitempty"`
}

// Version defines model for Version.
type Version struct {

	// the current algod build version information.
	Build       VersionBuild `json:"build"`
	GenesisHash []byte       `json:"genesis-hash"`
	GenesisId   string       `json:"genesis-id"`
	Versions    []string     `json:"versions"`
}

// VersionBuild defines model for the current algod build version information.
type VersionBuild struct {
	Branch      string `json:"branch"`
	BuildNumber uint64 `json:"build-number"`
	Channel     string `json:"channel"`
	CommitHash  []byte `json:"commit-hash"`
	Major       uint64 `json:"major"`
	Minor       uint64 `json:"minor"`
}

// AccountId defines model for account-id.
type AccountId string

// Address defines model for address.
type Address string

// AddressGreaterThan defines model for address-greater-than.
type AddressGreaterThan string

// AddressRole defines model for address-role.
type AddressRole string

// AfterAddress defines model for after-address.
type AfterAddress string

// AfterAsset defines model for after-asset.
type AfterAsset uint64

// AfterTime defines model for after-time.
type AfterTime string

// AlgosGreaterThan defines model for algos-greater-than.
type AlgosGreaterThan uint64

// AlgosLessThan defines model for algos-less-than.
type AlgosLessThan uint64

// AssetId defines model for asset-id.
type AssetId uint64

// BeforeTime defines model for before-time.
type BeforeTime string

// CurrencyGreaterThan defines model for currency-greater-than.
type CurrencyGreaterThan uint64

// CurrencyLessThan defines model for currency-less-than.
type CurrencyLessThan uint64

// ExcludeCloseTo defines model for exclude-close-to.
type ExcludeCloseTo bool

// Limit defines model for limit.
type Limit uint64

// MaxRound defines model for max-round.
type MaxRound uint64

// MinRound defines model for min-round.
type MinRound uint64

// NotePrefix defines model for note-prefix.
type NotePrefix []byte

// Offset defines model for offset.
type Offset uint64

// Round defines model for round.
type Round uint64

// RoundNumber defines model for round-number.
type RoundNumber uint64

// SigType defines model for sig-type.
type SigType string

// TxId defines model for tx-id.
type TxId []byte

// TxType defines model for tx-type.
type TxType string

// AccountResponse defines model for AccountResponse.
type AccountResponse struct {

	// Account information at a given round.
	//
	// Definition:
	// data/basics/userBalance.go : AccountData
	Account Account `json:"account"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token"`
}

// AccountsResponse defines model for AccountsResponse.
type AccountsResponse struct {
	Accounts []Account `json:"accounts"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token"`
}

// AssetBalancesResponse defines model for AssetBalancesResponse.
type AssetBalancesResponse struct {

	// A simplified version of AssetHolding
	Balances []MiniAssetHolding `json:"balances"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token"`
}

// AssetResponse defines model for AssetResponse.
type AssetResponse struct {

	// Specifies both the unique identifier and the parameters for an asset
	Asset Asset `json:"asset"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`
}

// AssetsResponse defines model for AssetsResponse.
type AssetsResponse struct {
	Assets []Asset `json:"assets"`

	// Round at which the results were computed.
	CurrentRound uint64 `json:"current-round"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token"`
}

// BlockResponse defines model for BlockResponse.
type BlockResponse Block

// HealthCheckResponse defines model for HealthCheckResponse.
type HealthCheckResponse HealthCheck

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Data    *map[string]interface{} `json:"data,omitempty"`
	Message string                  `json:"message"`
}

// HealthCheck defines model for HealthCheck.
type HealthCheck struct {
	Data    *map[string]interface{} `json:"data,omitempty"`
	Message string                  `json:"message"`
}

// TransactionsResponse defines model for TransactionsResponse.
type TransactionsResponse struct {

	// Round at which the results were computed.
	CurrentRound uint64        `json:"current-round"`
	Transactions []Transaction `json:"transactions"`

	// Used for pagination, when making another request provide this token with the next parameter.
	NextToken string `json:"next-token"`
}

// PendingTransactionInfoResponse is returned by Get Pending Transaction by TXID
type PendingTransactionInfoResponse = struct {
	Transaction types.SignedTxn `codec:"txn"`
	//PoolError indicates that the transaction was kicked out of this node's transaction pool (and specifies why that happened).  An empty string indicates the transaction wasn't kicked out of this node's txpool due to an error.
	PoolError string `codec:"pool-error"`
	//ConfirmedRound is the round where this transaction was confirmed, if present.
	ConfirmedRound uint64 `codec:"confirmed-round,omitempty"`
	//AssetIndex is the index of the newly created asset, if this was an asset creation transaction.
	AssetIndex uint64 `codec:"asset-index,omitempty"`
	//Closing amount for the transaction.
	ClosingAmount uint64 `codec:"closing-amount,omitempty"`
	//Rewards, in microAlgos, applied to sender.
	SenderRewards uint64 `codec:"sender-rewards,omitempty"`
	//Rewards, in microAlgos, applied to receiver.
	ReceiverRewards uint64 `codec:"receiver-rewards,omitempty"`
	//Rewards, in microAlgos, applied to close-to.
	CloseRewards uint64 `codec:"close-rewards,omitempty"`
}

// PendingTransactionsResponse is returned by PendingTransactions and by Txid
type PendingTransactionsResponse = struct {
	TopTransactions   []types.SignedTxn `codec:"top-transactions"`
	TotalTransactions uint64            `codec:"total-transactions"`
}

// GetBlock response is returned by Block
type GetBlockResponse = struct {
	Blockb64 string `json:"block"`
}

type LookupAccountByIDResponse struct {
	CurrentRound uint64  `json:"current-round"`
	Account      Account `json:"account"`
}

type LookupAssetByIDResponse struct {
	CurrentRound uint64 `json:"current-round"`
	Asset        Asset  `json:"asset"`
}

type ApplicationResponse struct {
	/**
	 * Application index and its parameters
	 */
	Application Application `json:"application,omitempty"`

	/**
	 * Round at which the results were computed.
	 */
	CurrentRound uint64 `json:"current-round,omitempty"`
}

type ApplicationsResponse struct {
	Applications []Application `json:"applications,omitempty"`

	/**
	 * Round at which the results were computed.
	 */
	CurrentRound uint64 `json:"current-round,omitempty"`

	/**
	 * Used for pagination, when making another request provide this token with the
	 * next parameter.
	 */
	NextToken string `json:"next-token,omitempty"`
}


/**
 * Application index and its parameters
 */
type Application struct {
	/**
	 * (appidx) application index.
	 */
	AppIndex uint64 `json:"app-index,omitempty"`

	/**
	 * (appparams) application parameters.
	 */
	AppParams ApplicationParams `json:"app-params,omitempty"`
}

/**
 * Stores the global information associated with an application.
 */
type ApplicationParams struct {
	/**
	 * (approv) approval program.
	 */
	ApprovalProgram string `json:"approval-program,omitempty"`

	/**
	 * (clearp) approval program.
	 */
	ClearStateProgram string `json:"clear-state-program,omitempty"`

	/**
	 * [\gs) global schema
	 */
	GlobalState []TealKeyValue `json:"global-state,omitempty"`

	/**
	 * [\lsch) global schema
	 */
	GlobalStateSchema ApplicationStateSchema `json:"global-state-schema,omitempty"`

	/**
	 * [\lsch) local schema
	 */
	LocalStateSchema ApplicationStateSchema `json:"local-state-schema,omitempty"`
}


/**
 * Specifies maximums on the number of each type that may be stored.
 */
type ApplicationStateSchema struct {
	/**
	 * (nbs) num of byte slices.
	 */
	NumByteSlice uint64 `json:"num-byte-slice,omitempty"`

	/**
	 * (nui) num of uints.
	 */
	NumUint uint64 `json:"num-uint,omitempty"`
}

/**
 * Pair of application index and application local state
 */
type ApplicationLocalStates struct {
	AppIndex uint64 `json:"app-index,omitempty"`

	/**
	 * Stores local state associated with an application.
	 */
	State ApplicationLocalState `json:"state,omitempty"`
}

/**
 * Stores local state associated with an application.
 */
type ApplicationLocalState struct {
	/**
	 * (tkv) storage.
	 */
	KeyValue []TealKeyValue `json:"key-value,omitempty"`

	/**
	 * (hsch) schema.
	 */
	Schema ApplicationStateSchema `json:"schema,omitempty"`
}


/**
 * Represents a key-value pair in an application store.
 */
type TealKeyValue struct {
	Key string `json:"key,omitempty"`

	/**
	 * Represents a TEAL value.
	 */
	Value TealValue `json:"value,omitempty"`
}

/**
 * Represents a TEAL value.
 */
type TealValue struct {
	/**
	 * (tb) bytes value.
	 */
	Bytes string `json:"bytes,omitempty"`

	/**
	 * (tt) value type.
	 */
	Type uint64 `json:"type,omitempty"`

	/**
	 * (ui) uint value.
	 */
	Uint uint64 `json:"uint,omitempty"`
}
