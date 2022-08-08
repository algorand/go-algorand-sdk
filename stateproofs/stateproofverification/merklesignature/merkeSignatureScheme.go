package merklesignature

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/stateproofs/stateprooftypes"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/merklearray"
	"github.com/algorand/go-algorand-sdk/stateproofs/stateproofverification/stateproofcrypto"
)

// Errors for the merkle signature scheme
var (
	ErrKeyLifetimeIsZero                 = errors.New("received zero KeyLifetime")
	ErrSignatureSchemeVerificationFailed = errors.New("merkle signature verification failed")
	ErrSignatureSaltVersionMismatch      = errors.New("the signature's salt version does not match")
)

type (
	// Signature represents a signature in the merkle signature scheme using falcon signatures as an underlying crypto scheme.
	// It consists of an ephemeral public key, a signature, a merkle verification path and an index.
	// The merkle signature considered valid only if the Signature is verified under the ephemeral public key and
	// the Merkle verification path verifies that the ephemeral public key is located at the given index of the tree
	// (for the root given in the long-term public key).
	// More details can be found on Algorand's spec
	Signature struct {
		_struct struct{} `codec:",omitempty,omitemptyarray"`

		Signature             stateproofcrypto.FalconSignature `codec:"sig"`
		VectorCommitmentIndex uint64                           `codec:"idx"`
		Proof                 merklearray.SingleLeafProof      `codec:"prf"`
		VerifyingKey          stateproofcrypto.FalconVerifier  `codec:"vkey"`
	}

	// Commitment represents the root of the vector commitment tree built upon the MSS keys.
	Commitment [MerkleSignatureSchemeRootSize]byte

	// Verifier is used to verify a merklesignature.Signature produced by merklesignature.Secrets.
	Verifier struct {
		_struct struct{} `codec:",omitempty,omitemptyarray"`

		Commitment  Commitment `codec:"cmt"`
		KeyLifetime uint64     `codec:"lf"`
	}
)

// ValidateSaltVersion validates that the version of the signature is matching the expected version
func (s *Signature) ValidateSaltVersion(version byte) error {
	if !s.Signature.IsSaltVersionEqual(version) {
		return ErrSignatureSaltVersionMismatch
	}
	return nil
}

// FirstRoundInKeyLifetime calculates the round of the valid key for a given round by lowering to the closest KeyLiftime divisor.
func (v *Verifier) FirstRoundInKeyLifetime(round uint64) (uint64, error) {
	if v.KeyLifetime == 0 {
		return 0, ErrKeyLifetimeIsZero
	}

	return firstRoundInKeyLifetime(round, v.KeyLifetime), nil
}

// VerifyBytes verifies that a merklesignature sig is valid, on a specific round, under a given public key
func (v *Verifier) VerifyBytes(round uint64, msg []byte, sig *Signature) error {
	validKeyRound, err := v.FirstRoundInKeyLifetime(round)
	if err != nil {
		return err
	}

	ephkey := CommittablePublicKey{
		VerifyingKey: sig.VerifyingKey,
		Round:        validKeyRound,
	}

	// verify the merkle tree verification path using the ephemeral public key, the
	// verification path and the index.
	err = merklearray.VerifyVectorCommitment(
		v.Commitment[:],
		map[uint64]stateprooftypes.Hashable{sig.VectorCommitmentIndex: &ephkey},
		sig.Proof.ToProof(),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSignatureSchemeVerificationFailed, err)
	}

	// verify that the signature is valid under the ephemeral public key
	err = sig.VerifyingKey.VerifyBytes(msg, sig.Signature)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSignatureSchemeVerificationFailed, err)
	}
	return nil
}

// GetFixedLengthHashableRepresentation returns the signature as a hashable byte sequence.
// the format details can be found in the Algorand's spec.
func (s *Signature) GetFixedLengthHashableRepresentation() ([]byte, error) {
	var schemeType [2]byte
	binary.LittleEndian.PutUint16(schemeType[:], CryptoPrimitivesID)
	sigBytes, err := s.Signature.GetFixedLengthHashableRepresentation()
	if err != nil {
		return nil, err
	}

	verifierBytes := s.VerifyingKey.GetFixedLengthHashableRepresentation()

	var binaryVectorCommitmentIndex [8]byte
	binary.LittleEndian.PutUint64(binaryVectorCommitmentIndex[:], s.VectorCommitmentIndex)

	proofBytes := s.Proof.GetFixedLengthHashableRepresentation()

	merkleSignatureBytes := make([]byte, 0, len(schemeType)+len(sigBytes)+len(verifierBytes)+len(binaryVectorCommitmentIndex)+len(proofBytes))
	merkleSignatureBytes = append(merkleSignatureBytes, schemeType[:]...)
	merkleSignatureBytes = append(merkleSignatureBytes, sigBytes...)
	merkleSignatureBytes = append(merkleSignatureBytes, verifierBytes...)
	merkleSignatureBytes = append(merkleSignatureBytes, binaryVectorCommitmentIndex[:]...)
	merkleSignatureBytes = append(merkleSignatureBytes, proofBytes...)
	return merkleSignatureBytes, nil
}
