package crypto

import (
	"github.com/algorand/go-algorand-sdk/v2/internal/signing"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// PQAddress returns the account address for the given pq public key and scheme,
// using the canonical salt for that pair.
func PQAddress(pk []byte, scheme types.PQScheme) (addr types.Address, err error) {
	salt, err := signing.CanonicalSaltForPQAddress(pk, scheme)
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

// CanonicalSaltForPQAddress returns the canonical salt for the given pq public key
// and scheme: the lowest salt whose derived address is not a valid ed25519
// point. Only the canonical salt is currently supported when signing.
func CanonicalSaltForPQAddress(pk []byte, scheme types.PQScheme) (types.PQAddressSalt, error) {
	return signing.CanonicalSaltForPQAddress(pk, scheme)
}
