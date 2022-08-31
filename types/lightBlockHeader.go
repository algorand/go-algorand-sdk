package types

// A Seed contains cryptographic entropy which can be used to determine a
// committee.
type Seed [32]byte

// LightBlockHeader represents a minimal block header. It contains all the necessary fields
// for verifying proofs on transactions.
// In addition, this struct is designed to be used on environments where only SHA256 function exists
type LightBlockHeader struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Seed                Seed   `codec:"0"`
	RoundNumber         Round  `codec:"r"`
	GenesisHash         Digest `codec:"gh"`
	Sha256TxnCommitment Digest `codec:"tc,allocbound=Sha256Size"`
}
