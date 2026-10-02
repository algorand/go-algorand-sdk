package signing

import (
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Ed25519RawSignature signs toBeSigned and returns the result as a
// types.Signature, returning ErrInvalidSignatureReturned if the signer produced
// a signature of an unexpected length.
//
// toBeSigned is signed as-is: no domain-separation prefix is added, so callers
// are responsible for building the full byte sequence (see for example
// TransactionBytesToSign).
func Ed25519RawSignature(sgnr Ed25519Signer, toBeSigned []byte) (s types.Signature, err error) {
	signature, err := sgnr.Ed25519Sign(toBeSigned)
	if err != nil {
		return
	}

	if len(signature) != len(s) {
		return s, ErrInvalidSignatureReturned
	}
	copy(s[:], signature)
	return
}

// Ed25519SignTransaction signs tx, returning its txid and the encoded
// SignedTxn. The signer's address is set as the AuthAddr if it is not the
// sender of tx.
func Ed25519SignTransaction(sgnr Ed25519Signer, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	toBeSigned := TransactionBytesToSign(tx)
	sig, err := Ed25519RawSignature(sgnr, toBeSigned)
	if err != nil {
		return
	}

	stxBytes = EncodeSignedTxn(types.SignedTxn{Sig: sig, Txn: tx}, types.Address(sgnr.Ed25519PublicKey()))
	return TxIDFromBytesToSign(toBeSigned), stxBytes, nil
}

// Ed25519SignBytes signs the bytes and returns the signature
func Ed25519SignBytes(sgnr Ed25519Signer, message []byte) (signature []byte, err error) {
	return sgnr.Ed25519Sign(BytesToSign(message))
}

// Ed25519TealSign creates a signature compatible with ed25519verify opcode from
// contract address
func Ed25519TealSign(sgnr Ed25519Signer, data []byte, contractAddress types.Address) (types.Signature, error) {
	return Ed25519RawSignature(sgnr, TealSignData(data, contractAddress))
}
