package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base32"
	"fmt"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/logic"
	"github.com/algorand/go-algorand-sdk/types"
)

// txidPrefix is prepended to a transaction when computing its txid
var txidPrefix = []byte("TX")

// tgidPrefix is prepended to a transaction group when computing the group ID
var tgidPrefix = []byte("TG")

// bidPrefix is prepended to a bid when signing it
var bidPrefix = []byte("aB")

// bytesPrefix is prepended to a message when signing
var bytesPrefix = []byte("MX")

// programPrefix is prepended to a logic program when computing a hash
var programPrefix = []byte("Program")

// RandomBytes fills the passed slice with randomness, and panics if it is
// unable to do so
func RandomBytes(s []byte) {
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
}

// SignTransaction accepts a private key and a transaction, and returns the
// bytes of a signed transaction ready to be broadcasted to the network
func SignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	s, txid, err := rawSignTransaction(sk, tx)
	if err != nil {
		return
	}
	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

// rawTransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from.
func rawTransactionBytesToSign(tx types.Transaction) []byte {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{txidPrefix, encodedTx}
	return bytes.Join(msgParts, nil)
}

// txID computes a transaction id from raw transaction bytes
func txIDFromRawTxnBytesToSign(toBeSigned []byte) (txid string) {
	txidBytes := sha512.Sum512_256(toBeSigned)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}

// txIDFromTransaction is a convenience function for generating txID from txn
func txIDFromTransaction(tx types.Transaction) (txid string) {
	txid = txIDFromRawTxnBytesToSign(rawTransactionBytesToSign(tx))
	return
}

// rawSignTransaction signs the msgpack-encoded tx (with prepended "TX" prefix), and returns the sig and txid
func rawSignTransaction(sk ed25519.PrivateKey, tx types.Transaction) (s types.Signature, txid string, err error) {
	toBeSigned := rawTransactionBytesToSign(tx)

	// Sign the encoded transaction
	signature := ed25519.Sign(sk, toBeSigned)

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	n := copy(s[:], signature)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}
	// Populate txID
	txid = txIDFromRawTxnBytesToSign(toBeSigned)
	return
}

// SignBytes signs the bytes and returns the signature
func SignBytes(sk ed25519.PrivateKey, bytesToSign []byte) (signature []byte, err error) {
	// prepend the prefix for signing bytes
	toBeSigned := bytes.Join([][]byte{bytesPrefix, bytesToSign}, nil)

	// sign the bytes
	signature = ed25519.Sign(sk, toBeSigned)
	return
}

//VerifyBytes verifies that the signature is valid
func VerifyBytes(pk ed25519.PublicKey, message, signature []byte) bool {
	msgParts := [][]byte{bytesPrefix, message}
	toBeVerified := bytes.Join(msgParts, nil)
	return ed25519.Verify(pk, toBeVerified, signature)
}

// SignBid accepts a private key and a bid, and returns the signature of the
// bid under that key
func SignBid(sk ed25519.PrivateKey, bid types.Bid) (signedBid []byte, err error) {
	// Encode the bid as msgpack
	encodedBid := msgpack.Encode(bid)

	// Prepend the hashable prefix
	msgParts := [][]byte{bidPrefix, encodedBid}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded bid
	sig := ed25519.Sign(sk, toBeSigned)

	var s types.Signature
	n := copy(s[:], sig)
	if n != len(s) {
		err = errInvalidSignatureReturned
		return
	}

	sb := types.SignedBid{
		Bid: bid,
		Sig: s,
	}

	nf := types.NoteField{
		Type:      types.NoteBid,
		SignedBid: sb,
	}

	signedBid = msgpack.Encode(nf)
	return
}

/* Multisig Support */

type signer func() (signature types.Signature, err error)

// Service function to make a single signature in Multisig
func multisigSingle(sk ed25519.PrivateKey, ma MultisigAccount, customSigner signer) (msig types.MultisigSig, myIndex int, err error) {
	// check that sk.pk exists in the list of public keys in MultisigAccount ma
	myIndex = len(ma.Pks)
	myPublicKey := sk.Public().(ed25519.PublicKey)
	for i := 0; i < len(ma.Pks); i++ {
		if bytes.Equal(myPublicKey, ma.Pks[i]) {
			myIndex = i
		}
	}
	if myIndex == len(ma.Pks) {
		err = errMsigInvalidSecretKey
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
	rawSig, err := customSigner()
	if err != nil {
		return
	}
	msig.Subsigs[myIndex].Sig = rawSig
	return
}

// SignMultisigTransaction signs the given transaction, and multisig preimage, with the
// private key, returning the bytes of a signed transaction with the multisig field
// partially populated, ready to be passed to other multisig signers to sign or broadcast.
func SignMultisigTransaction(sk ed25519.PrivateKey, ma MultisigAccount, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	err = ma.Validate()
	if err != nil {
		return
	}
	// check that the address of txn matches the preimage
	maAddress, err := ma.Address()
	if err != nil {
		return
	}
	if tx.Sender != maAddress { // array value comparison is fine
		err = errMsigBadTxnSender
		return
	}

	// this signer signs a transaction and sets txid from the closure
	customSigner := func() (rawSig types.Signature, err error) {
		rawSig, txid, err = rawSignTransaction(sk, tx)
		return rawSig, err
	}

	sig, _, err := multisigSingle(sk, ma, customSigner)
	if err != nil {
		return
	}

	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig: sig,
		Txn:  tx,
	}
	stxBytes = msgpack.Encode(stx)
	return
}

// MergeMultisigTransactions merges the given (partially) signed multisig transactions, and
// returns an encoded signed multisig transaction with the component signatures.
func MergeMultisigTransactions(stxsBytes ...[]byte) (txid string, stxBytes []byte, err error) {
	if len(stxsBytes) < 2 {
		err = errMsigMergeLessThanTwo
		return
	}
	var sig types.MultisigSig
	var refAddr *types.Address
	var refTx types.Transaction
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
		} else {
			if partAddr != *refAddr {
				err = errMsigMergeKeysMismatch
				return
			}
		}
		// now, add subsignatures appropriately
		zeroSig := types.Signature{}
		for i := 0; i < len(sig.Subsigs); i++ {
			mSubsig := partStx.Msig.Subsigs[i]
			if mSubsig.Sig != zeroSig {
				if sig.Subsigs[i].Sig == zeroSig {
					sig.Subsigs[i].Sig = mSubsig.Sig
				} else if sig.Subsigs[i].Sig != mSubsig.Sig {
					err = errMsigMergeInvalidDups
					return
				}
			}
		}
	}
	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig: sig,
		Txn:  refTx,
	}
	stxBytes = msgpack.Encode(stx)
	// let's also compute the txid.
	txid = txIDFromTransaction(refTx)
	return
}

// AppendMultisigTransaction appends the signature corresponding to the given private key,
// returning an encoded signed multisig transaction including the signature.
// While we could compute the multisig preimage from the multisig blob, we ask the caller
// to pass it back in, to explicitly check that they know who they are signing as.
func AppendMultisigTransaction(sk ed25519.PrivateKey, ma MultisigAccount, preStxBytes []byte) (txid string, stxBytes []byte, err error) {
	preStx := types.SignedTxn{}
	err = msgpack.Decode(preStxBytes, &preStx)
	if err != nil {
		return
	}
	_, partStxBytes, err := SignMultisigTransaction(sk, ma, preStx.Txn)
	if err != nil {
		return
	}
	txid, stxBytes, err = MergeMultisigTransactions(partStxBytes, preStxBytes)
	return
}

// VerifyMultisig verifies an assembled MultisigSig
func VerifyMultisig(addr types.Address, message []byte, msig types.MultisigSig) bool {
	msigAccount, err := MultisigAccountFromSig(msig)
	if err != nil {
		return false
	}

	if msigAddress, err := msigAccount.Address(); err != nil || msigAddress != addr {
		return false
	}

	// check that we don't have too many multisig subsigs
	if len(msig.Subsigs) > 255 {
		return false
	}

	// check that we don't have too few multisig subsigs
	if len(msig.Subsigs) < int(msig.Threshold) {
		return false
	}

	// checks the number of non-blank signatures is no less than threshold
	var counter int
	for _, subsigi := range msig.Subsigs {
		if (subsigi.Sig != types.Signature{}) {
			counter++
		}
	}
	if counter < int(msig.Threshold) {
		return false
	}

	// checks individual signature verifies
	var verifiedCount uint8
	for _, subsigi := range msig.Subsigs {
		if (subsigi.Sig != types.Signature{}) {
			if !ed25519.Verify(subsigi.Key, message, subsigi.Sig[:]) {
				return false
			}
			verifiedCount++
		}
	}

	if verifiedCount < msig.Threshold {
		return false
	}

	return true
}

// ComputeGroupID returns group ID for a group of transactions
func ComputeGroupID(txgroup []types.Transaction) (gid types.Digest, err error) {
	// MaxTxGroupSize is max number of transactions in a single group
	const MaxTxGroupSize = 16
	if len(txgroup) > MaxTxGroupSize {
		err = fmt.Errorf("txgroup too large, %v > max size %v", len(txgroup), MaxTxGroupSize)
		return
	}
	var group types.TxGroup
	empty := types.Digest{}
	for _, tx := range txgroup {
		if tx.Group != empty {
			err = fmt.Errorf("transaction %v already has a group %v", tx, tx.Group)
			return
		}

		txID := sha512.Sum512_256(rawTransactionBytesToSign(tx))
		group.TxGroupHashes = append(group.TxGroupHashes, txID)
	}

	encoded := msgpack.Encode(group)

	// Prepend the hashable prefix and hash it
	msgParts := [][]byte{tgidPrefix, encoded}
	return sha512.Sum512_256(bytes.Join(msgParts, nil)), nil
}

/* LogicSig support */

// VerifyLogicSig verifies LogicSig against assumed sender address
func VerifyLogicSig(lsig types.LogicSig, sender types.Address) (result bool) {
	if err := logic.CheckProgram(lsig.Logic, lsig.Args); err != nil {
		return false
	}

	hasSig := lsig.Sig != (types.Signature{})
	hasMsig := !lsig.Msig.Blank()

	// require only one or zero sig
	if hasSig && hasMsig {
		return false
	}

	result = false
	toBeSigned := programToSign(lsig.Logic)
	// logic sig, compare hashes
	if !hasSig && !hasMsig {
		result = types.Digest(sha512.Sum512_256(toBeSigned)) == types.Digest(sender)
		return
	}

	if hasSig {
		result = ed25519.Verify(sender[:], toBeSigned, lsig.Sig[:])
		return
	}

	result = VerifyMultisig(sender, toBeSigned, lsig.Msig)
	return
}

// SignLogicsigTransaction takes LogicSig object and a transaction and returns the
// bytes of a signed transaction ready to be broadcasted to the network
// Note, LogicSig actually can be attached to any transaction (with matching sender field for Sig and Multisig cases)
// and it is a program's responsibility to approve/decline the transaction
func SignLogicsigTransaction(lsig types.LogicSig, tx types.Transaction) (txid string, stxBytes []byte, err error) {

	if !VerifyLogicSig(lsig, tx.Header.Sender) {
		err = errLsigInvalidSignature
		return
	}

	txid = txIDFromTransaction(tx)
	// Construct the SignedTxn
	stx := types.SignedTxn{
		Lsig: lsig,
		Txn:  tx,
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

func programToSign(program []byte) []byte {
	parts := [][]byte{programPrefix, program}
	toBeSigned := bytes.Join(parts, nil)
	return toBeSigned
}

func signProgram(sk ed25519.PrivateKey, program []byte) (sig types.Signature, err error) {
	toBeSigned := programToSign(program)
	rawSig := ed25519.Sign(sk, toBeSigned)
	n := copy(sig[:], rawSig)
	if n != len(sig) {
		err = errInvalidSignatureReturned
		return
	}
	return
}

func AddressFromProgram(program []byte) types.Address {
	toBeHashed := programToSign(program)
	hash := sha512.Sum512_256(toBeHashed)
	return types.Address(hash)
}

// MakeLogicSig produces a new LogicSig signature.
// The function can work in three modes:
// 1. If no sk and ma provided then it returns contract-only LogicSig
// 2. If no ma provides, it returns Sig delegated LogicSig
// 3. If both sk and ma specified the function returns Multisig delegated LogicSig
func MakeLogicSig(program []byte, args [][]byte, sk ed25519.PrivateKey, ma MultisigAccount) (lsig types.LogicSig, err error) {
	if len(program) == 0 {
		err = errLsigInvalidProgram
		return
	}
	if err = logic.CheckProgram(program, args); err != nil {
		return
	}

	if sk == nil && ma.Blank() {
		lsig.Logic = program
		lsig.Args = args
		return
	}

	if ma.Blank() {
		var sig types.Signature
		sig, err = signProgram(sk, program)
		if err != nil {
			return
		}

		lsig.Logic = program
		lsig.Args = args
		lsig.Sig = types.Signature(sig)
		return
	}

	// Format Multisig
	err = ma.Validate()
	if err != nil {
		return
	}

	// this signer signs a program
	customSigner := func() (rawSig types.Signature, err error) {
		return signProgram(sk, program)
	}

	msig, _, err := multisigSingle(sk, ma, customSigner)
	if err != nil {
		return
	}

	lsig.Logic = program
	lsig.Args = args
	lsig.Msig = msig

	return
}

// AppendMultisigToLogicSig adds a new signature to multisigned LogicSig
func AppendMultisigToLogicSig(lsig *types.LogicSig, sk ed25519.PrivateKey) error {
	if lsig.Msig.Blank() {
		return errLsigEmptyMsig
	}

	ma, err := MultisigAccountFromSig(lsig.Msig)
	if err != nil {
		return err
	}

	customSigner := func() (rawSig types.Signature, err error) {
		return signProgram(sk, lsig.Logic)
	}

	msig, idx, err := multisigSingle(sk, ma, customSigner)
	if err != nil {
		return err
	}

	lsig.Msig.Subsigs[idx] = msig.Subsigs[idx]

	return nil
}
