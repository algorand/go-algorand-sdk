package types

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
)

// MessageHash represents the message that a state proof will attest to.
type MessageHash [32]byte

// StateProofType identifies a particular configuration of state proofs.
type StateProofType uint64

const (
	// StateProofBasic is our initial state proof setup. using falcon keys and subset-sum hash
	StateProofBasic StateProofType = 0

	// NumStateProofTypes is the max number of types of state proofs
	// that we support.  This is used as an allocation bound for a map
	// containing different stateproof types in msgpack encoding.
	NumStateProofTypes int = 1

	// MaxReveals is a bound on allocation and on numReveals to limit log computation
	MaxReveals int = 640

	// MaxEncodedTreeDepth is the maximum tree depth (root only depth 0) for a tree which
	// is being encoded (either by msbpack or by the fixed length encoding)
	MaxEncodedTreeDepth = 16

	// MaxNumLeavesOnEncodedTree is the maximum number of leaves allowed for a tree which
	// is being encoded (either by msbpack or by the fixed length encoding)
	MaxNumLeavesOnEncodedTree = 1 << MaxEncodedTreeDepth
)

// GenericDigest is a digest that implements CustomSizeDigest, and can be used as hash output.
//msgp:allocbound GenericDigest MaxHashDigestSize
type GenericDigest []byte

// ToSlice is used inside the Tree itself when interacting with TreeDigest
func (d GenericDigest) ToSlice() []byte { return d }

// IsEqual compare two digests
func (d GenericDigest) IsEqual(other GenericDigest) bool {
	return bytes.Equal(d, other)
}

// IsEmpty checks wether the generic digest is an empty one or not
func (d GenericDigest) IsEmpty() bool {
	return len(d) == 0
}

// Sumhash512DigestSize  The size in bytes of the sumhash checksum
const Sumhash512DigestSize = 64

// Sizes of each hash
const (
	Sha512_256Size    = sha512.Size256
	SumhashDigestSize = Sumhash512DigestSize
	Sha256Size        = sha256.Size
)

// HashType represents different hash functions
type HashType uint16

// types of hashes
const (
	Sha512_256 HashType = iota
	Sumhash
	Sha256
	MaxHashType
)

// HashFactory is responsible for generating new hashes accordingly to the type it stores.
//msgp:postunmarshalcheck HashFactory Validate
type HashFactory struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	HashType HashType `codec:"t"`
}

// Proof is used to convince a verifier about membership of leaves: h0,h1...hn
// at indexes i0,i1...in on a tree. The verifier has a trusted value of the tree
// root hash.
type Proof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Path is bounded by MaxNumLeavesOnEncodedTree since there could be multiple reveals, and
	// given the distribution of the elt positions and the depth of the tree,
	// the path length can increase up to 2^MaxEncodedTreeDepth / 2
	Path        []GenericDigest `codec:"pth,allocbound=MaxNumLeavesOnEncodedTree/2"`
	HashFactory HashFactory     `codec:"hsh"`
	// TreeDepth represents the depth of the tree that is being proven.
	// It is the number of edges from the root to a leaf.
	TreeDepth uint8 `codec:"td"`
}

// MerkleSignatureSchemeRootSize is the size of the root of the merkle tree.
const MerkleSignatureSchemeRootSize = SumhashDigestSize

// Commitment represents the root of the vector commitment tree built upon the MSS keys.
type Commitment [MerkleSignatureSchemeRootSize]byte

// Verifier is used to verify a merklesignature.Signature produced by merklesignature.Secrets.
type Verifier struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Commitment  Commitment `codec:"cmt"`
	KeyLifetime uint64     `codec:"lf"`
}

// A Participant corresponds to an account whose AccountData.Status
// is Online, and for which the expected sigRound satisfies
// AccountData.VoteFirstValid <= sigRound <= AccountData.VoteLastValid.
//
// In the Algorand ledger, it is possible for multiple accounts to have
// the same PK.  Thus, the PK is not necessarily unique among Participants.
// However, each account will produce a unique Participant struct, to avoid
// potential DoS attacks where one account claims to have the same VoteID PK
// as another account.
type Participant struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// PK is the identifier used to verify the signature for a specific participant
	PK Verifier `codec:"p"`

	// Weight is AccountData.MicroAlgos.
	Weight uint64 `codec:"w"`
}

// MerkleSignature represents a Falcon signature in a compressed-form
//msgp:allocbound MerkleSignature FalconMaxSignatureSize
type MerkleSignature []byte

// SingleLeafProof is used to convince a verifier about membership of a specific
// leaf h at index i on a tree. The verifier has a trusted value of the tree
// root hash. it corresponds to merkle verification path.
type SingleLeafProof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Proof
}

// FalconPublicKeySize pulled out of falcon.go
const FalconPublicKeySize = 0x701

// FalconPublicKey is a wrapper for cfalcon.PublicKeySizey (used for packing)
type FalconPublicKey [FalconPublicKeySize]byte

// FalconVerifier implements the type Verifier interface for the falcon signature scheme.
type FalconVerifier struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	PublicKey FalconPublicKey `codec:"k"`
}

// FalconSignatureStruct represents a signature in the merkle signature scheme using falcon signatures as an underlying crypto scheme.
// It consists of an ephemeral public key, a signature, a merkle verification path and an index.
// The merkle signature considered valid only if the Signature is verified under the ephemeral public key and
// the Merkle verification path verifies that the ephemeral public key is located at the given index of the tree
// (for the root given in the long-term public key).
// More details can be found on Algorand's spec
type FalconSignatureStruct struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Signature             MerkleSignature `codec:"sig"`
	VectorCommitmentIndex uint64          `codec:"idx"`
	Proof                 SingleLeafProof `codec:"prf"`
	VerifyingKey          FalconVerifier  `codec:"vkey"`
}

// A sigslotCommit is a single slot in the sigs array that forms the state proof.
type sigslotCommit struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Sig is a signature by the participant on the expected message.
	Sig FalconSignatureStruct `codec:"s"`

	// L is the total weight of signatures in lower-numbered slots.
	// This is initialized once the builder has collected a sufficient
	// number of signatures.
	L uint64 `codec:"l"`
}

// Reveal is a single array position revealed as part of a state
// proof.  It reveals an element of the signature array and
// the corresponding element of the participants array.
type Reveal struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	SigSlot sigslotCommit `codec:"s"`
	Part    Participant   `codec:"p"`
}

// StateProof represents a proof on Algorand's state.
type StateProof struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	SigCommit                  GenericDigest `codec:"c"`
	SignedWeight               uint64        `codec:"w"`
	SigProofs                  Proof         `codec:"S"`
	PartProofs                 Proof         `codec:"P"`
	MerkleSignatureSaltVersion byte          `codec:"v"`
	// Reveals is a sparse map from the position being revealed
	// to the corresponding elements from the sigs and participants
	// arrays.
	Reveals           map[uint64]Reveal `codec:"r,allocbound=MaxReveals"`
	PositionsToReveal []uint64          `codec:"pr,allocbound=MaxReveals"`
}

// Message represents the message that the state proofs are attesting to. This message can be
// used by lightweight client and gives it the ability to verify proofs on the Algorand's state.
// In addition to that proof, this message also contains fields that
// are needed in order to verify the next state proofs (VotersCommitment and LnProvenWeight).
type Message struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`
	// BlockHeadersCommitment contains a commitment on all light block headers within a state proof interval.
	BlockHeadersCommitment []byte `codec:"b,allocbound=Sha256Size"`
	VotersCommitment       []byte `codec:"v,allocbound=SumhashDigestSize"`
	LnProvenWeight         uint64 `codec:"P"`
	FirstAttestedRound     uint64 `codec:"f"`
	LastAttestedRound      uint64 `codec:"l"`
}

// StateProofTxnFields captures the fields used for stateproof transactions.
type StateProofTxnFields struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	StateProofType StateProofType `codec:"sptype"`
	StateProof     StateProof     `codec:"sp"`
	Message        Message        `codec:"spmsg"`
}
