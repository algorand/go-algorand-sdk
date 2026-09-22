package transaction

import (
	"errors"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

var errNoMultisigSigners = errors.New("multisig signer has no signing keys")

// TransactionSigner represents a function which can sign transactions from an atomic transaction group.
// @param txnGroup - The atomic group containing transactions to be signed
// @param indexesToSign - An array of indexes in the atomic transaction group that should be signed
// @returns An array of encoded signed transactions. The length of the
//
//	array will be the same as the length of indexesToSign, and each index i in the array
//	corresponds to the signed transaction from txnGroup[indexesToSign[i]]
type TransactionSigner interface { //nolint:revive // Ignore stuttering for backwards compatibility
	SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error)
	Equals(other TransactionSigner) bool
}

// DelegatableSigner represents a signer that can delegate its authority to a LogicSig program.
type DelegatableSigner interface {
	// SignDelegationTo signs a delegation to the given LogicSig program. This
	// program will have the authority to sign transactions on behalf of the
	// signing account, called the delegating account.
	SignDelegationTo(program []byte, args [][]byte) (lsa crypto.LogicSigAccount, err error)
}

// Ed25519TransactionSigner represents a signer that can perform operations
// exclusive to elliptic curve accounts (like signing ed25519 multisig
// accounts).
type Ed25519TransactionSigner interface {
	// SignBytes signs the bytes and returns the signature
	SignBytes(bytesToSign []byte) (signature []byte, err error)
	// TealSign creates a signature compatible with ed25519verify opcode from
	// contract address
	TealSign(data []byte, contractAddress types.Address) (rawSig types.Signature, err error)
	// AppendSignature appends the signature corresponding to the given signer,
	// returning an encoded signed multisig transaction including the signature.
	AppendSignature(ma crypto.MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error)
	// AppendDelegationSignature adds an additional signature from a member of
	// the delegating multisig account.
	AppendDelegationSignature(lsa *crypto.LogicSigAccount) error
}

// Ed25519AccountTransactionSigner is a TransactionSigner that can sign
// transactions using the provided Ed25519 signer.
type Ed25519AccountTransactionSigner struct {
	Signer crypto.Ed25519Signer
}

// SignTransactions signs the provided transactions with the Ed25519Signer.
func (txSigner Ed25519AccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	return signTransactions(txGroup, indexesToSign, func(tx types.Transaction) ([]byte, error) {
		return ed25519SignTransaction(txSigner.Signer, tx)
	})
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner Ed25519AccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(Ed25519AccountTransactionSigner); ok {
		return equalSignerImplementations(txSigner.Signer, castedSigner.Signer)
	}
	return false
}

// SignDelegationTo signs a delegation to the given LogicSig program. This
// program will have the authority to sign transactions on behalf of the signing
// account, called the delegating account.
func (txSigner Ed25519AccountTransactionSigner) SignDelegationTo(program []byte, args [][]byte) (lsa crypto.LogicSigAccount, err error) {
	return crypto.Ed25519MakeLogicSigAccountDelegated(program, args, txSigner.Signer)
}

// SignBytes signs the bytes and returns the signature
func (txSigner Ed25519AccountTransactionSigner) SignBytes(bytesToSign []byte) (signature []byte, err error) {
	return crypto.Ed25519SignBytes(txSigner.Signer, bytesToSign)
}

// TealSign creates a signature compatible with ed25519verify opcode from
// contract address
func (txSigner Ed25519AccountTransactionSigner) TealSign(data []byte, contractAddress types.Address) (rawSig types.Signature, err error) {
	return crypto.Ed25519TealSign(txSigner.Signer, data, contractAddress)
}

// AppendSignature appends the signature corresponding to the given signer,
// returning an encoded signed multisig transaction including the signature.
func (txSigner Ed25519AccountTransactionSigner) AppendSignature(ma crypto.MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error) {
	return ed25519AppendMultisigTransaction(txSigner.Signer, ma, preStxBytes)
}

// AppendDelegationSignature adds an additional signature from a member of the
// delegating multisig account.
func (txSigner Ed25519AccountTransactionSigner) AppendDelegationSignature(lsa *crypto.LogicSigAccount) error {
	return lsa.Ed25519AppendMultisigSignature(txSigner.Signer)
}

// MultiSigEd25519AccountTransactionSigner is a TransactionSigner that can sign
// transactions for the provided MultiSig Account
type MultiSigEd25519AccountTransactionSigner struct {
	Msig    crypto.MultisigAccount
	Signers []crypto.Ed25519Signer
}

// SignTransactions signs the provided transactions with the Ed25519Signer.
func (txSigner MultiSigEd25519AccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	if len(txSigner.Signers) == 0 {
		return nil, errNoMultisigSigners
	}

	return signTransactions(txGroup, indexesToSign, func(tx types.Transaction) ([]byte, error) {
		unmergedStxs := make([][]byte, len(txSigner.Signers))
		for i, sgnr := range txSigner.Signers {
			stxBytes, err := ed25519SignMultisigTransaction(sgnr, txSigner.Msig, tx)
			if err != nil {
				return nil, err
			}

			unmergedStxs[i] = stxBytes
		}

		if len(unmergedStxs) == 1 {
			return unmergedStxs[0], nil
		}

		_, stxBytes, err := crypto.MergeMultisigTransactions(unmergedStxs...)
		return stxBytes, err
	})
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner MultiSigEd25519AccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(MultiSigEd25519AccountTransactionSigner); ok {
		if !equalBySerialization(txSigner.Msig, castedSigner.Msig) {
			return false
		}

		if len(txSigner.Signers) != len(castedSigner.Signers) {
			return false
		}

		for idx, sgnr := range txSigner.Signers {
			if !equalSignerImplementations(sgnr, castedSigner.Signers[idx]) {
				return false
			}
		}

		return true
	}
	return false
}

// SignDelegationTo signs a delegation to the given LogicSig program. This
// program will have the authority to sign transactions on behalf of the signing
// account, called the delegating account.
func (txSigner MultiSigEd25519AccountTransactionSigner) SignDelegationTo(program []byte, args [][]byte) (lsa crypto.LogicSigAccount, err error) {
	if len(txSigner.Signers) == 0 {
		return crypto.LogicSigAccount{}, errNoMultisigSigners
	}
	firstSigner := txSigner.Signers[0]
	lsa, err = crypto.Ed25519MakeLogicSigAccountDelegatedMsig(program, args, txSigner.Msig, firstSigner)
	if err != nil {
		return
	}
	for _, signer := range txSigner.Signers[1:] {
		err = lsa.Ed25519AppendMultisigSignature(signer)
		if err != nil {
			return crypto.LogicSigAccount{}, err
		}
	}
	return
}

// LogicSigAccountTransactionSigner is a TransactionSigner that can
// sign transactions for the provided LogicSigAccount.
type LogicSigAccountTransactionSigner struct {
	LogicSigAccount crypto.LogicSigAccount
}

// SignTransactions signs the provided transactions with the private key of the account.
func (txSigner LogicSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	return signTransactions(txGroup, indexesToSign, func(tx types.Transaction) ([]byte, error) {
		return signLogicSigAccountTransaction(txSigner.LogicSigAccount, tx)
	})
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner LogicSigAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(LogicSigAccountTransactionSigner); ok {
		return equalBySerialization(txSigner, castedSigner)
	}
	return false
}

// PQAccountTransactionSigner is a TransactionSigner that can
// sign transactions using the provided PQSigner
type PQAccountTransactionSigner struct {
	Signer crypto.PQSigner
}

// SignTransactions signs the provided transactions with the PQSigner signer.
func (txSigner PQAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	return signTransactions(txGroup, indexesToSign, func(tx types.Transaction) ([]byte, error) {
		return pqSignTransaction(txSigner.Signer, tx)
	})
}

// SignDelegationTo signs a delegation to the given LogicSig program. This
// program will have the authority to sign transactions on behalf of the signing
// account, called the delegating account.
func (txSigner PQAccountTransactionSigner) SignDelegationTo(program []byte, args [][]byte) (lsa crypto.LogicSigAccount, err error) {
	return crypto.MakeLogicSigAccountDelegatedPQ(program, args, txSigner.Signer)
}

// equalPQSigners reports whether two PQ signer interface values refer to the
// same signer implementation.
func equalPQSigners(signer, other crypto.PQSigner) bool {
	return equalSignerImplementations(signer, other)
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner PQAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(PQAccountTransactionSigner); ok {
		return equalPQSigners(txSigner.Signer, castedSigner.Signer)
	}
	return false
}

// EmptyTransactionSigner is a TransactionSigner that produces signed transaction objects without
// signatures. This is useful for simulating transactions, but it won't work for actual submission.
type EmptyTransactionSigner struct{}

// SignTransactions returns SignedTxn bytes but does not sign them.
func (txSigner EmptyTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	return signTransactions(txGroup, indexesToSign, func(tx types.Transaction) ([]byte, error) {
		return msgpack.Encode(&types.SignedTxn{Txn: tx}), nil
	})
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner EmptyTransactionSigner) Equals(other TransactionSigner) bool {
	_, ok := other.(EmptyTransactionSigner)
	return ok
}

// PQEmptyTransactionSigner is a TransactionSigner that produces signed transaction
// objects with an empty post-quantum signature envelope (scheme, salt, and public key populated,
// but signature bytes empty). This is useful for simulating transactions with the
// allowEmptySignatures option enabled, allowing algod to charge the post-quantum fee surcharge
// without paying the computational cost of generating a real post-quantum signature.
type PQEmptyTransactionSigner struct {
	Signer crypto.PQSigner
}

// SignTransactions returns SignedTxn bytes with placeholder PQ signatures.
func (txSigner PQEmptyTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	return signTransactions(txGroup, indexesToSign, func(tx types.Transaction) ([]byte, error) {
		return pqSignedTxn(txSigner.Signer, tx, []byte{})
	})
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner PQEmptyTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(PQEmptyTransactionSigner); ok {
		return equalPQSigners(txSigner.Signer, castedSigner.Signer)
	}
	return false
}
