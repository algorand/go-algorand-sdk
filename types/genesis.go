package types

import (
	"crypto/sha512"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
)

// GenesisHashID is the Genesis HashID defined in go-algorand/protocol/hash.go
var GenesisHashID = "GE"

// A Genesis object defines an Algorand "universe"
type Genesis struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// The SchemaID allows nodes to store data specific to a particular
	// universe
	SchemaID string `codec:"id"`

	// Network identifies the unique algorand network for which the ledger
	// is valid.
	Network string `codec:"network"`

	// Proto is the consensus protocol in use at the genesis block.
	Proto string `codec:"proto"`

	// Allocation determines the initial accounts and their state.
	Allocation []GenesisAllocation `codec:"alloc"`

	// RewardsPool is the address of the rewards pool.
	RewardsPool string `codec:"rwd"`

	// FeeSink is the address of the fee sink.
	FeeSink string `codec:"fees"`

	// Timestamp for the genesis block
	Timestamp int64 `codec:"timestamp"`

	// Arbitrary genesis comment string - will be excluded from file if empty
	Comment string `codec:"comment"`

	// DevMode defines whether this network operates in a developer mode or not. Developer mode networks
	// are a single node network, that operates without the agreement service being active. In liue of the
	// agreement service, a new block is generated each time a node receives a transaction group. The
	// default value for this field is "false", which makes this field empty from it's encoding, and
	// therefore backward compatible.
	DevMode bool `codec:"devmode"`
}

// GenesisAllocation object represents an allocation of algos to
// an address in the genesis block.  Address is the checksummed
// short address.  Comment is a note about what this address is
// representing, and is purely informational.  State is the initial
// account state.
type GenesisAllocation struct {
	_struct struct{} `codec:""`

	Address string  `codec:"addr"`
	Comment string  `codec:"comment"`
	State   Account `codec:"state"`
}

// Account contains the data associated with a given address.
// This includes the account balance, cryptographic public keys,
// consensus delegation status.
type Account struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Status          byte     `codec:"onl"`
	MicroAlgos      uint64   `codec:"algo"`
	VoteID          [32]byte `codec:"vote"`
	SelectionID     [32]byte `codec:"sel"`
	StateProofID    [64]byte `codec:"stprf"`
	VoteLastValid   uint64   `codec:"voteLst"`
	VoteKeyDilution uint64   `codec:"voteKD"`
}

// ID is the effective Genesis identifier - the combination
// of the network and the ledger schema version
func (genesis Genesis) ID() string {
	return string(genesis.Network) + "-" + genesis.SchemaID
}

// Hash is the genesis hash.
func (genesis Genesis) Hash() Digest {
	hashRep := []byte(GenesisHashID)
	data := msgpack.Encode(genesis)
	hashRep = append(hashRep, data...)
	return sha512.Sum512_256(hashRep)
}

// GenesisBalances contains the information needed to generate a new ledger
type GenesisBalances struct {
	Balances    map[Address]Account
	FeeSink     Address
	RewardsPool Address
	Timestamp   int64
}

// Balances returns the genesis account balances.
func (genesis Genesis) Balances() (GenesisBalances, error) {
	genalloc := make(map[Address]Account)
	for _, entry := range genesis.Allocation {
		addr, err := DecodeAddress(entry.Address)
		if err != nil {
			return GenesisBalances{}, fmt.Errorf("cannot parse genesis addr %s: %w", entry.Address, err)
		}

		_, present := genalloc[addr]
		if present {
			return GenesisBalances{}, fmt.Errorf("repeated allocation to %s", entry.Address)
		}

		genalloc[addr] = entry.State
	}

	feeSink, err := DecodeAddress(genesis.FeeSink)
	if err != nil {
		return GenesisBalances{}, fmt.Errorf("cannot parse fee sink addr %s: %w", genesis.FeeSink, err)
	}

	rewardsPool, err := DecodeAddress(genesis.RewardsPool)
	if err != nil {
		return GenesisBalances{}, fmt.Errorf("cannot parse rewards pool addr %s: %w", genesis.RewardsPool, err)
	}

	return MakeTimestampedGenesisBalances(genalloc, feeSink, rewardsPool, genesis.Timestamp), nil
}

// MakeTimestampedGenesisBalances returns the information needed to bootstrap the ledger based on a given time
func MakeTimestampedGenesisBalances(balances map[Address]Account, feeSink, rewardsPool Address, timestamp int64) GenesisBalances {
	return GenesisBalances{Balances: balances, FeeSink: feeSink, RewardsPool: rewardsPool, Timestamp: timestamp}
}
