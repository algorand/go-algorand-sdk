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

// GenerateAddressFromSK take a secret key and returns the corresponding Address
func GenerateAddressFromSK(sk []byte) (types.Address, error) {
	edsk := ed25519.PrivateKey(sk)

	var a types.Address
	pk := edsk.Public()
	n := copy(a[:], []byte(pk.(ed25519.PublicKey)))
	if n != ed25519.PublicKeySize {
		return [32]byte{}, fmt.Errorf("generated public key has the wrong size, expected %d, got %d", ed25519.PublicKeySize, n)
	}
	return a, nil
}

// GetTxID returns the txid of a transaction
func GetTxID(tx types.Transaction) string {
	rawTx := rawTransactionBytesToSign(tx)
	return txIDFromRawTxnBytesToSign(rawTx)
}

// SignTransaction accepts a private key and a transaction, and returns the
// bytes of a signed transaction ready to be broadcasted to the network
// If the SK's corresponding address is different than the txn sender's, the SK's
// corresponding address will be assigned as AuthAddr
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

	a, err := GenerateAddressFromSK(sk)
	if err != nil {
		return
	}

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

// VerifyBytes verifies that the signature is valid
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
// multsig account). In that case, it should be the address of the delegating
// account.
func VerifyLogicSig(lsig types.LogicSig, singleSigner types.Address) (result bool) {
	if err := sanityCheckProgram(lsig.Logic); err != nil {
		return false
	}

	hasSig := lsig.Sig != (types.Signature{})
	hasMsig := !lsig.Msig.Blank()

	// require only one or zero sig
	if hasSig && hasMsig {
		return false
	}

	toBeSigned := programToSign(lsig.Logic)

	if hasSig {
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
		return VerifyMultisig(addr, toBeSigned, lsig.Msig)
	}

	// the lsig account is the hash of its program bytes, nothing left to verify
	return true
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
	hasMsig := !lsig.Msig.Blank()

	// the address that the LogicSig represents
	var lsigAddress types.Address
	if hasSig {
		// For a LogicSig with a non-multisig delegating account, we cannot derive
		// the address of that account from only its signature, so assume the
		// delegating account is the sender. If that's not the case, the signing
		// will fail.
		lsigAddress = tx.Header.Sender
	} else if hasMsig {
		var msigAccount MultisigAccount
		msigAccount, err = MultisigAccountFromSig(lsig.Msig)
		if err != nil {
			return
		}
		lsigAddress, err = msigAccount.Address()
		if err != nil {
			return
		}
	} else {
		lsigAddress = LogicSigAddress(lsig)
	}

	txid, stxBytes, err = signLogicSigTransactionWithAddress(lsig, lsigAddress, tx)
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

// AddressFromProgram returns escrow account address derived from TEAL bytecode
func AddressFromProgram(program []byte) types.Address {
	toBeHashed := programToSign(program)
	hash := sha512.Sum512_256(toBeHashed)
	return types.Address(hash)
}

// makeLogicSig produces a new LogicSig signature.
//
// The function can work in three modes:
// 1. If no sk and ma provided then it returns contract-only LogicSig
// 2. If no ma provides, it returns Sig delegated LogicSig
// 3. If both sk and ma specified the function returns Multisig delegated LogicSig
func makeLogicSig(program []byte, args [][]byte, sk ed25519.PrivateKey, ma MultisigAccount) (lsig types.LogicSig, err error) {
	if err = sanityCheckProgram(program); err != nil {
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

// TealSign creates a signature compatible with ed25519verify opcode from contract address
func TealSign(sk ed25519.PrivateKey, data []byte, contractAddress types.Address) (rawSig types.Signature, err error) {
	msgParts := [][]byte{programDataPrefix, contractAddress[:], data}
	toBeSigned := bytes.Join(msgParts, nil)

	signature := ed25519.Sign(sk, toBeSigned)
	// Copy the resulting signature into a Signature, and check that it's
	// the expected length
	n := copy(rawSig[:], signature)
	if n != len(rawSig) {
		err = errInvalidSignatureReturned
	}
	return
}

// TealSignFromProgram creates a signature compatible with ed25519verify opcode from raw program bytes
func TealSignFromProgram(sk ed25519.PrivateKey, data []byte, program []byte) (rawSig types.Signature, err error) {
	addr := AddressFromProgram(program)
	return TealSign(sk, data, addr)
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
