package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"filippo.io/edwards25519"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/types"
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

// msigProgramPrefix is prepended to a logic program when computing a hash for a program signed by multisig
var msigProgramPrefix = []byte("MsigProgram")

// programDataPrefix is prepended to teal sign data
var programDataPrefix = []byte("ProgData")

// appIDPrefix is prepended to application IDs in order to compute addresses
var appIDPrefix = []byte("appID")

// StateProofMessagePrefix is prepended to the canonical msgpack encoded state proof message when computing its hash.
var StateProofMessagePrefix = []byte("spm")

// LightBlockHeaderPrefix is prepended to the canonical msgpack encoded light block header when computing its vector commitment leaf.
var LightBlockHeaderPrefix = []byte("B256")

// RandomBytes fills the passed slice with randomness, and panics if it is
// unable to do so
func RandomBytes(s []byte) {
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
}

// GetTxID returns the txid of a transaction
func GetTxID(tx types.Transaction) string {
	rawTx := rawTransactionBytesToSign(tx)
	return txIDFromRawTxnBytesToSign(rawTx)
}

// Ed25519SignTransaction accepts an elliptic curve signer and a transaction,
// and returns the bytes of a signed transaction ready to be broadcasted to the
// network
// If the SK's corresponding address is different than the txn sender's, the SK's
// corresponding address will be assigned as AuthAddr
func Ed25519SignTransaction(sgnr Ed25519Signer, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	s, txid, err := rawSignTransaction(sgnr, tx)
	if err != nil {
		return
	}
	// Construct the SignedTxn
	stx := types.SignedTxn{
		Sig: s,
		Txn: tx,
	}

	a := types.Address(sgnr.Ed25519PublicKey())
	if stx.Txn.Sender != a {
		stx.AuthAddr = a
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

// txID computes a transaction id base32 string from raw transaction bytes
func txIDFromRawTxnBytesToSign(toBeSigned []byte) (txid string) {
	txidBytes := sha512.Sum512_256(toBeSigned)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}

// txIDFromTransaction is a convenience function for generating txID from txn
func txIDFromTransaction(tx types.Transaction) (txid string) {
	txidBytes := TransactionID(tx)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}

// TransactionID is the unique identifier for a Transaction in progress
func TransactionID(tx types.Transaction) (txid []byte) {
	toBeSigned := rawTransactionBytesToSign(tx)
	txid32 := sha512.Sum512_256(toBeSigned)
	txid = txid32[:]
	return
}

// TransactionIDString is a base32 representation of a TransactionID
func TransactionIDString(tx types.Transaction) (txid string) {
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(TransactionID(tx))
	return
}

// rawSignTransaction signs the msgpack-encoded tx (with prepended "TX" prefix), and returns the sig and txid
func rawSignTransaction(sgnr Ed25519Signer, tx types.Transaction) (s types.Signature, txid string, err error) {
	toBeSigned := rawTransactionBytesToSign(tx)

	// Sign the encoded transaction
	signature, err := sgnr.Ed25519Sign(toBeSigned)
	if err != nil {
		return
	}

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

// Ed25519SignBytes signs the bytes and returns the signature
func Ed25519SignBytes(sgnr Ed25519Signer, bytesToSign []byte) (signature []byte, err error) {
	// prepend the prefix for signing bytes
	toBeSigned := bytes.Join([][]byte{bytesPrefix, bytesToSign}, nil)

	// sign the bytes
	signature, err = sgnr.Ed25519Sign(toBeSigned)
	return
}

// VerifyBytes verifies that the signature is valid
func VerifyBytes(pk ed25519.PublicKey, message, signature []byte) bool {
	msgParts := [][]byte{bytesPrefix, message}
	toBeVerified := bytes.Join(msgParts, nil)
	return ed25519.Verify(pk, toBeVerified, signature)
}

// Ed25519SignBid accepts an Ed25519Signer and a bid, and returns the signature
// of the bid under that key
func Ed25519SignBid(sgnr Ed25519Signer, bid types.Bid) (signedBid []byte, err error) {
	// Encode the bid as msgpack
	encodedBid := msgpack.Encode(bid)

	// Prepend the hashable prefix
	msgParts := [][]byte{bidPrefix, encodedBid}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded bid
	sig, err := sgnr.Ed25519Sign(toBeSigned)

	if err != nil {
		return
	}

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
func multisigSingle(sgnr Ed25519Signer, ma MultisigAccount, customSigner signer) (msig types.MultisigSig, myIndex int, err error) {
	// check that sgnr.pk exists in the list of public keys in MultisigAccount ma
	myIndex = len(ma.Pks)
	myPublicKey := sgnr.Ed25519PublicKey()
	for i := 0; i < len(ma.Pks); i++ {
		if bytes.Equal(myPublicKey[:], ma.Pks[i]) {
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

// Ed25519SignMultisigTransaction signs the given transaction, and multisig
// preimage, with the signer, returning the bytes of a signed transaction with
// the multisig field partially populated, ready to be passed to other multisig
// signers to sign or broadcast.
func Ed25519SignMultisigTransaction(sgnr Ed25519Signer, ma MultisigAccount, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	err = ma.Validate()
	if err != nil {
		return
	}

	// this signer signs a transaction and sets txid from the closure
	customSigner := func() (rawSig types.Signature, err error) {
		rawSig, txid, err = rawSignTransaction(sgnr, tx)
		return rawSig, err
	}

	sig, _, err := multisigSingle(sgnr, ma, customSigner)
	if err != nil {
		return
	}

	// Encode the signedTxn
	stx := types.SignedTxn{
		Msig: sig,
		Txn:  tx,
	}

	maAddress, err := ma.Address()
	if err != nil {
		return
	}

	if stx.Txn.Sender != maAddress {
		stx.AuthAddr = maAddress
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
			err = errMsigMergeKeysMismatch
			return
		}

		if partStx.AuthAddr != refAuthAddr {
			err = errMsigMergeAuthAddrMismatch
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
					err = errMsigMergeInvalidDups
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
	txid = txIDFromTransaction(refTx)
	return
}

// Ed25519AppendMultisigTransaction appends the signature corresponding to the
// given signer, returning an encoded signed multisig transaction including the
// signature.  While we could compute the multisig preimage from the multisig
// blob, we ask the caller to pass it back in, to explicitly check that they
// know who they are signing as.
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
	txid, stxBytes, err = MergeMultisigTransactions(partStxBytes, preStxBytes)
	return
}

// VerifyMultisig verifies an assembled MultisigSig
//
// addr is the address of the Multisig account
// message is the bytes there were signed
// msig is the Multisig signature to verify
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

	return verifiedCount >= msig.Threshold
}

// ComputeGroupID returns group ID for a group of transactions
func ComputeGroupID(txgroup []types.Transaction) (gid types.Digest, err error) {
	if len(txgroup) > types.MaxTxGroupSize {
		err = fmt.Errorf("txgroup too large, %v > max size %v", len(txgroup), types.MaxTxGroupSize)
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

func isASCIIPrintableByte(symbol byte) bool {
	isBreakLine := symbol == '\n'
	isStdPrintable := symbol >= ' ' && symbol <= '~'
	return isBreakLine || isStdPrintable
}

func isASCIIPrintable(program []byte) bool {
	for _, b := range program {
		if !isASCIIPrintableByte(b) {
			return false
		}
	}
	return true
}

// sanityCheckProgram performs heuristic program validation:
// check if passed in bytes are Algorand address or is B64 encoded, rather than Teal bytes
func sanityCheckProgram(program []byte) error {
	if len(program) == 0 {
		return fmt.Errorf("empty program")
	}
	if isASCIIPrintable(program) {
		if _, err := types.DecodeAddress(string(program)); err == nil {
			return fmt.Errorf("requesting program bytes, get Algorand address")
		}
		if _, err := base64.StdEncoding.DecodeString(string(program)); err == nil {
			return fmt.Errorf("program should not be b64 encoded")
		}
		return fmt.Errorf("program bytes are all ASCII printable characters, not looking like Teal byte code")
	}
	return nil
}

// VerifyLogicSig verifies that a LogicSig contains a valid program and, if a
// delegated signature is present, that the signature is valid.
//
// The singleSigner argument is only used in the case of a delegated LogicSig
// whose delegating account is backed by a single private key (i.e. not a
// multisig account). In that case, it should be the address of the delegating
// account.
//
// Deprecated: This function is unsupported and unmantained. PQ signatures will
// not be validated and will always be treated as valid
func VerifyLogicSig(lsig types.LogicSig, singleSigner types.Address) (result bool) {
	if err := sanityCheckProgram(lsig.Logic); err != nil {
		return false
	}

	hasSig, hasMsig, hasLMsig, hasPQsig, count := lsig.SignatureCount()
	if count > 1 {
		return false
	}

	if hasSig {
		toBeSigned := programToSign(lsig.Logic)
		return ed25519.Verify(singleSigner[:], toBeSigned, lsig.Sig[:])
	}

	if hasMsig {
		msigAccount, err := MultisigAccountFromSig(lsig.Msig)
		if err != nil {
			return false
		}
		addr, err := msigAccount.Address()
		if err != nil {
			return false
		}
		toBeSigned := programToSign(lsig.Logic)
		return VerifyMultisig(addr, toBeSigned, lsig.Msig)
	}

	if hasLMsig {
		msigAccount, err := MultisigAccountFromSig(lsig.LMsig)
		if err != nil {
			return false
		}
		addr, err := msigAccount.Address()
		if err != nil {
			return false
		}
		toBeSigned := msigProgramToSign(addr, lsig.Logic)
		return VerifyMultisig(addr, toBeSigned, lsig.LMsig)
	}

	if hasPQsig {
		return true
	}
	// the lsig account is the hash of its program bytes, nothing left to verify
	return true
}

// pqsigProgramToSign returns the bytes a post-quantum scheme signs when
// delegating a LogicSig to a PQ account: ("PQProgram" || address ||
// program).
func pqsigProgramToSign(addr types.Address, program []byte) []byte {
	parts := [][]byte{pqProgramPrefix, addr[:], program}
	return bytes.Join(parts, nil)
}

// signLogicSigTransactionWithAddress signs a transaction with a LogicSig.
//
// lsigAddress is the address of the account that the LogicSig represents.
func signLogicSigTransactionWithAddress(lsig types.LogicSig, lsigAddress types.Address, tx types.Transaction) (txid string, stxBytes []byte, err error) {

	if !VerifyLogicSig(lsig, lsigAddress) {
		err = errLsigInvalidSignature
		return
	}

	txid = txIDFromTransaction(tx)
	// Construct the SignedTxn
	stx := types.SignedTxn{
		Lsig: lsig,
		Txn:  tx,
	}

	if stx.Txn.Sender != lsigAddress {
		stx.AuthAddr = lsigAddress
	}

	// Encode the SignedTxn
	stxBytes = msgpack.Encode(stx)
	return
}

// SignLogicSigAccountTransaction signs a transaction with a LogicSigAccount. It
// returns the TxID of the signed transaction and the raw bytes ready to be
// broadcast to the network. Note: any type of transaction can be signed by a
// LogicSig, but the network will reject the transaction if the LogicSig's
// program declines the transaction.
func SignLogicSigAccountTransaction(logicSigAccount LogicSigAccount, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	addr, err := logicSigAccount.Address()
	if err != nil {
		return
	}

	txid, stxBytes, err = signLogicSigTransactionWithAddress(logicSigAccount.Lsig, addr, tx)
	return
}

// SignLogicSigTransaction takes LogicSig object and a transaction and returns the
// bytes of a signed transaction ready to be broadcasted to the network
// Note, LogicSig actually can be attached to any transaction and it is a
// program's responsibility to approve/decline the transaction
//
// This function supports signing transactions with a sender that differs from
// the LogicSig's address, EXCEPT IF the LogicSig is delegated to a non-multisig
// account. In order to properly handle that case, create a LogicSigAccount and
// use SignLogicSigAccountTransaction instead.
func SignLogicSigTransaction(lsig types.LogicSig, tx types.Transaction) (txid string, stxBytes []byte, err error) {
	hasSig := lsig.Sig != (types.Signature{})
	hasLMsig := !lsig.LMsig.Blank()
	hasPQsig := !lsig.PQsig.Blank()

	// the address that the LogicSig represents
	var lsigAddress types.Address
	if hasSig {
		// For a LogicSig with a non-multisig delegating account, we cannot derive
		// the address of that account from only its signature, so assume the
		// delegating account is the sender. If that's not the case, the signing
		// will fail.
		lsigAddress = tx.Header.Sender
	} else if hasLMsig {
		var msigAccount MultisigAccount
		msigAccount, err = MultisigAccountFromSig(lsig.LMsig)
		if err != nil {
			return
		}
		lsigAddress, err = msigAccount.Address()
		if err != nil {
			return
		}
	} else if hasPQsig {
		lsigAddress = PQAddressFromSig(lsig.PQsig)
	} else {
		lsigAddress = LogicSigAddress(lsig)
	}

	txid, stxBytes, err = signLogicSigTransactionWithAddress(lsig, lsigAddress, tx)
	return
}

// PQAddressFromSig returns the address of the account that performed a given PQ signature
func PQAddressFromSig(sig types.PQSig) (addr types.Address) {
	buf := make([]byte, 0, len(pqAddressPrefix)+len(sig.Scheme)+1+len(sig.PublicKey))
	buf = append(buf, pqAddressPrefix...)
	buf = append(buf, sig.Scheme[:]...)
	buf = append(buf, uint8(sig.Salt))
	buf = append(buf, sig.PublicKey[:]...)

	digest := sha512.Sum512_256(buf)

	copy(addr[:], digest[:])

	return
}

func programToSign(program []byte) []byte {
	parts := [][]byte{programPrefix, program}
	toBeSigned := bytes.Join(parts, nil)
	return toBeSigned
}

func msigProgramToSign(msigAddr types.Address, program []byte) []byte {
	parts := [][]byte{msigProgramPrefix, msigAddr[:], program}
	toBeSigned := bytes.Join(parts, nil)
	return toBeSigned
}

func signProgram(sgnr Ed25519Signer, program []byte) (sig types.Signature, err error) {
	toBeSigned := programToSign(program)
	rawSig, err := sgnr.Ed25519Sign(toBeSigned)
	if err != nil {
		return
	}

	n := copy(sig[:], rawSig)
	if n != len(sig) {
		err = errInvalidSignatureReturned
		return
	}
	return
}

// AddressFromProgram returns escrow account address derived from TEAL bytecode
func AddressFromProgram(program []byte) types.Address {
	toBeHashed := programToSign(program)
	hash := sha512.Sum512_256(toBeHashed)
	return types.Address(hash)
}

// makeLogicSig produces a new LogicSig signature.
//
// The function can work in three modes:
// 1. If no sgnr and ma provided then it returns contract-only LogicSig
// 2. If no ma provides, it returns Sig delegated LogicSig
// 3. If both sgnr and ma specified the function returns Multisig delegated LogicSig
func makeLogicSig(program []byte, args [][]byte, sgnr Ed25519Signer, ma MultisigAccount) (lsig types.LogicSig, err error) {
	if err = sanityCheckProgram(program); err != nil {
		return
	}

	if sgnr == nil && ma.Blank() {
		lsig.Logic = program
		lsig.Args = args
		return
	}

	if ma.Blank() {
		var sig types.Signature
		sig, err = signProgram(sgnr, program)
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

	multisigAddr, err := ma.Address()
	if err != nil {
		return
	}

	// this signer signs a program
	customSigner := func() (rawSig types.Signature, err error) {
		toBeSigned := msigProgramToSign(multisigAddr, program)
		sigBytes, err := sgnr.Ed25519Sign(toBeSigned)
		if err != nil {
			return
		}

		copy(rawSig[:], sigBytes)
		return
	}

	msig, _, err := multisigSingle(sgnr, ma, customSigner)
	if err != nil {
		return
	}

	lsig.Logic = program
	lsig.Args = args
	lsig.LMsig = msig

	return
}

// Ed25519AppendMultisigToLogicSig adds a new signature to multisigned LogicSig
func Ed25519AppendMultisigToLogicSig(lsig *types.LogicSig, sgnr Ed25519Signer) error {
	if lsig.LMsig.Blank() {
		return errLsigEmptyMsig
	}

	ma, err := MultisigAccountFromSig(lsig.LMsig)
	if err != nil {
		return err
	}

	multisigAddr, err := ma.Address()
	if err != nil {
		return err
	}

	customSigner := func() (rawSig types.Signature, err error) {
		toBeSigned := msigProgramToSign(multisigAddr, lsig.Logic)
		sigBytes, err := sgnr.Ed25519Sign(toBeSigned)
		if err != nil {
			return
		}

		copy(rawSig[:], sigBytes)
		return
	}

	msig, idx, err := multisigSingle(sgnr, ma, customSigner)
	if err != nil {
		return err
	}

	lsig.LMsig.Subsigs[idx] = msig.Subsigs[idx]

	return nil
}

// Ed25519TealSign creates a signature compatible with ed25519verify opcode from
// contract address
func Ed25519TealSign(sgnr Ed25519Signer, data []byte, contractAddress types.Address) (rawSig types.Signature, err error) {
	msgParts := [][]byte{programDataPrefix, contractAddress[:], data}
	toBeSigned := bytes.Join(msgParts, nil)

	signature, err := sgnr.Ed25519Sign(toBeSigned)
	if err != nil {
		return
	}

	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	n := copy(rawSig[:], signature)
	if n != len(rawSig) {
		err = errInvalidSignatureReturned
	}
	return
}

// Ed25519TealSignFromProgram creates a signature compatible with ed25519verify
// opcode from raw program bytes
func Ed25519TealSignFromProgram(sgnr Ed25519Signer, data []byte, program []byte) (rawSig types.Signature, err error) {
	addr := AddressFromProgram(program)
	return Ed25519TealSign(sgnr, data, addr)
}

// TealVerify verifies signatures generated by TealSign and TealSignFromProgram
func TealVerify(pk ed25519.PublicKey, data []byte, contractAddress types.Address, rawSig types.Signature) bool {
	msgParts := [][]byte{programDataPrefix, contractAddress[:], data}
	toBeVerified := bytes.Join(msgParts, nil)

	return ed25519.Verify(pk, toBeVerified, rawSig[:])
}

// GetApplicationAddress returns the address corresponding to an application's escrow account.
func GetApplicationAddress(appID uint64) types.Address {
	encodedAppID := make([]byte, 8)
	binary.BigEndian.PutUint64(encodedAppID, appID)

	parts := [][]byte{appIDPrefix, encodedAppID}
	toBeHashed := bytes.Join(parts, nil)

	hash := sha512.Sum512_256(toBeHashed)
	return types.Address(hash)
}

// HashStateProofMessage returns the hash of a state proof message.
func HashStateProofMessage(stateProofMessage *types.Message) types.MessageHash {
	msgPackedStateProofMessage := msgpack.Encode(stateProofMessage)

	stateProofMessageData := make([]byte, 0, len(StateProofMessagePrefix)+len(msgPackedStateProofMessage))
	stateProofMessageData = append(stateProofMessageData, StateProofMessagePrefix...)
	stateProofMessageData = append(stateProofMessageData, msgPackedStateProofMessage...)

	return sha256.Sum256(stateProofMessageData)
}

// HashLightBlockHeader returns the hash of a light block header.
func HashLightBlockHeader(lightBlockHeader types.LightBlockHeader) types.Digest {
	msgPackedLightBlockHeader := msgpack.Encode(lightBlockHeader)

	lightBlockHeaderData := make([]byte, 0, len(LightBlockHeaderPrefix)+len(msgPackedLightBlockHeader))
	lightBlockHeaderData = append(lightBlockHeaderData, LightBlockHeaderPrefix...)
	lightBlockHeaderData = append(lightBlockHeaderData, msgpack.Encode(lightBlockHeader)...)

	return sha256.Sum256(lightBlockHeaderData)
}

// IsEdwards25519Point reports whether encoded can be decoded as an
// Edwards25519 curve point.
func IsEdwards25519Point(encoded []byte) bool {
	if len(encoded) != 32 {
		return false
	}
	_, err := new(edwards25519.Point).SetBytes(encoded)
	return err == nil
}
