package crypto

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

// prefix for multisig transaction signing
const msigAddrPrefix = "MultisigAddr"

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
func MultisigAccountFromSig(sig types.MultisigSig) (ma MultisigAccount, err error) {
	ma.Version = sig.Version
	ma.Threshold = sig.Threshold
	ma.Pks = make([]ed25519.PublicKey, len(sig.Subsigs))
	for i := 0; i < len(sig.Subsigs); i++ {
		c := make([]byte, len(sig.Subsigs[i].Key))
		copy(c, sig.Subsigs[i].Key)
		ma.Pks[i] = c
	}
	err = ma.Validate()
	return
}

// Address takes this multisig preimage data, and generates the corresponding identifying
// address, committing to the exact group, version, and public keys that it requires to sign.
// Hash("MultisigAddr" || version uint8 || threshold uint8 || PK1 || PK2 || ...)
func (ma MultisigAccount) Address() (addr types.Address, err error) {
	// See go-algorand/crypto/multisig.go
	err = ma.Validate()
	if err != nil {
		return
	}
	buffer := append([]byte(msigAddrPrefix), byte(ma.Version), byte(ma.Threshold))
	for _, pki := range ma.Pks {
		buffer = append(buffer, pki[:]...)
	}
	return sha512.Sum512_256(buffer), nil
}

// Validate ensures that this multisig setup is a valid multisig account
func (ma MultisigAccount) Validate() (err error) {
	if ma.Version != 1 {
		err = errMsigUnknownVersion
		return
	}
	if ma.Threshold == 0 || len(ma.Pks) == 0 || int(ma.Threshold) > len(ma.Pks) {
		err = errMsigInvalidThreshold
		return
	}
	return
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
	toBeSigned := programToSign(lsig.Logic)
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

// Ed25519MakeLogicSigAccountDelegated creates a new delegated LogicSigAccount.
// This type of LogicSig has the authority to sign transactions on behalf of
// another account, called the delegating account. If the delegating account is
// a multisig account, use Ed25519MakeLogicSigAccountDelegatedMsig instead.
//
// The parameter signer is an Ed25519Signer for the delegating account.
func Ed25519MakeLogicSigAccountDelegated(program []byte, args [][]byte, signer Ed25519Signer) (lsa LogicSigAccount, err error) {
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

// Ed25519MakeLogicSigAccountDelegatedMsig creates a new delegated
// LogicSigAccount.  This type of LogicSig has the authority to sign
// transactions on behalf of another account, called the delegating account. Use
// this function if the delegating account is a multisig account, otherwise use
// Ed25519MakeLogicSigAccountDelegated.
//
// The parameter msigAccount is the delegating multisig account.
//
// The parameter signer is an Ed25519Signer of one of the members of the
// delegating multisig account. Use the method Ed25519AppendMultisigSignature on
// the returned LogicSigAccount to add additional signatures from other members.
func Ed25519MakeLogicSigAccountDelegatedMsig(program []byte, args [][]byte, msigAccount MultisigAccount, signer Ed25519Signer) (lsa LogicSigAccount, err error) {
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

// Ed25519AppendMultisigSignature adds an additional signature from a member of the
// delegating multisig account.
//
// The LogicSigAccount must represent a delegated LogicSig backed by a multisig
// account.
func (lsa *LogicSigAccount) Ed25519AppendMultisigSignature(signer Ed25519Signer) error {
	return Ed25519AppendMultisigToLogicSig(&lsa.Lsig, signer)
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
	hasSig, _, _, _, count := lsig.SignatureCount()

	if count > 1 {
		err = errLsigTooManySignatures
		return
	}

	if hasSig {
		if signerPublicKey == nil {
			err = errLsigNoPublicKey
			return
		}

		toBeSigned := programToSign(lsig.Logic)
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
	hasSig := lsa.Lsig.Sig != (types.Signature{})
	hasMsig := !lsa.Lsig.Msig.Blank()
	hasLMsig := !lsa.Lsig.LMsig.Blank()
	hasPQsig := !lsa.Lsig.PQsig.Blank()
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
	var count int
	if hasSig, hasMsig, hasLMsig, hasPQsig, count = lsa.Lsig.SignatureCount(); count > 1 {
		err = errLsigTooManySignatures
	}
	return
}
