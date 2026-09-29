package crypto

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/internal/signing"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

/* Multisig Support */

// MultisigAccount is a convenience type for holding multisig preimage data
type MultisigAccount struct {
	// Version is the version of this multisig
	Version uint8
	// Threshold is how many signatures are needed to fully sign as this address
	Threshold uint8
	// Pks is an ordered list of public keys that could potentially sign a message
	Pks []ed25519.PublicKey
}

// MultisigAccountWithParams creates a MultisigAccount with the given parameters
func MultisigAccountWithParams(version uint8, threshold uint8, addrs []types.Address) (ma MultisigAccount, err error) {
	ma.Version = version
	ma.Threshold = threshold
	ma.Pks = make([]ed25519.PublicKey, len(addrs))
	for i := 0; i < len(addrs); i++ {
		ma.Pks[i] = addrs[i][:]
	}
	err = ma.Validate()
	return
}

// MultisigAccountFromSig is a convenience method that creates an account
// from a sig in a signed tx. Useful for getting addresses from signed msig txs, etc.
func MultisigAccountFromSig(sig types.MultisigSig) (MultisigAccount, error) {
	ma, err := signing.MultisigAccountFromSig(sig)
	return MultisigAccount(ma), err
}

// Address takes this multisig preimage data, and generates the corresponding identifying
// address, committing to the exact group, version, and public keys that it requires to sign.
// Hash("MultisigAddr" || version uint8 || threshold uint8 || PK1 || PK2 || ...)
func (ma MultisigAccount) Address() (addr types.Address, err error) {
	return signing.MultisigAccount(ma).Address()
}

// Validate ensures that this multisig setup is a valid multisig account
func (ma MultisigAccount) Validate() (err error) {
	return signing.MultisigAccount(ma).Validate()
}

// Blank return true if MultisigAccount is empty
// struct containing []ed25519.PublicKey cannot be compared
func (ma MultisigAccount) Blank() bool {
	if ma.Version != 0 {
		return false
	}
	if ma.Threshold != 0 {
		return false
	}
	if ma.Pks != nil {
		return false
	}
	return true
}

/* LogicSig support */

// LogicSigAddress returns the contract (escrow) address for a LogicSig.
//
// NOTE: If the LogicSig is delegated to another account this will not
// return the delegated address of the LogicSig.
func LogicSigAddress(lsig types.LogicSig) types.Address {
	toBeSigned := signing.ProgramToSign(lsig.Logic)
	checksum := sha512.Sum512_256(toBeSigned)

	var addr types.Address
	n := copy(addr[:], checksum[:])
	if n != ed25519.PublicKeySize {
		panic(fmt.Sprintf("Generated public key has length of %d, expected %d", n, ed25519.PublicKeySize))
	}
	return addr
}

// LogicSigAccount represents an account that can sign with a LogicSig program.
type LogicSigAccount struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// The underlying LogicSig object
	Lsig types.LogicSig `codec:"lsig"`

	// The key that provided Lsig.Sig, if any
	SigningKey ed25519.PublicKey `codec:"sigkey"`
}

// MakeLogicSigAccountEscrowChecked creates a new escrow LogicSigAccount.
// The address of this account will be a hash of its program.
func MakeLogicSigAccountEscrowChecked(program []byte, args [][]byte) (LogicSigAccount, error) {
	lsig, err := makeLogicSig(program, args, nil, MultisigAccount{})
	if err != nil {
		return LogicSigAccount{}, err
	}
	return LogicSigAccount{Lsig: lsig}, nil
}

// ed25519MakeLogicSigAccountDelegated backs the deprecated in-memory
// MakeLogicSigAccountDelegated. Delegation signing itself lives in the
// transaction package.
func ed25519MakeLogicSigAccountDelegated(program []byte, args [][]byte, signer Ed25519Signer) (lsa LogicSigAccount, err error) {
	var ma MultisigAccount
	lsig, err := makeLogicSig(program, args, signer, ma)
	if err != nil {
		return
	}

	pk := signer.Ed25519PublicKey()
	lsa = LogicSigAccount{
		Lsig: lsig,
		// attach SigningKey to remember which account the signature belongs to
		SigningKey: pk[:],
	}
	return
}

// ed25519MakeLogicSigAccountDelegatedMsig backs the deprecated in-memory
// MakeLogicSigAccountDelegatedMsig. Delegation signing itself lives in the
// transaction package.
func ed25519MakeLogicSigAccountDelegatedMsig(program []byte, args [][]byte, msigAccount MultisigAccount, signer Ed25519Signer) (lsa LogicSigAccount, err error) {
	lsig, err := makeLogicSig(program, args, signer, msigAccount)
	if err != nil {
		return
	}
	lsa = LogicSigAccount{
		Lsig: lsig,
		// do not attach SigningKey, since that doesn't apply to an msig signature
	}
	return
}

// LogicSigAccountFromLogicSig creates a LogicSigAccount from an existing
// LogicSig object.
//
// The parameter signerPublicKey must be present if the LogicSig is delegated
// and the delegating account is backed by a single private key (i.e. not a
// multisig account). In this case, signerPublicKey must be the public key of
// the delegating account. In all other cases, an error will be returned if
// signerPublicKey is present.
func LogicSigAccountFromLogicSig(lsig types.LogicSig, signerPublicKey *ed25519.PublicKey) (lsa LogicSigAccount, err error) {
	hasSig, _, _, _, err := lsigSignatures(lsig)
	if err != nil {
		return
	}

	if hasSig {
		if signerPublicKey == nil {
			err = errLsigNoPublicKey
			return
		}

		toBeSigned := signing.ProgramToSign(lsig.Logic)
		valid := ed25519.Verify(*signerPublicKey, toBeSigned, lsig.Sig[:])
		if !valid {
			err = errLsigInvalidPublicKey
			return
		}

		lsa.Lsig = lsig
		lsa.SigningKey = make(ed25519.PublicKey, len(*signerPublicKey))
		copy(lsa.SigningKey, *signerPublicKey)
		return
	}

	if signerPublicKey != nil {
		err = errLsigAccountPublicKeyNotNeeded
		return
	}

	lsa.Lsig = lsig
	return
}

// IsDelegated returns true if and only if the LogicSig has been delegated to
// another account with a signature.
//
// Note this function only checks for the presence of a delegation signature. To
// verify the delegation signature, use VerifyLogicSig.
func (lsa LogicSigAccount) IsDelegated() bool {
	hasSig, hasMsig, hasLMsig, hasPQsig, _ := lsigSignatures(lsa.Lsig)
	return hasSig || hasMsig || hasLMsig || hasPQsig
}

// Address returns the address of this LogicSigAccount.
//
// If the LogicSig is delegated to another account, this will return the address
// of that account.
//
// If the LogicSig is not delegated to another account, this will return an
// escrow address that is the hash of the LogicSig's program code.
func (lsa LogicSigAccount) Address() (addr types.Address, err error) {
	hasSig, hasMsig, hasLMsig, hasPQsig, err := lsa.hasSignatures()
	if err != nil {
		return types.Address{}, err
	}

	if hasSig {
		n := copy(addr[:], lsa.SigningKey)
		if n != ed25519.PublicKeySize {
			err = fmt.Errorf("generated public key has length of %d, expected %d", n, ed25519.PublicKeySize)
		}
		return
	}

	if hasMsig {
		var msigAccount MultisigAccount
		msigAccount, err = MultisigAccountFromSig(lsa.Lsig.Msig)
		if err != nil {
			return
		}
		addr, err = msigAccount.Address()
		return
	}

	if hasLMsig {
		var msigAccount MultisigAccount
		msigAccount, err = MultisigAccountFromSig(lsa.Lsig.LMsig)
		if err != nil {
			return
		}
		addr, err = msigAccount.Address()
		return
	}

	if hasPQsig {
		addr = PQAddressFromSig(lsa.Lsig.PQsig)
		return
	}

	addr = LogicSigAddress(lsa.Lsig)
	return
}

func (lsa LogicSigAccount) hasSignatures() (hasSig, hasMsig, hasLMsig, hasPQsig bool, err error) {
	return lsigSignatures(lsa.Lsig)
}
