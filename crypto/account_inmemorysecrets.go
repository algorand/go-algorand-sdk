package crypto

import (
	"errors"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// Account holds both the public and private information associated with an
// Algorand address
//
// Note: usage of in-memory cryptographic APIs is discouraged
type Account struct {
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
	Address    types.Address
}

func init() {
	addrLen := len(types.Address{})
	pkLen := ed25519.PublicKeySize
	if addrLen != pkLen {
		panic("address and public key are different sizes")
	}
}

// GenerateAccount generates a random Account
//
// Note: usage of in-memory cryptographic APIs is discouraged
func GenerateAccount() (kp Account) {
	// Generate an ed25519 keypair. This should never fail
	pk, sk, err := ed25519.GenerateKey(nil)
	if err != nil {
		panic(err)
	}

	// Convert the public key to an address
	var a types.Address
	n := copy(a[:], pk)
	if n != ed25519.PublicKeySize {
		panic("generated public key is the wrong size")
	}

	// Build the account
	kp.PublicKey = pk
	kp.PrivateKey = sk
	kp.Address = a
	return
}

// AccountFromPrivateKey derives the remaining Account fields from only a
// private key. The argument sk must have a length equal to
// ed25519.PrivateKeySize.
//
// Note: usage of in-memory cryptographic APIs is discouraged
func AccountFromPrivateKey(sk ed25519.PrivateKey) (account Account, err error) {
	if len(sk) != ed25519.PrivateKeySize {
		err = errInvalidPrivateKey
		return
	}

	// copy sk
	account.PrivateKey = make(ed25519.PrivateKey, len(sk))
	copy(account.PrivateKey, sk)

	account.PublicKey = sk.Public().(ed25519.PublicKey)
	if len(account.PublicKey) != ed25519.PublicKeySize {
		err = errors.New("generated public key is the wrong size")
		return
	}

	copy(account.Address[:], account.PublicKey)

	return
}

// AsSigner transforms this account to an Ed25519Signer using only the
// information provided by its PrivateKey.
//
// Panics if the private key is not well-formed.
func (acc Account) AsSigner() Ed25519Signer {
	return inMemoryEd25519Signer{
		sk: acc.PrivateKey,
		pk: Ed25519PublicKey(acc.PrivateKey.Public().(ed25519.PublicKey)),
	}
}

// MakeLogicSigAccountDelegated creates a new delegated LogicSigAccount. This
// type of LogicSig has the authority to sign transactions on behalf of another
// account, called the delegating account. If the delegating account is a
// multisig account, use MakeLogicSigAccountDelegated instead.
//
// The parameter signer is the private key of the delegating account.
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519MakeLogicSigAccountDelegated instead
func MakeLogicSigAccountDelegated(program []byte, args [][]byte, sk ed25519.PrivateKey) (lsa LogicSigAccount, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519MakeLogicSigAccountDelegated(program, args, sgnr)
}

// MakeLogicSigAccountDelegatedMsig creates a new delegated LogicSigAccount.
// This type of LogicSig has the authority to sign transactions on behalf of
// another account, called the delegating account. Use this function if the
// delegating account is a multisig account, otherwise use
// MakeLogicSigAccountDelegated.
//
// The parameter msigAccount is the delegating multisig account.
//
// The parameter signer is the private key of one of the members of the
// delegating multisig account. Use the method AppendMultisigSignature on the
// returned LogicSigAccount to add additional signatures from other members.
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519MakeLogicSigAccountDelegatedMsig instead
func MakeLogicSigAccountDelegatedMsig(program []byte, args [][]byte, msigAccount MultisigAccount, sk ed25519.PrivateKey) (lsa LogicSigAccount, err error) {
	sgnr, err := SKToInMemorySigner(sk)
	if err != nil {
		return
	}

	return Ed25519MakeLogicSigAccountDelegatedMsig(program, args, msigAccount, sgnr)
}

// AppendMultisigSignature adds an additional signature from a member of the
// delegating multisig account.
//
// The LogicSigAccount must represent a delegated LogicSig backed by a multisig
// account.
//
// Deprecated: usage of in-memory cryptographic APIs is discouraged, use
// Ed25519AppendMultisigSignature instead
func (lsa *LogicSigAccount) AppendMultisigSignature(sk ed25519.PrivateKey) error {
	signer, err := SKToInMemorySigner(sk)
	if err != nil {
		return err
	}
	return lsa.Ed25519AppendMultisigSignature(signer)
}
