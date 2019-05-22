package transaction

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

const minFee = 1000

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
func MakePaymentTxn(from, to string, fee, amount, firstRound, lastRound int64, note []byte, closeRemainderTo, genesisID string) (encoded []byte, err error) {

	// Sanity check for int64
	if fee < 0 ||
		amount < 0 ||
		firstRound < 0 ||
		lastRound < 0 {
		err = fmt.Errorf("all numbers must not be negative")
		return
	}

	// Decode from address
	fromAddr, err := types.DecodeAddress(from)
	if err != nil {
		return
	}

	// Decode to address
	toAddr, err := types.DecodeAddress(to)
	if err != nil {
		return
	}

	// Decode the CloseRemainderTo address, if present
	var closeRemainderToAddr types.Address
	if closeRemainderTo != "" {
		closeRemainderToAddr, err = types.DecodeAddress(closeRemainderTo)
		if err != nil {
			return
		}
	}

	// Build the transaction
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     fromAddr,
			Fee:        types.Algos(200 * fee),
			FirstValid: types.Round(firstRound),
			LastValid:  types.Round(lastRound),
			Note:       note,
			GenesisID:  genesisID,
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver:         toAddr,
			Amount:           types.Algos(amount),
			CloseRemainderTo: closeRemainderToAddr,
		},
	}

	// Get the right fee
	l, err := getEstimatedSize(tx)
	if err != nil {
		return nil, err
	}

	tmpFee := types.Algos(uint64(fee) * l)

	if tmpFee < minFee {
		tx.Fee = minFee
	}

	encoded = msgpack.Encode(tx)

	return
}

func getEstimatedSize(tx types.Transaction) (uint64, error) {
	key := crypto.GenerateSK()
	en, err := crypto.SignTransaction(key, msgpack.Encode(tx))
	if err != nil {
		return 0, err
	}

	return uint64(len(en)), nil
}
