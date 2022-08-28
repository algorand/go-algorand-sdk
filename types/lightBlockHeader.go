package types

import (
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
)

const BlockHeader256 HashID = "B256"

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

// ToBeHashed implements the crypto.Hashable interface
func (bh LightBlockHeader) ToBeHashed() []byte {
	return append([]byte(BlockHeader256), msgpack.Encode(&bh)...)
}
