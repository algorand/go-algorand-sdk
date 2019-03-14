package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"

	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"
	"github.com/algorand/go-algorand-sdk/types"
)

func TestSignTransaction(t *testing.T) {
	// Reference transaction + keys generated using goal
	const referenceSignedTxn = "82a3736967c4402f7d02826bc77dcd2a6e4d098ddcb619c4670c1dd98eba9a96f8d9a56e4fe8ff9868cee08ef1eae822bca9e99353244402717ad5850fd8136e0652f7295bd10da374786e87a3616d74cd04d2a366656501a26676ce0001a04fa26c76ce0001a437a3726376c4207d3f99e53d34ae49eb2f458761cf538408ffdaee35c70d8234166de7abe3e517a3736e64c4201bd63dc672b0bb29d42fcafa3422a4d385c0c8169bb01595babf8855cf596979a474797065a3706179"
	const referenceTxid = "YGE4O2RBSMVPSPPXBK3SR45M453TRQA3L6U3GG7VYFLZL54Y4EZQ"
	const fromAddr = "DPLD3RTSWC5STVBPZL5DIIVE2OC4BSAWTOYBLFN2X6EFLT2ZNF4SMX64UA"
	const fromSK = "actress tongue harbor tray suspect odor load topple vocal avoid ignore apple lunch unknown tissue museum once switch captain place lemon sail outdoor absent creek"
	const toAddr = "PU7ZTZJ5GSXET2ZPIWDWDT2TQQEP7WXOGXDQ3ARUCZW6PK7D4ULSE6NYCE"

	// Build the unsigned transaction
	tx, err := transaction.MakePaymentTxn(fromAddr, toAddr, 1, 1234, 106575, 107575, nil, "")
	require.NoError(t, err)

	// Decode the secret key for the sender
	seed, err := mnemonic.ToKey(fromSK)
	require.NoError(t, err)
	sk := ed25519.NewKeyFromSeed(seed)

	// Check that we have the correct pk
	derivedFromPK := sk.Public()
	var derivedFromAddr types.Address
	copy(derivedFromAddr[:], derivedFromPK.(ed25519.PublicKey))
	require.Equal(t, fromAddr, derivedFromAddr.String())

	// Sign the transaction
	stxBytes, err := SignTransaction(sk, tx)
	require.NoError(t, err)
	stxHex := fmt.Sprintf("%x", stxBytes)
	require.Equal(t, referenceSignedTxn, stxHex)

	txid := getTxID(tx)
	require.Equal(t, txid, referenceTxid)
}
