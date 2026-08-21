package transaction

import (
	"encoding/json"

	"github.com/algorand/go-algorand-sdk/v2/crypto"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// BasicAccountTransactionSigner that can sign transactions for the provided basic Account.
//
// Deprecated: having in-memory cryptographic secrets is discouraged, use
// Ed25519AccountTransactionSigner instead
type BasicAccountTransactionSigner struct {
	Account crypto.Account
}

// SignTransactions signs the provided transactions with the private key of the account.
func (txSigner BasicAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.SignTransactions(txGroup, indexesToSign)
}

// SignDelegationTo signs a delegation to the given LogicSig program. This
// program will have the authority to sign transactions on behalf of the signing
// account, called the delegating account.
func (txSigner BasicAccountTransactionSigner) SignDelegationTo(program []byte, args [][]byte) (lsa crypto.LogicSigAccount, err error) {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.SignDelegationTo(program, args)
}

// SignBytes signs the bytes and returns the signature
func (txSigner BasicAccountTransactionSigner) SignBytes(bytesToSign []byte) (signature []byte, err error) {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.SignBytes(bytesToSign)
}

// TealSign creates a signature compatible with ed25519verify opcode from
// contract address
func (txSigner BasicAccountTransactionSigner) TealSign(data []byte, contractAddress types.Address) (rawSig types.Signature, err error) {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.TealSign(data, contractAddress)
}

// AppendSignature appends the signature corresponding to the given signer,
// returning an encoded signed multisig transaction including the signature.
func (txSigner BasicAccountTransactionSigner) AppendSignature(ma crypto.MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error) {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.AppendSignature(ma, preStxBytes)
}

// AppendDelegationSignature adds an additional signature from a member of the
// delegating multisig account.
func (txSigner BasicAccountTransactionSigner) AppendDelegationSignature(lsa *crypto.LogicSigAccount) error {
	transactionSigner := Ed25519AccountTransactionSigner{Signer: txSigner.Account.AsSigner()}
	return transactionSigner.AppendDelegationSignature(lsa)
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner BasicAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(BasicAccountTransactionSigner); ok {
		otherJSON, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJSON, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJSON) == string(selfJSON)
	}
	return false
}

// MultiSigAccountTransactionSigner is a TransactionSigner that can
// sign transactions for the provided MultiSig Account
//
// Deprecated: having in-memory cryptographic secrets is discouraged, use
// MultiSigEd25519AccountTransactionSigner instead
type MultiSigAccountTransactionSigner struct {
	Msig crypto.MultisigAccount
	Sks  [][]byte
}

func (txSigner MultiSigAccountTransactionSigner) asNewSigner() (MultiSigEd25519AccountTransactionSigner, error) {
	signers := make([]crypto.Ed25519Signer, len(txSigner.Sks))
	for i, sk := range txSigner.Sks {
		signer, err := crypto.SKToInMemorySigner(sk)
		if err != nil {
			return MultiSigEd25519AccountTransactionSigner{}, err
		}

		signers[i] = signer
	}
	return MultiSigEd25519AccountTransactionSigner{Msig: txSigner.Msig, Signers: signers}, nil
}

// SignTransactions signs the provided transactions with the private keys of the account.
func (txSigner MultiSigAccountTransactionSigner) SignTransactions(txGroup []types.Transaction, indexesToSign []int) ([][]byte, error) {
	transactionSigner, err := txSigner.asNewSigner()
	if err != nil {
		return nil, err
	}
	return transactionSigner.SignTransactions(txGroup, indexesToSign)
}

// SignDelegationTo signs a delegation to the given LogicSig program. This
// program will have the authority to sign transactions on behalf of the signing
// account, called the delegating account.
func (txSigner MultiSigAccountTransactionSigner) SignDelegationTo(program []byte, args [][]byte) (lsa crypto.LogicSigAccount, err error) {
	transactionSigner, err := txSigner.asNewSigner()
	if err != nil {
		return crypto.LogicSigAccount{}, err
	}
	return transactionSigner.SignDelegationTo(program, args)
}

// Equals returns true if the other TransactionSigner equals this one.
func (txSigner MultiSigAccountTransactionSigner) Equals(other TransactionSigner) bool {
	if castedSigner, ok := other.(MultiSigAccountTransactionSigner); ok {
		otherJSON, err := json.Marshal(castedSigner)
		if err != nil {
			return false
		}

		selfJSON, err := json.Marshal(txSigner)
		if err != nil {
			return false
		}

		return string(otherJSON) == string(selfJSON)
	}
	return false
}
