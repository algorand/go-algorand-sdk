package transaction

import (
	"fmt"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
func MakePaymentTxn(from, to string, fee, amount, firstRound, lastRound int64, note []byte, genesisID string) (encoded []byte, err error) {

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

	// Build the transaction
	tx := types.Transaction{
		Type: types.PaymentTx,
		Header: types.Header{
			Sender:     fromAddr,
			Fee:        types.Algos(fee),
			FirstValid: types.Round(firstRound),
			LastValid:  types.Round(lastRound),
			Note:       note,
			GenesisID:  genesisID,
		},
		PaymentTxnFields: types.PaymentTxnFields{
			Receiver: toAddr,
			Amount:   types.Algos(amount),
		},
	}

	encoded = msgpack.Encode(tx)
	return
}
