package transaction

import (
	"bytes"
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

func transactionBytesToSign(tx types.Transaction) []byte {
	return append([]byte("TX"), msgpack.Encode(tx)...)
}

func ed25519SignTransaction(signer crypto.Ed25519Signer, tx types.Transaction) ([]byte, error) {
	signature, err := signer.Ed25519Sign(transactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}

	var sig types.Signature
	if copy(sig[:], signature) != len(sig) {
		return nil, errInvalidSignatureReturned
	}

	stx := types.SignedTxn{Sig: sig, Txn: tx}
	address := types.Address(signer.Ed25519PublicKey())
	if tx.Sender != address {
		stx.AuthAddr = address
	}
	return msgpack.Encode(stx), nil
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

	signature, err := signer.Ed25519Sign(transactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}
	if copy(msig.Subsigs[signerIndex].Sig[:], signature) != len(msig.Subsigs[signerIndex].Sig) {
		return nil, errInvalidSignatureReturned
	}

	stx := types.SignedTxn{Msig: msig, Txn: tx}
	address, err := account.Address()
	if err != nil {
		return nil, err
	}
	if tx.Sender != address {
		stx.AuthAddr = address
	}
	return msgpack.Encode(stx), nil
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

	stx := types.SignedTxn{Lsig: account.Lsig, Txn: tx}
	if tx.Sender != address {
		stx.AuthAddr = address
	}
	return msgpack.Encode(stx), nil
}

func signPQAccountTransaction(signer crypto.PQSigner, tx types.Transaction) ([]byte, error) {
	signature, err := signer.PQSign(transactionBytesToSign(tx))
	if err != nil {
		return nil, err
	}
	salt, err := crypto.SaltForPQSigner(signer)
	if err != nil {
		return nil, err
	}

	stx := types.SignedTxn{
		Txn: tx,
		PQsig: types.PQSig{
			Scheme:    signer.PQScheme(),
			Salt:      salt,
			PublicKey: signer.PQPublicKey(),
			Signature: signature,
		},
	}
	address := crypto.PQAddress(signer.PQPublicKey(), signer.PQScheme(), salt)
	if tx.Sender != address {
		stx.AuthAddr = address
	}
	return msgpack.Encode(stx), nil
}
