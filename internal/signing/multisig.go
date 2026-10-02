package signing

import (
	"bytes"
	"crypto/sha512"
	"errors"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// prefix for multisig transaction signing
const msigAddrPrefix = "MultisigAddr"

// Multisig errors, re-exported under their original unexported names by the
// crypto package.
var (
	ErrMsigUnknownVersion        = errors.New("unknown version != 1")
	ErrMsigInvalidThreshold      = errors.New("invalid threshold")
	ErrMsigInvalidSecretKey      = errors.New("secret key has no corresponding public identity in multisig preimage")
	ErrMsigMergeLessThanTwo      = errors.New("cannot merge fewer than two multisig transactions")
	ErrMsigMergeKeysMismatch     = errors.New("multisig parameters do not match")
	ErrMsigMergeInvalidDups      = errors.New("mismatched duplicate signatures")
	ErrMsigMergeAuthAddrMismatch = errors.New("mismatched AuthAddrs")
)

// MultisigAccount holds multisig preimage data.
//
// Its fields are identical to crypto.MultisigAccount so that the two can be
// converted into one another.
type MultisigAccount struct {
	Version   uint8
	Threshold uint8
	Pks       []ed25519.PublicKey
}

// MultisigAccountFromSig creates a MultisigAccount from the preimage carried by
// a MultisigSig.
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
		err = ErrMsigUnknownVersion
		return
	}
	if ma.Threshold == 0 || len(ma.Pks) == 0 || int(ma.Threshold) > len(ma.Pks) {
		err = ErrMsigInvalidThreshold
		return
	}
	return
}

// ed25519MultisigSig returns a MultisigSig for the multisig account ma, with
// every subsig keyed but only the one belonging to sgnr populated by signing
// toBeSigned, alongside the index of that subsig. It returns an error if sgnr
// is not a member of ma.
func ed25519MultisigSig(sgnr Ed25519Signer, ma MultisigAccount, toBeSigned []byte) (msig types.MultisigSig, myIndex int, err error) {
	// check that sgnr.pk exists in the list of public keys in MultisigAccount ma
	myIndex = len(ma.Pks)
	myPublicKey := sgnr.Ed25519PublicKey()
	for i := 0; i < len(ma.Pks); i++ {
		if bytes.Equal(myPublicKey[:], ma.Pks[i]) {
			myIndex = i
		}
	}
	if myIndex == len(ma.Pks) {
		err = ErrMsigInvalidSecretKey
		return
	}

	// now, create the signed transaction
	msig.Version = ma.Version
	msig.Threshold = ma.Threshold
	msig.Subsigs = make([]types.MultisigSubsig, len(ma.Pks))
	for i := 0; i < len(ma.Pks); i++ {
		c := make([]byte, len(ma.Pks[i]))
		copy(c, ma.Pks[i])
		msig.Subsigs[i].Key = c
	}
	rawSig, err := Ed25519RawSignature(sgnr, toBeSigned)
	if err != nil {
		return
	}
	msig.Subsigs[myIndex].Sig = rawSig
	return
}

// Ed25519SignMultisigTransaction signs tx on behalf of the multisig account ma,
// returning its txid and the encoded SignedTxn with only sgnr's subsig
// populated, ready to be passed to other multisig signers to sign or broadcast.
func Ed25519SignMultisigTransaction(sgnr Ed25519Signer, ma MultisigAccount, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	maAddress, err := ma.Address()
	if err != nil {
		return
	}

	toBeSigned := TransactionBytesToSign(tx)
	msig, _, err := ed25519MultisigSig(sgnr, ma, toBeSigned)
	if err != nil {
		return
	}

	stxBytes = EncodeSignedTxn(types.SignedTxn{Msig: msig, Txn: tx}, maAddress)
	return TxIDFromBytesToSign(toBeSigned), stxBytes, nil
}

// Ed25519AppendMultisigTransaction appends the signature of sgnr to the encoded
// multisig SignedTxn preStxBytes, returning an encoded signed multisig
// transaction including the signature.
func Ed25519AppendMultisigTransaction(sgnr Ed25519Signer, ma MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error) {
	preStx := types.SignedTxn{}
	err = msgpack.Decode(preStxBytes, &preStx)
	if err != nil {
		return
	}
	_, partStxBytes, err := Ed25519SignMultisigTransaction(sgnr, ma, preStx.Txn)
	if err != nil {
		return
	}
	return MergeMultisigTransactions(partStxBytes, preStxBytes)
}

// MergeMultisigTransactions merges the given (partially) signed multisig transactions, and
// returns an encoded signed multisig transaction with the component signatures.
func MergeMultisigTransactions(stxsBytes ...[]byte) (txid string, stxBytes []byte, err error) {
	if len(stxsBytes) < 2 {
		err = ErrMsigMergeLessThanTwo
		return
	}
	var sig types.MultisigSig
	var refAddr *types.Address
	var refTx types.Transaction
	var refAuthAddr types.Address
	for _, partStxBytes := range stxsBytes {
		partStx := types.SignedTxn{}
		err = msgpack.Decode(partStxBytes, &partStx)
		if err != nil {
			return
		}
		// check that multisig parameters match
		partMa, innerErr := MultisigAccountFromSig(partStx.Msig)
		if innerErr != nil {
			err = innerErr
			return
		}
		partAddr, innerErr := partMa.Address()
		if innerErr != nil {
			err = innerErr
			return
		}
		if refAddr == nil {
			refAddr = &partAddr
			// add parameters to new merged txn
			sig.Version = partStx.Msig.Version
			sig.Threshold = partStx.Msig.Threshold
			sig.Subsigs = make([]types.MultisigSubsig, len(partStx.Msig.Subsigs))
			for i := 0; i < len(sig.Subsigs); i++ {
				c := make([]byte, len(partStx.Msig.Subsigs[i].Key))
				copy(c, partStx.Msig.Subsigs[i].Key)
				sig.Subsigs[i].Key = c
			}
			refTx = partStx.Txn
			refAuthAddr = partStx.AuthAddr
		}

		if partAddr != *refAddr {
			err = ErrMsigMergeKeysMismatch
			return
		}

		if partStx.AuthAddr != refAuthAddr {
			err = ErrMsigMergeAuthAddrMismatch
			return
		}

		// now, add subsignatures appropriately
		zeroSig := types.Signature{}
		for i := 0; i < len(sig.Subsigs); i++ {
			mSubsig := partStx.Msig.Subsigs[i]
			if mSubsig.Sig != zeroSig {
				if sig.Subsigs[i].Sig == zeroSig {
					sig.Subsigs[i].Sig = mSubsig.Sig
				} else if sig.Subsigs[i].Sig != mSubsig.Sig {
					err = ErrMsigMergeInvalidDups
					return
				}
			}
		}
	}
	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig:     sig,
		Txn:      refTx,
		AuthAddr: refAuthAddr,
	}
	stxBytes = msgpack.Encode(stx)
	// let's also compute the txid.
	txid = TxIDFromBytesToSign(TransactionBytesToSign(refTx))
	return
}
