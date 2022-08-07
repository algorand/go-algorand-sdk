package stateproofcrypto

import (
	cfalcon "github.com/algorand/falcon"
)

const (
	// FalconSeedSize Represents the size in bytes of the random bytes used to generate Falcon keys
	FalconSeedSize = 48

	// FalconMaxSignatureSize Represents the max possible size in bytes of a falcon signature
	FalconMaxSignatureSize = cfalcon.CTSignatureSize
)

type (
	// FalconPublicKey is a wrapper for cfalcon.PublicKeySizey (used for packing)
	FalconPublicKey [cfalcon.PublicKeySize]byte
	// FalconPrivateKey is a wrapper for cfalcon.PrivateKeySize (used for packing)
	FalconPrivateKey [cfalcon.PrivateKeySize]byte
	// FalconSeed represents the seed which is being used to generate Falcon keys
	FalconSeed [FalconSeedSize]byte
	// FalconSignature represents a Falcon signature in a compressed-form
	//msgp:allocbound FalconSignature FalconMaxSignatureSize
	FalconSignature []byte
)

// FalconSigner is the implementation of Signer for the Falcon signature scheme.
type FalconSigner struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	PublicKey  FalconPublicKey  `codec:"pk"`
	PrivateKey FalconPrivateKey `codec:"sk"`
}

// FalconVerifier implements the type Verifier interface for the falcon signature scheme.
type FalconVerifier struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	PublicKey FalconPublicKey `codec:"k"`
}

// VerifyBytes follows falcon algorithm to verify a signature.
func (d *FalconVerifier) VerifyBytes(data []byte, sig FalconSignature) error {
	// The wrapper, currently, support only the compress form signature. so we can
	// assume that the signature given is in a compress form
	falconSig := cfalcon.CompressedSignature(sig)
	return (*cfalcon.PublicKey)(&d.PublicKey).Verify(falconSig, data)
}

// GetFixedLengthHashableRepresentation is used to fetch a plain serialized version of the public data (without the use of the msgpack).
func (d *FalconVerifier) GetFixedLengthHashableRepresentation() []byte {
	return d.PublicKey[:]
}

// GetFixedLengthHashableRepresentation returns a serialized version of the signature
func (s FalconSignature) GetFixedLengthHashableRepresentation() ([]byte, error) {
	compressedSignature := cfalcon.CompressedSignature(s)
	ctSignature, err := compressedSignature.ConvertToCT()
	return ctSignature[:], err
}

// IsSaltVersionEqual of the signature matches the given version
func (s FalconSignature) IsSaltVersionEqual(version byte) bool {
	return (*cfalcon.CompressedSignature)(&s).SaltVersion() == version
}
