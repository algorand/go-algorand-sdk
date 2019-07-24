package config

import (
	"time"
)

// ConsensusParams specifies settings that might vary based on the
// particular version of the consensus protocol.
type ConsensusParams struct {
	// Consensus protocol upgrades.  Votes for upgrades are collected for
	// UpgradeVoteRounds.  If the number of positive votes is over
	// UpgradeThreshold, the proposal is accepted.
	//
	// UpgradeVoteRounds needs to be long enough to collect an
	// accurate sample of participants, and UpgradeThreshold needs
	// to be high enough to ensure that there are sufficient participants
	// after the upgrade.
	//
	// There is a delay of UpgradeWaitRounds between approval of
	// an upgrade and its deployment, to give clients time to notify users.
	UpgradeVoteRounds   uint64
	UpgradeThreshold    uint64
	UpgradeWaitRounds   uint64
	MaxVersionStringLen int

	// MaxTxnBytesPerBlock determines the maximum number of bytes
	// that transactions can take up in a block.  Specifically,
	// the sum of the lengths of encodings of each transaction
	// in a block must not exceed MaxTxnBytesPerBlock.
	MaxTxnBytesPerBlock int

	// MaxTxnBytesPerBlock is the maximum size of a transaction's Note field.
	MaxTxnNoteBytes int

	// MaxTxnLife is how long a transaction can be live for:
	// the maximum difference between LastValid and FirstValid.
	//
	// Note that in a protocol upgrade, the ledger must first be upgraded
	// to hold more past blocks for this value to be raised.
	MaxTxnLife uint64

	// ApprovedUpgrades describes the upgrade proposals that this protocol
	// implementation will vote for.
	ApprovedUpgrades map[ConsensusVersion]bool

	// SupportGenesisHash indicates support for the GenesisHash
	// fields in transactions (and requires them in blocks).
	SupportGenesisHash bool

	// RequireGenesisHash indicates that GenesisHash must be present
	// in every transaction.
	RequireGenesisHash bool

	// DefaultKeyDilution specifies the granularity of top-level ephemeral
	// keys. KeyDilution is the number of second-level keys in each batch,
	// signed by a top-level "batch" key.  The default value can be
	// overriden in the account state.
	DefaultKeyDilution uint64

	// MinBalance specifies the minimum balance that can appear in
	// an account.  To spend money below MinBalance requires issuing
	// an account-closing transaction, which transfers all of the
	// money from the account, and deletes the account state.
	MinBalance uint64

	// MinTxnFee specifies the minimum fee allowed on a transaction.
	// A minimum fee is necessary to prevent DoS. In some sense this is
	// a way of making the spender subsidize the cost of storing this transaction.
	MinTxnFee uint64

	// RewardUnit specifies the number of MicroAlgos corresponding to one reward
	// unit.
	//
	// Rewards are received by whole reward units.  Fractions of
	// RewardUnits do not receive rewards.
	RewardUnit uint64

	// RewardsRateRefreshInterval is the number of rounds after which the
	// rewards level is recomputed for the next RewardsRateRefreshInterval rounds.
	RewardsRateRefreshInterval uint64

	// seed-related parameters
	SeedLookback        uint64 // how many blocks back we use seeds from in sortition. delta_s in the spec
	SeedRefreshInterval uint64 // how often an old block hash is mixed into the seed. delta_r in the spec

	// ledger retention policy
	MaxBalLookback uint64 // (current round - MaxBalLookback) is the oldest round the ledger must answer balance queries for

	// sortition threshold factors
	NumProposers           uint64
	SoftCommitteeSize      uint64
	SoftCommitteeThreshold uint64
	CertCommitteeSize      uint64
	CertCommitteeThreshold uint64
	NextCommitteeSize      uint64 // for any non-FPR votes >= deadline step, committee sizes and thresholds are constant
	NextCommitteeThreshold uint64
	LateCommitteeSize      uint64
	LateCommitteeThreshold uint64
	RedoCommitteeSize      uint64
	RedoCommitteeThreshold uint64
	DownCommitteeSize      uint64
	DownCommitteeThreshold uint64

	FastRecoveryLambda    time.Duration // time between fast recovery attempts
	FastPartitionRecovery bool          // set when fast partition recovery is enabled

	// commit to payset using a hash of entire payset,
	// instead of txid merkle tree
	PaysetCommitFlat bool

	MaxTimestampIncrement int64 // maximum time between timestamps on successive blocks

	// support for the efficient encoding in SignedTxnInBlock
	SupportSignedTxnInBlock bool

	// force the FeeSink address to be non-participating in the genesis balances.
	ForceNonParticipatingFeeSink bool

	// support for ApplyData in SignedTxnInBlock
	ApplyData bool

	// track reward distributions in ApplyData
	RewardsInApplyData bool

	// domain-separated credentials
	CredentialDomainSeparationEnabled bool
}

// ConsensusVersion is a string that identifies a version of the
// consensus protocol.
type ConsensusVersion string

// ConsensusV7 increases MaxBalLookback to 320 in preparation for
// the twin seeds change.
const ConsensusV7 = ConsensusVersion("v7")

// ConsensusV8 uses the new parameters and seed derivation policy
// from the agreement protocol's security analysis.
const ConsensusV8 = ConsensusVersion("v8")

// ConsensusV9 increases min balance to 100,000 microAlgos.
const ConsensusV9 = ConsensusVersion("v9")

// ConsensusV10 introduces fast partition recovery.
const ConsensusV10 = ConsensusVersion("v10")

// ConsensusV11 introduces efficient encoding of SignedTxn using SignedTxnInBlock.
const ConsensusV11 = ConsensusVersion("v11")

// ConsensusV12 increases the maximum length of a version string.
const ConsensusV12 = ConsensusVersion("v12")

// ConsensusV13 makes the consensus version a meaningful string.
const ConsensusV13 = ConsensusVersion(
	// Points to version of the Algorand spec as of May 21, 2019.
	"https://github.com/algorand/spec/tree/0c8a9dc44d7368cc266d5407b79fb3311f4fc795",
)

// ConsensusV14 adds tracking of closing amounts in ApplyData,
// and enables genesis hash in transactions.
const ConsensusV14 = ConsensusVersion(
	"https://github.com/algorand/spec/tree/2526b6ae062b4fe5e163e06e41e1d9b9219135a9",
)

// ConsensusV15 adds tracking of reward distributions in ApplyData.
const ConsensusV15 = ConsensusVersion(
	"https://github.com/algorand/spec/tree/a26ed78ed8f834e2b9ccb6eb7d3ee9f629a6e622",
)

// ConsensusV16 fixes domain separation in Credentials and requires GenesisHash.
const ConsensusV16 = ConsensusVersion(
	"https://github.com/algorand/spec/tree/22726c9dcd12d9cddce4a8bd7e8ccaa707f74101",
)

// ConsensusV17 points to 'final' spec commit
const ConsensusV17 = ConsensusVersion(
	"https://github.com/algorandfoundation/specs/tree/5615adc36bad610c7f165fa2967f4ecfa75125f0",
)

// Consensus tracks the protocol-level settings for different versions of the
// consensus protocol.
var Consensus map[ConsensusVersion]ConsensusParams

func init() {
	Consensus = make(map[ConsensusVersion]ConsensusParams)

	// WARNING: copying a ConsensusParams by value into a new variable
	// does not copy the ApprovedUpgrades map.  Make sure that each new
	// ConsensusParams structure gets a fresh ApprovedUpgrades map.

	// Base consensus protocol version, v7.
	v7 := ConsensusParams{
		UpgradeVoteRounds:   10000,
		UpgradeThreshold:    9000,
		UpgradeWaitRounds:   10000,
		MaxVersionStringLen: 64,

		MinBalance:          10000,
		MinTxnFee:           1000,
		MaxTxnLife:          1000,
		MaxTxnNoteBytes:     1024,
		MaxTxnBytesPerBlock: 1000000,
		DefaultKeyDilution:  10000,

		MaxTimestampIncrement: 25,

		RewardUnit:                 1e6,
		RewardsRateRefreshInterval: 5e5,

		ApprovedUpgrades: map[ConsensusVersion]bool{},

		NumProposers:           30,
		SoftCommitteeSize:      2500,
		SoftCommitteeThreshold: 1870,
		CertCommitteeSize:      1000,
		CertCommitteeThreshold: 720,
		NextCommitteeSize:      10000,
		NextCommitteeThreshold: 7750,
		LateCommitteeSize:      10000,
		LateCommitteeThreshold: 7750,
		RedoCommitteeSize:      10000,
		RedoCommitteeThreshold: 7750,
		DownCommitteeSize:      10000,
		DownCommitteeThreshold: 7750,

		FastRecoveryLambda: 5 * time.Minute,

		SeedLookback:        2,
		SeedRefreshInterval: 100,

		MaxBalLookback: 320,
	}

	v7.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV7] = v7

	// v8 uses parameters and a seed derivation policy (the "twin seeds") from Georgios' new analysis
	v8 := v7

	v8.SeedRefreshInterval = 80
	v8.NumProposers = 9
	v8.SoftCommitteeSize = 2990
	v8.SoftCommitteeThreshold = 2267
	v8.CertCommitteeSize = 1500
	v8.CertCommitteeThreshold = 1112
	v8.NextCommitteeSize = 5000
	v8.NextCommitteeThreshold = 3838
	v8.LateCommitteeSize = 5000
	v8.LateCommitteeThreshold = 3838
	v8.RedoCommitteeSize = 5000
	v8.RedoCommitteeThreshold = 3838
	v8.DownCommitteeSize = 5000
	v8.DownCommitteeThreshold = 3838

	v8.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV8] = v8

	// v7 can be upgraded to v8.
	v7.ApprovedUpgrades[ConsensusV8] = true

	// v9 increases the minimum balance to 100,000 microAlgos.
	v9 := v8
	v9.MinBalance = 100000
	v9.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV9] = v9

	// v8 can be upgraded to v9.
	v8.ApprovedUpgrades[ConsensusV9] = true

	// v10 introduces fast partition recovery (and also raises NumProposers).
	v10 := v9
	v10.FastPartitionRecovery = true
	v10.NumProposers = 20
	v10.LateCommitteeSize = 500
	v10.LateCommitteeThreshold = 320
	v10.RedoCommitteeSize = 2400
	v10.RedoCommitteeThreshold = 1768
	v10.DownCommitteeSize = 6000
	v10.DownCommitteeThreshold = 4560
	v10.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV10] = v10

	// v9 can be upgraded to v10.
	v9.ApprovedUpgrades[ConsensusV10] = true

	// v11 introduces SignedTxnInBlock.
	v11 := v10
	v11.SupportSignedTxnInBlock = true
	v11.PaysetCommitFlat = true
	v11.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV11] = v11

	// v10 can be upgraded to v11.
	v10.ApprovedUpgrades[ConsensusV11] = true

	// v12 increases the maximum length of a version string.
	v12 := v11
	v12.MaxVersionStringLen = 128
	v12.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV12] = v12

	// v11 can be upgraded to v12.
	v11.ApprovedUpgrades[ConsensusV12] = true

	// v13 makes the consensus version a meaningful string.
	v13 := v12
	v13.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV13] = v13

	// v12 can be upgraded to v13.
	v12.ApprovedUpgrades[ConsensusV13] = true

	// v14 introduces tracking of closing amounts in ApplyData, and enables
	// GenesisHash in transactions.
	v14 := v13
	v14.ApplyData = true
	v14.SupportGenesisHash = true
	v14.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV14] = v14

	// v13 can be upgraded to v14.
	v13.ApprovedUpgrades[ConsensusV14] = true

	// v15 introduces tracking of reward distributions in ApplyData.
	v15 := v14
	v15.RewardsInApplyData = true
	v15.ForceNonParticipatingFeeSink = true
	v15.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV15] = v15

	// v14 can be upgraded to v15.
	v14.ApprovedUpgrades[ConsensusV15] = true

	// v16 fixes domain separation in credentials.
	v16 := v15
	v16.CredentialDomainSeparationEnabled = true
	v16.RequireGenesisHash = true
	v16.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV16] = v16

	// v15 can be upgraded to v16.
	v15.ApprovedUpgrades[ConsensusV16] = true

	// ConsensusV17 points to 'final' spec commit
	v17 := v16
	v17.ApprovedUpgrades = map[ConsensusVersion]bool{}
	Consensus[ConsensusV17] = v17

	// v16 can be upgraded to v17.
	v16.ApprovedUpgrades[ConsensusV17] = true
}