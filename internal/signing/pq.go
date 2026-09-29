package signing

import (
	"crypto/sha512"
	"fmt"

	"filippo.io/edwards25519"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// pqAddressPrefix is prepended when deriving a post-quantum account address.
var pqAddressPrefix = []byte("PQA")

// PQAddressWithSalt returns the account address for the given pq public key,
// scheme and salt.
// Hash("PQA" || scheme || salt || publicKey)
func PQAddressWithSalt(pk []byte, scheme types.PQScheme, salt types.PQAddressSalt) (addr types.Address) {
	buf := make([]byte, 0, len(pqAddressPrefix)+len(scheme)+1+len(pk))
	buf = append(buf, pqAddressPrefix...)
	buf = append(buf, scheme[:]...)
	buf = append(buf, uint8(salt))
	buf = append(buf, pk[:]...)

	digest := sha512.Sum512_256(buf)

	copy(addr[:], digest[:])
	return
}

// IsEdwards25519Point reports whether encoded can be decoded as an
// Edwards25519 curve point.
func IsEdwards25519Point(encoded []byte) bool {
	if len(encoded) != 32 {
		return false
	}
	_, err := new(edwards25519.Point).SetBytes(encoded)
	return err == nil
}

// CanonicalSaltForPQAddress returns the canonical salt for the given pq public key
// and scheme: the lowest salt whose address cannot be read as a point on the
// ed25519 curve, so that it can only ever be spent by the pq key.
func CanonicalSaltForPQAddress(pk []byte, scheme types.PQScheme) (types.PQAddressSalt, error) {
	for salt := 0; salt <= 0xff; salt++ {
		addr := PQAddressWithSalt(pk, scheme, types.PQAddressSalt(salt))
		if !IsEdwards25519Point(addr[:]) {
			return types.PQAddressSalt(salt), nil
		}
	}

	return 0, fmt.Errorf("no valid salt with an address outside the ed25519 curve exists for %x", pk)
}

// SaltForPQSigner returns the canonical salt that will be used when performing
// PQ signatures on behalf of the given signer.
func SaltForPQSigner(sgnr PQSigner) (types.PQAddressSalt, error) {
	return CanonicalSaltForPQAddress(sgnr.PQPublicKey(), sgnr.PQScheme())
}

// PQSigFor assembles the PQSig envelope that identifies the signer's account
// and carries the given signature bytes, alongside the address of that account.
//
// The signature is not checked, and may be empty: callers that need an envelope
// without a real signature (to size or simulate a transaction, say) can pass
// nil.
func PQSigFor(sgnr PQSigner, signature []byte) (sig types.PQSig, addr types.Address, err error) {
	salt, err := SaltForPQSigner(sgnr)
	if err != nil {
		return
	}

	sig = types.PQSig{
		Scheme:    sgnr.PQScheme(),
		Salt:      salt,
		PublicKey: sgnr.PQPublicKey(),
		Signature: signature,
	}
	addr = PQAddressWithSalt(sig.PublicKey, sig.Scheme, salt)
	return
}

// PQSignedTxn returns the encoded SignedTxn carrying the given post-quantum
// signature bytes on behalf of the signer's account.
//
// The signature may be empty, which produces an envelope suitable for
// simulating transactions with the allowEmptySignatures option enabled.
func PQSignedTxn(sgnr PQSigner, tx types.Transaction, signature []byte) ([]byte, error) {
	pqsig, address, err := PQSigFor(sgnr, signature)
	if err != nil {
		return nil, err
	}

	return EncodeSignedTxn(types.SignedTxn{Txn: tx, PQsig: pqsig}, address), nil
}

// PQSignTransaction signs tx, returning the encoded SignedTxn.
func PQSignTransaction(sgnr PQSigner, tx types.Transaction) ([]byte, error) {
	signature, err := sgnr.PQSign(TransactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	return PQSignedTxn(sgnr, tx, signature)
}
