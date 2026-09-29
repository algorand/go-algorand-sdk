package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/v2/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/v2/internal/signing"
	"github.com/algorand/go-algorand-sdk/v2/types"
)

// tgidPrefix is prepended to a transaction group when computing the group ID
var tgidPrefix = []byte("TG")

// bidPrefix is prepended to a bid when signing it
var bidPrefix = []byte("aB")

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
	return TransactionIDString(tx)
}

// TransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from: the canonical msgpack encoding of tx, prefixed with
// "TX" for domain separation.
//
// Signer implementations outside this package should build the bytes they sign
// with this function, so that all of them commit to the same byte sequence.
func TransactionBytesToSign(tx types.Transaction) []byte {
	return signing.TransactionBytesToSign(tx)
}

// TransactionID is the unique identifier for a Transaction in progress
func TransactionID(tx types.Transaction) (txid []byte) {
	txid32 := sha512.Sum512_256(TransactionBytesToSign(tx))
	return txid32[:]
}

// TransactionIDString is a base32 representation of a TransactionID
func TransactionIDString(tx types.Transaction) string {
	return signing.TxIDFromBytesToSign(TransactionBytesToSign(tx))
}

// VerifyBytes verifies that the signature is valid
func VerifyBytes(pk ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pk, signing.BytesToSign(message), signature)
}

// ed25519SignBid accepts an Ed25519Signer and a bid, and returns the signature
// of the bid under that key
func ed25519SignBid(sgnr Ed25519Signer, bid types.Bid) (signedBid []byte, err error) {
	// Encode the bid as msgpack
	encodedBid := msgpack.Encode(bid)

	// Prepend the hashable prefix
	msgParts := [][]byte{bidPrefix, encodedBid}
	toBeSigned := bytes.Join(msgParts, nil)

	// Sign the encoded bid
	s, err := signing.Ed25519RawSignature(sgnr, toBeSigned)
	if err != nil {
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

// MergeMultisigTransactions merges the given (partially) signed multisig transactions, and
// returns an encoded signed multisig transaction with the component signatures.
func MergeMultisigTransactions(stxsBytes ...[]byte) (txid string, stxBytes []byte, err error) {
	return signing.MergeMultisigTransactions(stxsBytes...)
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

		txID := sha512.Sum512_256(TransactionBytesToSign(tx))
		group.TxGroupHashes = append(group.TxGroupHashes, txID)
	}

	encoded := msgpack.Encode(group)

	// Prepend the hashable prefix and hash it
	msgParts := [][]byte{tgidPrefix, encoded}
	return sha512.Sum512_256(bytes.Join(msgParts, nil)), nil
}

/* LogicSig support */

// sanityCheckProgram performs heuristic program validation:
// check if passed in bytes are Algorand address or is B64 encoded, rather than Teal bytes
func sanityCheckProgram(program []byte) error {
	return signing.SanityCheckProgram(program)
}

// VerifyLogicSig verifies that a LogicSig contains a valid program and, if a
// delegated signature is present, that the signature is valid.
//
// The singleSigner argument is only used in the case of a delegated LogicSig
// whose delegating account is backed by a single private key (i.e. not a
// multisig account). In that case, it should be the address of the delegating
// account.
//
// Deprecated: This function is unsupported and unmaintained. Without the
// `falcon` build tag, PQ signatures are only checked against the delegating
// address and are not cryptographically validated.
func VerifyLogicSig(lsig types.LogicSig, singleSigner types.Address) (result bool) {
	if err := sanityCheckProgram(lsig.Logic); err != nil {
		return false
	}

	hasSig, hasMsig, hasLMsig, hasPQsig, err := lsigSignatures(lsig)
	if err != nil {
		return false
	}

	if hasSig {
		toBeSigned := signing.ProgramToSign(lsig.Logic)
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
		toBeSigned := signing.ProgramToSign(lsig.Logic)
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
		toBeSigned := signing.MsigProgramToSign(addr, lsig.Logic)
		return VerifyMultisig(addr, toBeSigned, lsig.LMsig)
	}

	if hasPQsig {
		addr := PQAddressFromSig(lsig.PQsig)
		if singleSigner != addr {
			return false
		}
		return verifyPQDelegation(signing.PQProgramToSign(addr, lsig.Logic), lsig.PQsig)
	}
	// the lsig account is the hash of its program bytes, nothing left to verify
	return true
}

// lsigSignatures reports which of the mutually exclusive delegation signatures
// the LogicSig carries. It errors out if more than one of them is present.
func lsigSignatures(lsig types.LogicSig) (hasSig, hasMsig, hasLMsig, hasPQsig bool, err error) {
	hasSig, hasMsig, hasLMsig, hasPQsig, count := lsig.SignatureCount()
	if count > 1 {
		err = errLsigTooManySignatures
	}
	return
}

// signLogicSigTransactionWithAddress signs a transaction with a LogicSig.
//
// lsigAddress is the address of the account that the LogicSig represents.
func signLogicSigTransactionWithAddress(lsig types.LogicSig, lsigAddress types.Address, tx types.Transaction) (txid string, stxBytes []byte, err error) {

	if !VerifyLogicSig(lsig, lsigAddress) {
		err = errLsigInvalidSignature
		return
	}

	txid = TransactionIDString(tx)
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
	hasSig, _, hasLMsig, hasPQsig, err := lsigSignatures(lsig)
	if err != nil {
		return "", nil, err
	}

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

// PQAddressFromSig returns the address of the account that performed a given PQ
// signature.
//
// The salt carried by the envelope is used as-is. Consensus does not require it
// to be the canonical salt for the envelope's scheme and public key, so an
// account on a non-canonical salt is a real account that this must resolve
// correctly, even though this SDK will only ever sign for canonical ones.
func PQAddressFromSig(sig types.PQSig) (addr types.Address) {
	return signing.PQAddressWithSalt(sig.PublicKey, sig.Scheme, sig.Salt)
}

// AddressFromProgram returns escrow account address derived from TEAL bytecode
func AddressFromProgram(program []byte) types.Address {
	toBeHashed := signing.ProgramToSign(program)
	hash := sha512.Sum512_256(toBeHashed)
	return types.Address(hash)
}

// makeLogicSig produces a new LogicSig signature.
//
// The function can work in three modes:
// 1. If no sgnr and ma provided then it returns contract-only LogicSig
// 2. If no ma provides, it returns Sig delegated LogicSig
// 3. If both sgnr and ma specified the function returns Multisig delegated LogicSig
func makeLogicSig(program []byte, args [][]byte, sgnr Ed25519Signer, ma MultisigAccount) (types.LogicSig, error) {
	if sgnr == nil && ma.Blank() {
		return signing.EscrowLogicSig(program, args)
	}

	if ma.Blank() {
		return signing.Ed25519DelegatedLogicSig(program, args, sgnr)
	}

	return signing.Ed25519MultisigDelegatedLogicSig(program, args, signing.MultisigAccount(ma), sgnr)
}

// TealVerify verifies signatures generated by TealSign and TealSignFromProgram
func TealVerify(pk ed25519.PublicKey, data []byte, contractAddress types.Address, rawSig types.Signature) bool {
	return ed25519.Verify(pk, signing.TealSignData(data, contractAddress), rawSig[:])
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
	return signing.IsEdwards25519Point(encoded)
}
