package transaction

import (
	"github.com/algorand/go-algorand-sdk/types"
)

// MakePaymentTxn constructs a payment transaction using the passed parameters.
// `from` and `to` addresses should be checksummed, human-readable addresses
func MakePaymentTxn(from, to string, fee, amount, firstRound, lastRound uint64, note []byte, closeRemainderTo, genesisID string) (tx types.Transaction, err error) {
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
	tx = types.Transaction{
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
			Receiver:         toAddr,
			Amount:           types.Algos(amount),
			CloseRemainderTo: closeRemainderToAddr,
		},
	}
	return
}
