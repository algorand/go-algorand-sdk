package transaction

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-algorand-sdk/crypto"
)

func TestMakePaymentTxn(t *testing.T) {
	const fromAddress = "R5KMWJSLZRBBJJBDDNCJLF7AR436L3OMLIC2AYGBOXCU7N244OCLNZFM2M"
	const referenceTxID = "SGQTHZ4NF47OEHNUKN4SPGJOEVLFIMW2GILZRVU347YOYYRDSVAA"
	var fromSK = []byte{242, 175, 163, 193, 109, 239, 243, 150, 57, 236, 107, 130, 11, 20, 250, 252, 116, 163, 125, 222, 50, 175, 14, 232, 7, 153, 82, 169, 228, 5, 76, 247, 143, 84, 203, 38, 75, 204, 66, 20, 164, 35, 27, 68, 149, 151, 224, 143, 55, 229, 237, 204, 90, 5, 160, 96, 193, 117, 197, 79, 183, 92, 227, 132}

	txn, err := MakePaymentTxn(fromAddress, fromAddress, 10, 10, 1000, 1000, nil, "", "")
	require.NoError(t, err)

	id, bytes, err := crypto.SignTransaction(fromSK, txn)

	stxBytes := []byte{130, 163, 115, 105, 103, 196, 64, 131, 118, 119, 11, 135, 24, 77, 7, 112, 40, 243, 142, 37, 135, 67, 134, 136, 191, 0, 29, 231, 196, 61, 179, 87, 218, 72, 35, 51, 136, 90, 21, 28, 20, 46, 187, 156, 253, 174, 221, 29, 32, 35, 191, 204, 151, 214, 104, 130, 179, 128, 91, 234, 165, 10, 125, 202, 69, 175, 56, 134, 162, 222, 13, 163, 116, 120, 110, 135, 163, 97, 109, 116, 10, 163, 102, 101, 101, 205, 7, 38, 162, 102, 118, 205, 3, 232, 162, 108, 118, 205, 3, 232, 163, 114, 99, 118, 196, 32, 143, 84, 203, 38, 75, 204, 66, 20, 164, 35, 27, 68, 149, 151, 224, 143, 55, 229, 237, 204, 90, 5, 160, 96, 193, 117, 197, 79, 183, 92, 227, 132, 163, 115, 110, 100, 196, 32, 143, 84, 203, 38, 75, 204, 66, 20, 164, 35, 27, 68, 149, 151, 224, 143, 55, 229, 237, 204, 90, 5, 160, 96, 193, 117, 197, 79, 183, 92, 227, 132, 164, 116, 121, 112, 101, 163, 112, 97, 121}
	require.Equal(t, stxBytes, bytes)

	require.Equal(t, referenceTxID, id)
}
