package crypto

import (
	"github.com/algorand/go-algorand-sdk/v2/internal/signing"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// PQAddress returns the account address for the given pq public key and scheme,
// using the canonical salt for that pair.
func PQAddress(pk []byte, scheme types.PQScheme) (addr types.Address, err error) {
	salt, err := signing.CanonicalSaltForPQPK(pk, scheme)
	if err != nil {
		return
	}
	return signing.PQAddressWithSalt(pk, scheme, salt), nil
}

// PQSignerAddress returns the address for a given PQSigner
func PQSignerAddress(signer PQSigner) (addr types.Address, err error) {
	if signer == nil {
		return types.Address{}, ErrNilPQSigner
	}
	return PQAddress(signer.PQPublicKey(), signer.PQScheme())
}

// SaltForPQSigner returns the canonical salt that will be used when performing
// PQ signatures on behalf of the given signer.
func SaltForPQSigner(sgnr PQSigner) (types.PQAddressSalt, error) {
	return signing.SaltForPQSigner(sgnr)
}
