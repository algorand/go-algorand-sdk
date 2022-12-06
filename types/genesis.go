package types

// A Genesis object defines an Algorand "universe" -- a set of nodes that can
// talk to each other, agree on the ledger contents, etc.  This is defined
// by the initial account states (GenesisAllocation), the initial
// consensus protocol (GenesisProto), and the schema of the ledger.
type Genesis struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Network identifies the unique algorand network for which the ledger
	// is valid.
	// Note the Network name should not include a '-', as we generate the
	// GenesisID from "<Network>-<SchemaID>"; the '-' makes it easy
	// to distinguish between the network and schema.
	Network string `codec:"network"`

	// Proto is the consensus protocol in use at the genesis block.
	Proto string `codec:"proto"`

	// Allocation determines the initial accounts and their state.
	Allocation []interface{} `codec:"alloc,allocbound=MaxInitialGenesisAllocationSize"`

	// RewardsPool is the address of the rewards pool.
	RewardsPool string `codec:"rwd"`

	// FeeSink is the address of the fee sink.
	FeeSink string `codec:"fees"`
}
