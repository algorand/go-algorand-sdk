package transaction

import (
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

const MinTxnFee = 1000 // v5 consensus params, in microAlgos

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
// fee is fee per byte as received from algod SuggestedFee API call
func MakePaymentTxn(from, to string, fee, amount, firstRound, lastRound uint64, note []byte, closeRemainderTo, genesisID string, genesisHash []byte) (types.Transaction, error) {
	// Decode from address
	fromAddr, err := types.DecodeAddress(from)
	if err != nil {
		return types.Transaction{}, err
	}

	// Decode to address
	toAddr, err := types.DecodeAddress(to)
	if err != nil {
		return types.Transaction{}, err
	}

	// Decode the CloseRemainderTo address, if present
	var closeRemainderToAddr types.Address
	if closeRemainderTo != "" {
		closeRemainderToAddr, err = types.DecodeAddress(closeRemainderTo)
		if err != nil {
			return types.Transaction{}, err
		}
	}

	// Decode GenesisHash
	if len(genesisHash) == 0 {
		return types.Transaction{}, fmt.Errorf("payment transaction must contain a genesisHash")
	}

	var gh types.Digest
	copy(gh[:], genesisHash)

	// Build the transaction
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:      fromAddr,
			Fee:         types.MicroAlgos(fee),
			FirstValid:  types.Round(firstRound),
			LastValid:   types.Round(lastRound),
			Note:        note,
			GenesisID:   genesisID,
			GenesisHash: gh,
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver:         toAddr,
			Amount:           types.MicroAlgos(amount),
			CloseRemainderTo: closeRemainderToAddr,
		},
	}

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * fee)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// MakeKeyRegTxn constructs a keyreg transaction using the passed parameters.
// - account is a checksummed, human-readable address for which we register the given participation key.
// - fee is fee per byte as received from algod SuggestedFee API call.
// - firstRound is the first round this txn is valid (txn semantics unrelated to key registration)
// - lastRound is the last round this txn is valid
// - genesis id corresponds to the id of the network
// - genesis hash corresponds to the base64-encoded hash of the genesis of the network
// KeyReg parameters:
// - votePK is a base64-encoded string corresponding to the root participation public key
// - selectionKey is a base64-encoded string corresponding to the vrf public key
// - voteFirst is the first round this participation key is valid
// - voteLast is the last round this participation key is valid
// - voteKeyDilution is the dilution for the 2-level participation key
func MakeKeyRegTxn(account string, feePerByte, firstRound, lastRound uint64, genesisID string, genesisHash string,
	voteKey, selectionKey string, voteFirst, voteLast, voteKeyDilution uint64) (types.Transaction, error) {
	// Decode account address
	accountAddr, err := types.DecodeAddress(account)
	if err != nil {
		return types.Transaction{}, err
	}

	gh, err := keysLenBytesFromBase64(genesisHash)
	if err != nil {
		return types.Transaction{}, err
	}

	votePK, err := keysLenBytesFromBase64(voteKey)
	if err != nil {
		return types.Transaction{}, err
	}

	selectionPK, err := keysLenBytesFromBase64(selectionKey)
	if err != nil {
		return types.Transaction{}, err
	}

	tx := types.Transaction{
		Type: types.KeyRegistrationTx,
		Header: types.Header{
			Sender:      accountAddr,
			Fee:         types.MicroAlgos(feePerByte),
			FirstValid:  types.Round(firstRound),
			LastValid:   types.Round(lastRound),
			GenesisHash: gh,
			GenesisID:   genesisID,
		},
		KeyregTxnFields: types.KeyregTxnFields{
			VotePK:          votePK,
			SelectionPK:     selectionPK,
			VoteFirst:       types.Round(voteFirst),
			VoteLast:        types.Round(voteLast),
			VoteKeyDilution: voteKeyDilution,
		},
	}

	// Update fee
	eSize, err := estimateSize(tx)
	if err != nil {
		return types.Transaction{}, err
	}
	tx.Fee = types.MicroAlgos(eSize * feePerByte)

	if tx.Fee < MinTxnFee {
		tx.Fee = MinTxnFee
	}

	return tx, nil
}

// EstimateSize returns the estimated length of the encoded transaction
func estimateSize(txn types.Transaction) (uint64, error) {
	key := crypto.GenerateAccount()
	_, stx, err := crypto.SignTransaction(key.PrivateKey, txn)
	if err != nil {
		return 0, err
	}
	return uint64(len(stx)), nil
}

// keysLenBytesFromBase64 decodes the input base64 string and outputs a
// 32 byte array, erroring if the input is the wrong length.
func keysLenBytesFromBase64(in string) (out [types.KeysLenBytes]byte, err error) {
	slice, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return
	}
	if len(slice) != types.KeysLenBytes {
		return out, fmt.Errorf("Input is not %v bytes", types.KeysLenBytes)
	}
	copy(out[:], slice)
	return
}
