package transaction

import (
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

// EstimateSize returns the estimated length of the encoded transaction
func estimateSize(txn types.Transaction) (uint64, error) {
	key := crypto.GenerateAccount()
	_, stx, err := crypto.SignTransaction(key.PrivateKey, txn)
	if err != nil {
		return 0, err
	}
	return uint64(len(stx)), nil
}
