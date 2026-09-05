package transaction

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

var errInvalidSignatureReturned = errors.New("ed25519 signer returned an invalid signature")

// SignTransaction signs one transaction with the provided TransactionSigner.
func SignTransaction(signer TransactionSigner, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	stxs, err := signer.SignTransactions([]types.Transaction{tx}, []int{0})
	if err != nil {
		return "", nil, err
	}
	if len(stxs) != 1 {
		return "", nil, fmt.Errorf("transaction signer returned %d signed transactions, expected 1", len(stxs))
	}
	return crypto.GetTxID(tx), stxs[0], nil
}

// signTransactions signs the transactions of the group selected by
// indexesToSign with signOne, returning the encoded signed transactions in the
// same order as indexesToSign.
func signTransactions(txGroup []types.Transaction, indexesToSign []int, signOne func(types.Transaction) ([]byte, error)) ([][]byte, error) {
	stxs := make([][]byte, len(indexesToSign))
	for i, pos := range indexesToSign {
		stxBytes, err := signOne(txGroup[pos])
		if err != nil {
			return nil, err
		}

		stxs[i] = stxBytes
	}

	return stxs, nil
}

// transactionBytesToSign returns the byte form of the tx that we actually sign.
func transactionBytesToSign(tx types.Transaction) []byte {
	return append([]byte("TX"), msgpack.Encode(tx)...)
}

// encodeSignedTxn encodes the SignedTxn, assigning the signing account as the
// AuthAddr when it is not the sender of the transaction.
func encodeSignedTxn(stx types.SignedTxn, signerAddress types.Address) []byte {
	if stx.Txn.Sender != signerAddress {
		stx.AuthAddr = signerAddress
	}
	return msgpack.Encode(&stx)
}

// equalBySerialization reports whether two signers hold equivalent parameters,
// by comparing their JSON encodings.
func equalBySerialization(signer, other interface{}) bool {
	signerJSON, err := json.Marshal(signer)
	if err != nil {
		return false
	}

	otherJSON, err := json.Marshal(other)
	if err != nil {
		return false
	}

	return string(signerJSON) == string(otherJSON)
}

// ed25519Signature signs the given bytes and returns the result as a
// types.Signature, erroring out if the signer returned a signature of an
// unexpected length.
func ed25519Signature(signer crypto.Ed25519Signer, toBeSigned []byte) (sig types.Signature, err error) {
	signature, err := signer.Ed25519Sign(toBeSigned)
	if err != nil {
		return
	}

	if copy(sig[:], signature) != len(sig) {
		err = errInvalidSignatureReturned
	}
	return
}

func ed25519SignTransaction(signer crypto.Ed25519Signer, tx types.Transaction) ([]byte, error) {
	sig, err := ed25519Signature(signer, transactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	return encodeSignedTxn(types.SignedTxn{Sig: sig, Txn: tx}, types.Address(signer.Ed25519PublicKey())), nil
}

func ed25519SignMultisigTransaction(signer crypto.Ed25519Signer, account crypto.MultisigAccount, tx types.Transaction) ([]byte, error) {
	if err := account.Validate(); err != nil {
		return nil, err
	}

	publicKey := signer.Ed25519PublicKey()
	signerIndex := -1
	msig := types.MultisigSig{
		Version:   account.Version,
		Threshold: account.Threshold,
		Subsigs:   make([]types.MultisigSubsig, len(account.Pks)),
	}
	for i, key := range account.Pks {
		msig.Subsigs[i].Key = append([]byte(nil), key...)
		if bytes.Equal(publicKey[:], key) {
			signerIndex = i
		}
	}
	if signerIndex == -1 {
		return nil, errors.New("secret key has no corresponding public identity in multisig preimage")
	}

	sig, err := ed25519Signature(signer, transactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}
	msig.Subsigs[signerIndex].Sig = sig

	address, err := account.Address()
	if err != nil {
		return nil, err
	}
	return encodeSignedTxn(types.SignedTxn{Msig: msig, Txn: tx}, address), nil
}

func ed25519AppendMultisigTransaction(signer crypto.Ed25519Signer, account crypto.MultisigAccount, encoded []byte) (txid string, stxBytes []byte, err error) {
	var stx types.SignedTxn
	if err = msgpack.Decode(encoded, &stx); err != nil {
		return
	}
	partial, err := ed25519SignMultisigTransaction(signer, account, stx.Txn)
	if err != nil {
		return "", nil, err
	}
	return crypto.MergeMultisigTransactions(partial, encoded)
}

func signLogicSigAccountTransaction(account crypto.LogicSigAccount, tx types.Transaction) ([]byte, error) {
	address, err := account.Address()
	if err != nil {
		return nil, err
	}
	if !crypto.VerifyLogicSig(account.Lsig, address) { //nolint:staticcheck // Preserve legacy signing validation behavior.
		return nil, errors.New("invalid logicsig signature")
	}

	return encodeSignedTxn(types.SignedTxn{Lsig: account.Lsig, Txn: tx}, address), nil
}

// pqSignedTxn returns the encoded SignedTxn carrying the given post-quantum
// signature bytes on behalf of the signer's account.
//
// The signature may be empty, which produces an envelope suitable for
// simulating transactions with the allowEmptySignatures option enabled.
func pqSignedTxn(signer crypto.PQSigner, tx types.Transaction, signature []byte) ([]byte, error) {
	salt, err := crypto.SaltForPQSigner(signer)
	if err != nil {
		return nil, err
	}

	publicKey, scheme := signer.PQPublicKey(), signer.PQScheme()
	stx := types.SignedTxn{
		Txn: tx,
		PQsig: types.PQSig{
			Scheme:    scheme,
			Salt:      salt,
			PublicKey: publicKey,
			Signature: signature,
		},
	}
	return encodeSignedTxn(stx, crypto.PQAddress(publicKey, scheme, salt)), nil
}

func pqSignTransaction(signer crypto.PQSigner, tx types.Transaction) ([]byte, error) {
	signature, err := signer.PQSign(transactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	return pqSignedTxn(signer, tx, signature)
}
